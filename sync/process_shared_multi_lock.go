// Copyright © 2026 Yoshiki Shibata. All rights reserved.

package sync

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"syscall"
	"testing"
)

const (
	// processSharedCoordinatorFileName is the well-known file used to
	// serialize acquisition attempts across processes.
	processSharedCoordinatorFileName = "courier_process_shared_multi_lock_coordinator"

	// processSharedLockFilePrefix matches ProcessSharedLock so that a
	// ProcessSharedMultiLock and a plain ProcessSharedLock with the same
	// name refer to the same underlying file and therefore contend.
	processSharedLockFilePrefix = "lock_file_"
)

// ProcessSharedMultiLock represents a set of locks that are acquired and
// released together across processes atomically. Unlike acquiring flocks
// sequentially, this avoids the "hold and wait" problem: either all
// requested locks are acquired or none of them, so no process holds a
// subset while waiting for the rest.
type ProcessSharedMultiLock struct {
	requests []LockRequest
	files    []*os.File // one entry per request, aligned by index
	stack    string
	acquired bool
}

// NewProcessSharedMultiLock creates a new ProcessSharedMultiLock with the
// given lock requests. Requests are sorted by name for stable ordering and
// consistent log output.
func NewProcessSharedMultiLock(t *testing.T, requests ...LockRequest) *ProcessSharedMultiLock {
	sortedRequests := make([]LockRequest, len(requests))
	copy(sortedRequests, requests)
	sort.Slice(sortedRequests, func(i, j int) bool {
		return sortedRequests[i].Name < sortedRequests[j].Name
	})

	_, file2, line2, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Cannot call runtime.Caller(1)")
	}
	_, file3, line3, ok := runtime.Caller(2)
	if !ok {
		t.Fatalf("Cannot call runtime.Caller(2)")
	}
	stack := fmt.Sprintf("%s:%d/%s:%d",
		filepath.Base(file3), line3, filepath.Base(file2), line2)

	return &ProcessSharedMultiLock{
		requests: sortedRequests,
		stack:    stack,
		acquired: false,
	}
}

