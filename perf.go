package main

import (
	"fmt"
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

func main() {
	fmt.Println(processTime())
	time.Sleep(time.Second)
	fmt.Println(processTime())

	timer := time.NewTimer(2 * time.Second)
	fmt.Println(processTime())
For:
	for {
		select {
		case <-timer.C:
			break For
		default:
		}
	}
	fmt.Println(processTime())
}
