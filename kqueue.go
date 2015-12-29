// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package kqueue

import (
	"os"
	"syscall"
)

type Kqueue struct {
	Fd int
}

func NewKqueue() (*Kqueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, os.NewSyscallError("kqueue", err)
	}
	return &Kqueue{
		Fd: fd,
	}, nil
}

func (k *Kqueue) Wait(buf []syscall.Kevent_t) (int, error) {
	n, err := syscall.Kevent(k.Fd, nil, buf, nil)
	if err != nil {
		return 0, os.NewSyscallError("kevent", err)
	}
	return n, nil
}

func (k *Kqueue) Close() error {
	return syscall.Close(k.Fd)
}
