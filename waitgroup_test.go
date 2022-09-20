// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debugsync

import (
	"sync/atomic"
	"testing"
)

func testWaitGroup(t *testing.T, wg1 *WaitGroup, wg2 *WaitGroup) {
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")
		default:
		}
		wg2.Done()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroupDebugOff(t *testing.T) {
	DebugIsOn = false

	wg1 := &WaitGroup{}
	wg2 := &WaitGroup{}

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func TestWaitGroupDebugOn(t *testing.T) {
	DebugIsOn = true

	wg1 := &WaitGroup{}
	wg2 := &WaitGroup{}

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func waitGroupMisuse(t *testing.T) {
	defer func() {
		err := recover()
		if err != "sync: negative WaitGroup counter" {
			t.Fatalf("Unexpected panic: %#v", err)
		}
	}()
	wg := &WaitGroup{}
	wg.Add(1)
	wg.Done()
	wg.Done()
	t.Fatal("Should panic")
}

func TestWaitGroupMisuseDebugOff(t *testing.T) {
	DebugIsOn = false
	waitGroupMisuse(t)
}

func TestWaitGroupMisuseDebugOn(t *testing.T) {
	DebugIsOn = true
	waitGroupMisuse(t)
}

func waitGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		wg := &WaitGroup{}
		n := new(int32)
		// spawn goroutine 1
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// spawn goroutine 2
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// Wait for goroutine 1 and 2
		wg.Wait()
		if atomic.LoadInt32(n) != 2 {
			t.Fatal("Spurious wakeup from Wait")
		}
	}
}

func TestWaitGroupRaceDebugOff(t *testing.T) {
	DebugIsOn = false
	waitGroupRace(t)
}

func TestWaitGroupRaceDebugOn(t *testing.T) {
	DebugIsOn = true
	waitGroupRace(t)
}

func waitGroupAlign(t *testing.T) {
	type X struct {
		wg WaitGroup
	}
	var x X
	x.wg.Add(1)
	go func(x *X) {
		x.wg.Done()
	}(&x)
	x.wg.Wait()
}

func TestWaitGroupAlignDebugOff(t *testing.T) {
	DebugIsOn = false
	waitGroupAlign(t)
}

func TestWaitGroupAlignDebugOn(t *testing.T) {
	DebugIsOn = true
	waitGroupAlign(t)
}
