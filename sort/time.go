package main

import (
	"syscall"
	"time"
)

func processTime() (time.Duration, error) {
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	var r syscall.Rusage
	err = syscall.GetProcessTimes(process,
		&r.CreationTime,
		&r.ExitTime,
		&r.KernelTime,
		&r.UserTime,
	)
	if err != nil {
		return 0, err
	}
	convInt64 := func(time syscall.Filetime) int64 {
		return int64(time.HighDateTime)<<32 + int64(time.LowDateTime)
	}
	nsec := convInt64(r.KernelTime) + convInt64(r.UserTime)
	nsec *= 100
	return time.Duration(nsec), nil
}

func mustProcessTime() time.Duration {
	t, err := processTime()
	if err != nil {
		panic(err)
	}
	return t
}

type timer struct {
	start time.Duration
}

func newTimer() *timer {
	return &timer{start: mustProcessTime()}
}

func (t *timer) stop() time.Duration {
	now := mustProcessTime()
	return now - t.start
}
