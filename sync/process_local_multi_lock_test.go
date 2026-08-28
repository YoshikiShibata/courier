// Copyright © 2026 Yoshiki Shibata. All rights reserved.

package sync

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessLocalMultiLock_BasicAcquireRelease(t *testing.T) {
	ml := NewProcessLocalMultiLock(t,
		Shared("lock1"),
		Exclusive("lock2"),
	)

	ml.AcquireAll(t)
	ml.ReleaseAll(t)
}

func TestProcessLocalMultiLock_MultipleSharedLocks(t *testing.T) {
	// Multiple goroutines should be able to acquire shared locks simultaneously
	var wg sync.WaitGroup
	var acquired int32

	for range 5 {
		wg.Go(func() {
			ml := NewProcessLocalMultiLock(t, Shared("shared_lock"))
			ml.AcquireAll(t)
			atomic.AddInt32(&acquired, 1)
			time.Sleep(50 * time.Millisecond)
			ml.ReleaseAll(t)
		})
	}

	// Wait a bit for all goroutines to acquire the lock
	time.Sleep(100 * time.Millisecond)

	// All 5 should have acquired the shared lock
	if atomic.LoadInt32(&acquired) != 5 {
		t.Fatalf("Expected 5 goroutines to acquire shared lock, got %d", acquired)
	}

	wg.Wait()
}

func TestProcessLocalMultiLock_ExclusiveLockBlocksOthers(t *testing.T) {
	var wg sync.WaitGroup
	var order []int
	var mu sync.Mutex

	appendOrder := func(n int) {
		mu.Lock()
		order = append(order, n)
		mu.Unlock()
	}

	// First goroutine acquires exclusive lock
	wg.Go(func() {
		ml := NewProcessLocalMultiLock(t, Exclusive("exclusive_test_lock"))
		ml.AcquireAll(t)
		appendOrder(1)
		time.Sleep(100 * time.Millisecond)
		ml.ReleaseAll(t)
	})

	// Wait for first goroutine to acquire
	time.Sleep(20 * time.Millisecond)

	// Second goroutine tries to acquire (should wait)
	wg.Go(func() {
		ml := NewProcessLocalMultiLock(t, Exclusive("exclusive_test_lock"))
		ml.AcquireAll(t)
		appendOrder(2)
		ml.ReleaseAll(t)
	})

	wg.Wait()

	// Verify order
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("Expected order [1, 2], got %v", order)
	}
}

func TestProcessLocalMultiLock_AtomicAcquisition(t *testing.T) {
	// This test verifies that ProcessLocalMultiLock acquires all locks atomically
	// and doesn't hold some locks while waiting for others

	var wg sync.WaitGroup
	var chainDetected atomic.Bool

	// Goroutine 1: Holds lock_a, lock_b
	wg.Go(func() {
		ml := NewProcessLocalMultiLock(t,
			Exclusive("atomic_lock_a"),
			Exclusive("atomic_lock_b"),
		)
		ml.AcquireAll(t)
		time.Sleep(100 * time.Millisecond)
		ml.ReleaseAll(t)
	})

	// Wait for goroutine 1 to acquire locks
	time.Sleep(20 * time.Millisecond)

	// Goroutine 2: Needs lock_b, lock_c
	// With the old implementation, this could hold lock_c while waiting for lock_b
	// causing a chain. With ProcessLocalMultiLock, it waits for both without holding either.
	wg.Go(func() {
		ml := NewProcessLocalMultiLock(t,
			Exclusive("atomic_lock_b"),
			Exclusive("atomic_lock_c"),
		)
		ml.AcquireAll(t)
		ml.ReleaseAll(t)
	})

	// Goroutine 3: Needs lock_c only
	// If goroutine 2 held lock_c while waiting, this would be blocked
	// With ProcessLocalMultiLock, this should acquire quickly
	wg.Go(func() {
		time.Sleep(30 * time.Millisecond) // Start after goroutine 2

		start := time.Now()
		ml := NewProcessLocalMultiLock(t, Exclusive("atomic_lock_c"))
		ml.AcquireAll(t)
		elapsed := time.Since(start)
		ml.ReleaseAll(t)

		// If goroutine 2 held lock_c while waiting for lock_b,
		// this would take ~70ms (100ms - 30ms). With atomic acquisition,
		// it should be much faster.
		if elapsed > 50*time.Millisecond {
			chainDetected.Store(true)
		}
	})

	wg.Wait()

	if chainDetected.Load() {
		t.Error("Lock chain detected - atomic acquisition may not be working correctly")
	}
}

func TestProcessLocalMultiLock_SortedByName(t *testing.T) {
	// Verify that locks are sorted by name regardless of input order
	ml1 := NewProcessLocalMultiLock(t,
		Exclusive("z_lock"),
		Shared("a_lock"),
		Exclusive("m_lock"),
	)

	// Check that requests are sorted
	if ml1.requests[0].Name != "a_lock" {
		t.Fatalf("Expected first lock to be 'a_lock', got '%s'", ml1.requests[0].Name)
	}
	if ml1.requests[1].Name != "m_lock" {
		t.Fatalf("Expected second lock to be 'm_lock', got '%s'", ml1.requests[1].Name)
	}
	if ml1.requests[2].Name != "z_lock" {
		t.Fatalf("Expected third lock to be 'z_lock', got '%s'", ml1.requests[2].Name)
	}
}
