// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package kqueue

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func newfile(t *testing.T, dir string) *os.File {
	f, err := ioutil.TempFile(dir, "file")
	assert(t, err)
	return f
}

// setup creates a temporary working directory, and a cleanup function.
func setup(t *testing.T) (string, func()) {
	d, err := ioutil.TempDir("", "fswatch-")
	if err != nil {
		t.Fatal(err)
	}
	return d, func() {
		os.RemoveAll(d)
	}
}

func TestNewKqueue(t *testing.T) {
	kq, err := NewKqueue()
	assert(t, err)
	assert(t, kq.Close())
}

func TestKqueueAddSocket(t *testing.T) {
	kq, err := NewKqueue()
	assert(t, err)
	s, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	assert(t, err)
	defer syscall.Close(s)
	assert(t, kq.Add(uintptr(s), syscall.EVFILT_WRITE, 0, 0))
	assert(t, kq.Close())
}

func TestKqueueAddFile(t *testing.T) {
	d, done := setup(t)
	defer done()
	kq, err := NewKqueue()
	assert(t, err)
	f := newfile(t, d)
	defer f.Close()
	assert(t, kq.Add(f.Fd(), syscall.EVFILT_WRITE, 0, 0))
	assert(t, kq.Close())
}

func TestKqueueWaitPipeClose(t *testing.T) {
	kq, err := NewKqueue()
	assert(t, err)
	pr, pw, err := os.Pipe()
	assert(t, err)
	defer pw.Close()
	assert(t, kq.Add(uintptr(pr.Fd()), syscall.EVFILT_READ, syscall.EV_ONESHOT, 0))
	result := make(chan error)
	go func() {
		_, err := kq.Wait(make([]syscall.Kevent_t, 1))
		result <- err
	}()
	assert(t, pw.Close())
	assert(t, <-result)
	assert(t, kq.Close())
}