// AcquireAll acquires all requested locks across processes atomically.
// A single well-known coordinator lock file serializes acquisition
// attempts: while holding the coordinator (exclusive), each target lock
// is probed with a non-blocking flock. If any target is unavailable, all
// partial acquisitions are released, the coordinator is released, and
// the caller blocks on the failing target with a blocking flock. When
// that flock returns the target is released immediately and the whole
// protocol is retried under the coordinator - this replaces polling with
// an event-driven wait without ever holding a subset of the requested
// locks while waiting for the rest.
//
// The locks are automatically released when the test finishes via
// t.Cleanup.
func (ml *ProcessSharedMultiLock) AcquireAll(t *testing.T) {
	if ml.acquired {
		t.Fatalf("ProcessSharedMultiLock already acquired")
	}

	tempDir := os.TempDir()

	// Open all target lock files up front. Each request gets its own
	// open file description so flocks on the same name from different
	// ProcessSharedMultiLock instances are independent.
	files := make([]*os.File, len(ml.requests))
	for i, req := range ml.requests {
		path := filepath.Join(tempDir, processSharedLockFilePrefix+req.Name)
		f, err := os.Create(path)
		if err != nil {
			for j := range i {
				_ = files[j].Close()
			}
			t.Fatalf("os.Create(%s) failed: %v", path, err)
		}
		files[i] = f
	}
	ml.files = files

	ml.logWaiting()

	coordPath := filepath.Join(tempDir, processSharedCoordinatorFileName)
	coord, err := os.Create(coordPath)
	if err != nil {
		ml.closeFiles()
		t.Fatalf("os.Create(%s) failed: %v", coordPath, err)
	}
	defer func() { _ = coord.Close() }()

	for {
		// Take the coordinator so only one process is probing targets
		// at a time. This bounds the retry storm and makes the
		// "try all or none" check effectively atomic across processes.
		if err := syscall.Flock(int(coord.Fd()), syscall.LOCK_EX); err != nil {
			ml.closeFiles()
			t.Fatalf("syscall.Flock(coordinator, LOCK_EX) failed: %v", err)
		}

		gotAll := true
		heldCount := 0
		failingIdx := -1
		for i, req := range ml.requests {
			mode := syscall.LOCK_SH | syscall.LOCK_NB
			if req.Exclusive {
				mode = syscall.LOCK_EX | syscall.LOCK_NB
			}
			err := syscall.Flock(int(files[i].Fd()), mode)
			if err == nil {
				heldCount++
				continue
			}
			if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
				gotAll = false
				failingIdx = i
				break
			}
			// Unexpected error: unwind and fail the test.
			for j := range heldCount {
				_ = syscall.Flock(int(files[j].Fd()), syscall.LOCK_UN)
			}
			_ = syscall.Flock(int(coord.Fd()), syscall.LOCK_UN)
			ml.closeFiles()
			t.Fatalf("syscall.Flock(%s) failed: %v", req.Name, err)
		}

		if gotAll {
			if err := syscall.Flock(int(coord.Fd()), syscall.LOCK_UN); err != nil {
				for _, f := range files {
					_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
				}
				ml.closeFiles()
				t.Fatalf("syscall.Flock(coordinator, LOCK_UN) failed: %v", err)
			}
			break
		}

		// Partial acquisition: release what we grabbed and the coordinator
		// so other waiters can make progress.
		for j := range heldCount {
			_ = syscall.Flock(int(files[j].Fd()), syscall.LOCK_UN)
		}
		if err := syscall.Flock(int(coord.Fd()), syscall.LOCK_UN); err != nil {
			ml.closeFiles()
			t.Fatalf("syscall.Flock(coordinator, LOCK_UN) failed: %v", err)
		}

		// Instead of sleeping, block on the first target that was
		// unavailable using its desired mode. We hold no other locks at
		// this point, so this cannot introduce hold-and-wait. When the
		// blocking flock returns we release it immediately and retry the
		// full protocol under the coordinator - the brief hold gives
		// other pending waiters a chance to progress and avoids the
		// former 10ms polling interval.
		failingReq := ml.requests[failingIdx]
		waitMode := syscall.LOCK_SH
		if failingReq.Exclusive {
			waitMode = syscall.LOCK_EX
		}
		if err := syscall.Flock(int(files[failingIdx].Fd()), waitMode); err != nil {
			ml.closeFiles()
			t.Fatalf("syscall.Flock(%s, blocking) failed: %v", failingReq.Name, err)
		}
		if err := syscall.Flock(int(files[failingIdx].Fd()), syscall.LOCK_UN); err != nil {
			ml.closeFiles()
			t.Fatalf("syscall.Flock(%s, LOCK_UN) failed: %v", failingReq.Name, err)
		}
	}

	ml.acquired = true
	ml.logAcquired()

	t.Cleanup(func() { ml.ReleaseAll(t) })
}

// ReleaseAll releases every acquired lock and closes the underlying file
// handles. This method is idempotent - calling it multiple times is safe.
func (ml *ProcessSharedMultiLock) ReleaseAll(t *testing.T) {
	if !ml.acquired {
		return
	}

	for i, f := range ml.files {
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			t.Errorf("syscall.Flock(%s, LOCK_UN) failed: %v", ml.requests[i].Name, err)
		}
	}
	ml.closeFiles()
	ml.acquired = false
	ml.logReleased()
}

// closeFiles closes every open lock file and clears the slice.
func (ml *ProcessSharedMultiLock) closeFiles() {
	for _, f := range ml.files {
		if f != nil {
			_ = f.Close()
		}
	}
	ml.files = nil
}

func (ml *ProcessSharedMultiLock) logWaiting() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessSharedMultiLock[%s]: Wait for %s Lock(%s)", req.Name, lockType, ml.stack)
	}
}

func (ml *ProcessSharedMultiLock) logAcquired() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessSharedMultiLock[%s]: %s Lock is acquired(%s)", req.Name, lockType, ml.stack)
	}
}

func (ml *ProcessSharedMultiLock) logReleased() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessSharedMultiLock[%s]: %s Unlocked(%s)", req.Name, lockType, ml.stack)
	}
}
