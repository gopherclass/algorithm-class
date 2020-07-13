package main

import (
	"bufio"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestCall(t *testing.T) {
	process := newInternalProcess(internalFunc(echoModule))
	cmd := newCommand(process)

	s, err := cmd.Call(TextProtocol{}, "hello\n")
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello\n" {
		t.Fatalf("wrong output: %q != hello", s)
	}
}

func echoModule(in io.Reader, w io.Writer) error {
	br := bufio.NewReader(in)
	r := TextReader{Reader: br}
	_, lookaheadErr := br.Peek(0)
	for lookaheadErr == nil {
		s, err := r.String()
		if err != nil {
			return err
		}
		fmt.Fprintln(w, s)
		fmt.Fprintln(w, "//")
	}
	if lookaheadErr != io.EOF {
		return lookaheadErr
	}
	return nil
}

func TestTimer(t *testing.T) {
	process := newInternalProcess(internalFunc(echoModule))
	cmd := newCommand(process)

	s, err := cmd.Call(&Timer{
		ProtocolWriter: TextProtocol{},
		ProtocolReader: TextProtocol{},
		Timeout:        time.Second,
	}, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello\n" {
		t.Fatalf("wrong output: %q != hello", s)
	}
}

type internalCode interface {
	Main(r io.Reader, w io.Writer) error
}

type internalFunc func(io.Reader, io.Writer) error

func (fn internalFunc) Main(r io.Reader, w io.Writer) error {
	return fn(r, w)
}

type internalProcess struct {
	code internalCode
	w    *io.PipeWriter
	r    *io.PipeReader
}

func newInternalProcess(code internalCode) *internalProcess {
	return &internalProcess{code: code}
}

var _ Process = new(internalProcess)

func (i *internalProcess) Start() error {
	ir, iw := io.Pipe()
	or, ow := io.Pipe()
	i.w = iw
	i.r = or
	go func() {
		// TODO: Exit code도 받고 panic도 처리하기
		i.code.Main(ir, ow)
	}()
	return nil
}

func (i *internalProcess) Close() error {
	i.w.CloseWithError(io.EOF)
	i.r.Close()
	return nil
}

func (i *internalProcess) Read(p []byte) (int, error) {
	return i.r.Read(p)
}

func (i *internalProcess) Write(p []byte) (int, error) {
	return i.w.Write(p)
}

func (i *internalProcess) Stdin() io.Writer {
	return i
}

func (i *internalProcess) Stdout() io.Reader {
	return i
}

func (i *internalProcess) Stderr() io.Reader {
	return nil
}

type exitCode int

func (code exitCode) Error() string {
	return fmt.Sprintf("exit code = %d", int(code))
}
