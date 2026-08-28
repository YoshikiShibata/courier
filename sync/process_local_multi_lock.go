// Copyright © 2026 Yoshiki Shibata. All rights reserved.

package sync

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"testing"
)

// LockRequest represents a request to acquire a lock with a specific mode.
type LockRequest struct {
	Name      string
	Exclusive bool // true: exclusive (write) lock, false: shared (read) lock
}

// Shared creates a shared (read) lock request.
func Shared(name string) LockRequest {
	return LockRequest{Name: name, Exclusive: false}
}

// Exclusive creates an exclusive (write) lock request.
func Exclusive(name string) LockRequest {
	return LockRequest{Name: name, Exclusive: true}
}

// lockState tracks the state of a single lock.
type lockState struct {
	readersCount int
	writerActive bool
}

// processLocalMultiLockManager manages multiple locks and allows atomic
// acquisition of multiple locks. This prevents lock chain slowdowns by
// ensuring tests acquire all needed locks at once.
type processLocalMultiLockManager struct {
	mu    sync.Mutex
	cond  *sync.Cond
	locks map[string]*lockState
}

// globalProcessLocalMultiLockManager is the singleton instance used across all tests.
var globalProcessLocalMultiLockManager = newProcessLocalMultiLockManager()

func newProcessLocalMultiLockManager() *processLocalMultiLockManager {
	m := &processLocalMultiLockManager{
		locks: make(map[string]*lockState),
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// getOrCreateLockState returns the lock state for the given name, creating it if necessary.
// Must be called with m.mu held.
func (m *processLocalMultiLockManager) getOrCreateLockState(name string) *lockState {
	if ls, ok := m.locks[name]; ok {
		return ls
	}
	ls := &lockState{}
	m.locks[name] = ls
	return ls
}

// canAcquireAll checks if all requested locks can be acquired.
// Must be called with m.mu held.
// Note: ProcessLocalMultiLock does NOT use writer priority to avoid
// unnecessary blocking when waiting for multiple locks atomically.
func (m *processLocalMultiLockManager) canAcquireAll(requests []LockRequest) bool {
	for _, req := range requests {
		ls := m.getOrCreateLockState(req.Name)
		if req.Exclusive {
			// Exclusive lock requires no readers and no active writer
			if ls.readersCount > 0 || ls.writerActive {
				return false
			}
		} else {
			// Shared lock requires no active writer
			if ls.writerActive {
				return false
			}
		}
	}
	return true
}

// doAcquireAll acquires all requested locks.
// Must be called with m.mu held and after canAcquireAll returns true.
func (m *processLocalMultiLockManager) doAcquireAll(requests []LockRequest) {
	for _, req := range requests {
		ls := m.getOrCreateLockState(req.Name)
		if req.Exclusive {
			ls.writerActive = true
		} else {
			ls.readersCount++
		}
	}
}

// doReleaseAll releases all requested locks.
// Must be called with m.mu held.
func (m *processLocalMultiLockManager) doReleaseAll(requests []LockRequest) {
	for _, req := range requests {
		ls := m.getOrCreateLockState(req.Name)
		if req.Exclusive {
			ls.writerActive = false
		} else {
			ls.readersCount--
		}
	}
}

// ProcessLocalMultiLock represents a set of locks within a single process that
// are acquired and released together.
type ProcessLocalMultiLock struct {
	requests []LockRequest
	stack    string
	acquired bool
}

// NewProcessLocalMultiLock creates a new ProcessLocalMultiLock with the given
// lock requests. The locks will be sorted by name to ensure consistent
// ordering across all tests.
func NewProcessLocalMultiLock(t *testing.T, requests ...LockRequest) *ProcessLocalMultiLock {
	// Sort requests by name to ensure consistent ordering
	sortedRequests := make([]LockRequest, len(requests))
	copy(sortedRequests, requests)
	sort.Slice(sortedRequests, func(i, j int) bool {
		return sortedRequests[i].Name < sortedRequests[j].Name
	})

	// Collect caller information
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

	return &ProcessLocalMultiLock{
		requests: sortedRequests,
		stack:    stack,
		acquired: false,
	}
}

// AcquireAll acquires all locks atomically.
// It waits until all locks are available before acquiring any of them.
// This prevents lock chain slowdowns (Hold and Wait problem).
// The locks are automatically released when the test finishes via t.Cleanup.
//
// Note: ProcessLocalMultiLock does NOT use writer priority to avoid
// unnecessary blocking when waiting for multiple locks. This is a trade-off:
// we prevent lock chain slowdowns but may allow writer starvation in some
// scenarios.
func (ml *ProcessLocalMultiLock) AcquireAll(t *testing.T) {
	if ml.acquired {
		t.Fatalf("ProcessLocalMultiLock already acquired")
	}

	m := globalProcessLocalMultiLockManager

	// Log waiting
	ml.logWaiting()

	m.mu.Lock()
	defer m.mu.Unlock()

	for !m.canAcquireAll(ml.requests) {
		m.cond.Wait()
	}

	// Acquire all locks
	m.doAcquireAll(ml.requests)
	ml.acquired = true

	// Log acquired
	ml.logAcquired()

	// Register cleanup to release locks when test finishes
	t.Cleanup(func() { ml.ReleaseAll(t) })
}

// ReleaseAll releases all locks and wakes up waiting tests.
// This method is idempotent - calling it multiple times is safe.
func (ml *ProcessLocalMultiLock) ReleaseAll(t *testing.T) {
	if !ml.acquired {
		// Already released, do nothing (idempotent)
		return
	}

	m := globalProcessLocalMultiLockManager

	m.mu.Lock()
	defer m.mu.Unlock()

	// Release all locks
	m.doReleaseAll(ml.requests)
	ml.acquired = false

	// Log released
	ml.logReleased()

	// Wake up all waiting tests
	m.cond.Broadcast()
}

// logWaiting logs that we are waiting for locks.
func (ml *ProcessLocalMultiLock) logWaiting() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessLocalMultiLock[%s]: Wait for %s Lock(%s)", req.Name, lockType, ml.stack)
	}
}

// logAcquired logs that we have acquired the locks.
func (ml *ProcessLocalMultiLock) logAcquired() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessLocalMultiLock[%s]: %s Lock is acquired(%s)", req.Name, lockType, ml.stack)
	}
}

// logReleased logs that we have released the locks.
func (ml *ProcessLocalMultiLock) logReleased() {
	for _, req := range ml.requests {
		lockType := "Shared"
		if req.Exclusive {
			lockType = "Exclusive"
		}
		log.Printf("ProcessLocalMultiLock[%s]: %s Unlocked(%s)", req.Name, lockType, ml.stack)
	}
}
