// This file is adapted from the GO sync package.
// It originally contains the following license:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// GOMAXPROCS=10 go test

//lint:file-ignore SA2001 Empty critical section is acceptable in the context of a test

package sync

import (
	"fmt"

	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func HammerMutex(m *Mutex, loops int, cdone chan bool) {
	for i := 0; i < loops; i++ {
		if i%3 == 0 {
			if m.TryLock() {
				m.Unlock()
			}
			continue
		}
		m.Lock()
		m.Unlock()
	}
	cdone <- true
}

func mutex(t *testing.T) {
	DebugIsOn = false
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)

	m := new(Mutex)

	m.Lock()
	if m.TryLock() {
		t.Fatalf("TryLock succeeded with mutex locked")
	}
	m.Unlock()
	if !m.TryLock() {
		t.Fatalf("TryLock failed with mutex unlocked")
	}
	m.Unlock()

	c := make(chan bool)
	for i := 0; i < 10; i++ {
		go HammerMutex(m, 1000, c)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}

func TestMutexDebugOff(t *testing.T) {
	DebugIsOn = false
	mutex(t)
}

func TestMutexDebugOn(t *testing.T) {
	DebugIsOn = true
	mutex(t)
}

var misuseTests = []struct {
	name string
	f    func()
}{
	{
		"Mutex.Unlock",
		func() {
			var mu Mutex
			mu.Unlock()
		},
	},
	{
		"Mutex.Unlock2",
		func() {
			var mu Mutex
			mu.Lock()
			mu.Unlock()
			mu.Unlock()
		},
	},
	{
		"RWMutex.Unlock",
		func() {
			var mu RWMutex
			mu.Unlock()
		},
	},
	{
		"RWMutex.Unlock2",
		func() {
			var mu RWMutex
			mu.RLock()
			mu.Unlock()
		},
	},
	{
		"RWMutex.Unlock3",
		func() {
			var mu RWMutex
			mu.Lock()
			mu.Unlock()
			mu.Unlock()
		},
	},
	{
		"RWMutex.RUnlock",
		func() {
			var mu RWMutex
			mu.RUnlock()
		},
	},
	{
		"RWMutex.RUnlock2",
		func() {
			var mu RWMutex
			mu.Lock()
			mu.RUnlock()
		},
	},
	{
		"RWMutex.RUnlock3",
		func() {
			var mu RWMutex
			mu.RLock()
			mu.RUnlock()
			mu.RUnlock()
		},
	},
}

func init() {
	if len(os.Args) == 3 && os.Args[1] == "TESTMISUSE" {
		for _, test := range misuseTests {
			if test.name == os.Args[2] {
				func() {
					defer func() { recover() }()
					test.f()
				}()
				fmt.Printf("test completed\n")
				os.Exit(0)
			}
		}
		fmt.Printf("unknown test\n")
		os.Exit(0)
	}
}

func mutexMisuse(t *testing.T) {
	for _, test := range misuseTests {
		out, err := exec.Command(os.Args[0], "TESTMISUSE", test.name).CombinedOutput()
		if err == nil || !strings.Contains(string(out), "unlocked") {
			t.Errorf("%s: did not find failure with message about unlocked lock: %s\n%s\n", test.name, err, out)
		}
	}
}

func TestMutexMisuseDebugOff(t *testing.T) {
	DebugIsOn = false
	mutexMisuse(t)
}

func TestMutexMisuseDebugOn(t *testing.T) {
	DebugIsOn = true
	mutexMisuse(t)
}

func mutexFairness(t *testing.T) {
	var mu Mutex
	stop := make(chan bool)
	defer close(stop)
	go func() {
		for {
			mu.Lock()
			time.Sleep(100 * time.Microsecond)
			mu.Unlock()
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
	done := make(chan bool, 1)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Microsecond)
			mu.Lock()
			mu.Unlock()
		}
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatalf("can't acquire Mutex in 10 seconds")
	}
}

func TestMutexFairnessDebugOff(t *testing.T) {
	DebugIsOn = false
	mutexFairness(t)
}

func TestMutexFairnessDebugOn(t *testing.T) {
	DebugIsOn = true
	mutexFairness(t)
}
