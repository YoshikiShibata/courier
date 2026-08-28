// Copyright © 2026 Yoshiki Shibata. All rights reserved.

package sync

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestProcessSharedMultiLock_BasicAcquireRelease(t *testing.T) {
	ml := NewProcessSharedMultiLock(t,
		Shared("psml_basic_1"),
		Exclusive("psml_basic_2"),
	)

	ml.AcquireAll(t)
	ml.ReleaseAll(t)
}

func TestProcessSharedMultiLock_MultipleSharedLocks(t *testing.T) {
	// Multiple goroutines with independent file descriptors should be
	// able to hold shared locks on the same name simultaneously.
	var wg sync.WaitGroup
	var acquired int32

	for range 5 {
		wg.Go(func() {
			ml := NewProcessSharedMultiLock(t, Shared("psml_shared"))
			ml.AcquireAll(t)
			atomic.AddInt32(&acquired, 1)
			time.Sleep(50 * time.Millisecond)
			ml.ReleaseAll(t)
		})
	}

	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&acquired) != 5 {
		t.Fatalf("Expected 5 goroutines to acquire shared lock, got %d", acquired)
	}

	wg.Wait()
}

func TestProcessSharedMultiLock_ExclusiveLockBlocksOthers(t *testing.T) {
	var wg sync.WaitGroup
	var order []int
	var mu sync.Mutex

	appendOrder := func(n int) {
		mu.Lock()
		order = append(order, n)
		mu.Unlock()
	}

	wg.Go(func() {
		ml := NewProcessSharedMultiLock(t, Exclusive("psml_exclusive_block"))
		ml.AcquireAll(t)
		appendOrder(1)
		time.Sleep(100 * time.Millisecond)
		ml.ReleaseAll(t)
	})

	time.Sleep(20 * time.Millisecond)

	wg.Go(func() {
		ml := NewProcessSharedMultiLock(t, Exclusive("psml_exclusive_block"))
		ml.AcquireAll(t)
		appendOrder(2)
		ml.ReleaseAll(t)
	})

	wg.Wait()

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("Expected order [1, 2], got %v", order)
	}
}

func TestProcessSharedMultiLock_AtomicAcquisition(t *testing.T) {
	// Verifies that ProcessSharedMultiLock does not hold a subset of
	// requested locks while waiting for the rest.

	var wg sync.WaitGroup
	var chainDetected atomic.Bool

	// Goroutine 1 holds {a, b} for 100ms.
	wg.Go(func() {
		ml := NewProcessSharedMultiLock(t,
			Exclusive("psml_atomic_a"),
			Exclusive("psml_atomic_b"),
		)
		ml.AcquireAll(t)
		time.Sleep(100 * time.Millisecond)
		ml.ReleaseAll(t)
	})

	time.Sleep(20 * time.Millisecond)

	// Goroutine 2 wants {b, c}. If we had a hold-and-wait design it
	// would grab c and block waiting for b, keeping c unavailable to
	// goroutine 3.
	wg.Go(func() {
		ml := NewProcessSharedMultiLock(t,
			Exclusive("psml_atomic_b"),
			Exclusive("psml_atomic_c"),
		)
		ml.AcquireAll(t)
		ml.ReleaseAll(t)
	})

	// Goroutine 3 only wants c and should get it quickly, since
	// goroutine 2 must not hold c while waiting for b.
	wg.Go(func() {
		time.Sleep(30 * time.Millisecond)

		start := time.Now()
		ml := NewProcessSharedMultiLock(t, Exclusive("psml_atomic_c"))
		ml.AcquireAll(t)
		elapsed := time.Since(start)
		ml.ReleaseAll(t)

		// Allow generous slack: 60ms threshold vs the 70ms wait a
		// hold-and-wait design would impose, plus retry-interval
		// jitter (10ms).
		if elapsed > 60*time.Millisecond {
			chainDetected.Store(true)
		}
	})

	wg.Wait()

	if chainDetected.Load() {
		t.Error("Lock chain detected - atomic acquisition may not be working correctly")
	}
}

func TestProcessSharedMultiLock_SortedByName(t *testing.T) {
	ml := NewProcessSharedMultiLock(t,
		Exclusive("psml_sort_z"),
		Shared("psml_sort_a"),
		Exclusive("psml_sort_m"),
	)

	if ml.requests[0].Name != "psml_sort_a" {
		t.Fatalf("Expected first lock to be 'psml_sort_a', got '%s'", ml.requests[0].Name)
	}
	if ml.requests[1].Name != "psml_sort_m" {
		t.Fatalf("Expected second lock to be 'psml_sort_m', got '%s'", ml.requests[1].Name)
	}
	if ml.requests[2].Name != "psml_sort_z" {
		t.Fatalf("Expected third lock to be 'psml_sort_z', got '%s'", ml.requests[2].Name)
	}
}

// TestProcessSharedMultiLock_CrossProcess verifies that ProcessSharedMultiLock
// serializes correctly across separate OS processes and that the waiter
// truly sleeps in the kernel (event-driven) rather than polling.
//
// The test re-execs itself: the child acquires a well-known lock, writes
// "acquired" to stdout, holds the lock for holdDuration, then exits. The
// parent waits for the ack, times its own AcquireAll of the same lock,
// and checks two invariants:
//
//   - elapsed ≈ holdDuration (parent really did wait for the child)
//   - CPU used during the wait is tiny compared to the elapsed time
//     (regression guard: a busy-poll implementation would burn CPU)
const (
	psmlCrossProcLockName     = "psml_crossproc"
	psmlCrossProcHelperEnvKey = "COURIER_PSML_CROSSPROC_HELPER"
	psmlCrossProcHoldDuration = 20 * time.Second
)

func TestProcessSharedMultiLock_CrossProcess(t *testing.T) {
	if os.Getenv(psmlCrossProcHelperEnvKey) == "1" {
		// Child mode: acquire, notify parent, hold, return. The lock is
		// released by the t.Cleanup registered inside AcquireAll.
		ml := NewProcessSharedMultiLock(t, Exclusive(psmlCrossProcLockName))
		ml.AcquireAll(t)
		fmt.Println("acquired")
		_ = os.Stdout.Sync()
		time.Sleep(psmlCrossProcHoldDuration)
		return
	}

	if testing.Short() {
		t.Skip("skipping cross-process test in -short mode (takes ~20s)")
	}

	cmd := exec.Command(os.Args[0],
		"-test.run", "^TestProcessSharedMultiLock_CrossProcess$",
		"-test.timeout", "1m",
	)
	cmd.Env = append(os.Environ(), psmlCrossProcHelperEnvKey+"=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("StdoutPipe failed: %v", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("cmd.Start failed: %v", err)
	}

	// Wait for child to signal it has acquired the lock.
	type lineResult struct {
		line string
		err  error
	}
	lineC := make(chan lineResult, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if scanner.Text() == "acquired" {
				lineC <- lineResult{line: scanner.Text()}
				return
			}
		}
		lineC <- lineResult{err: fmt.Errorf("child stdout ended without 'acquired' (scanner err: %v)", scanner.Err())}
	}()

	select {
	case r := <-lineC:
		if r.err != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			t.Fatalf("waiting for child to acquire: %v", r.err)
		}
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatal("timed out waiting for child to signal 'acquired'")
	}

	// Drain remaining stdout so the child does not block on a full pipe.
	go func() { _, _ = io.Copy(io.Discard, stdout) }()

	// Measure elapsed and CPU while blocking on the same lock.
	var startRU, endRU syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &startRU); err != nil {
		t.Fatalf("Getrusage failed: %v", err)
	}
	start := time.Now()

	ml := NewProcessSharedMultiLock(t, Exclusive(psmlCrossProcLockName))
	ml.AcquireAll(t)
	elapsed := time.Since(start)
	ml.ReleaseAll(t)

	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &endRU); err != nil {
		t.Fatalf("Getrusage failed: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("child exited with error: %v", err)
	}

	// Elapsed should be close to holdDuration. Allow generous slack
	// on both sides for scheduling on loaded CI.
	minAcceptable := psmlCrossProcHoldDuration - 5*time.Second
	maxAcceptable := psmlCrossProcHoldDuration + 10*time.Second
	if elapsed < minAcceptable {
		t.Fatalf("parent acquired too early: elapsed=%v, expected around %v",
			elapsed, psmlCrossProcHoldDuration)
	}
	if elapsed > maxAcceptable {
		t.Fatalf("parent acquired too late: elapsed=%v, expected around %v",
			elapsed, psmlCrossProcHoldDuration)
	}

	// Busy-poll regression guard. A true event-driven wait burns
	// essentially no CPU. The former 10ms-sleep implementation would
	// burn tens of milliseconds; a wait with no sleep at all would
	// burn seconds. 1s is a generous threshold that catches the
	// egregious cases while staying flake-free.
	utime := time.Duration(endRU.Utime.Nano() - startRU.Utime.Nano())
	stime := time.Duration(endRU.Stime.Nano() - startRU.Stime.Nano())
	cpuUsed := utime + stime
	if cpuUsed > 1*time.Second {
		t.Fatalf("suspected busy polling: CPU=%v during a %v wait", cpuUsed, elapsed)
	}
	t.Logf("elapsed=%v, cpu=%v (utime=%v, stime=%v)",
		elapsed, cpuUsed, utime, stime)
}
