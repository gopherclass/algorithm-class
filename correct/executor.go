package main

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strconv"
	"time"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

// type Module interface {
// 	Call(input string) (string, error)
// }

type Executor interface {
	Call(protocol Protocol, input string) (string, error)
}

type Protocol interface {
	ProtocolReader
	ProtocolWriter
}

type ProtocolReader interface {
	ReadFrom(r *bufio.Reader) (string, error)
}

type ProtocolWriter interface {
	WriteTo(w *bufio.Writer, input string) error
}

type Process interface {
	Start() error
	Stdin() io.Writer
	Stdout() io.Reader
	Stderr() io.Reader
	Close() error
}

type Command struct {
	started  bool
	startErr error
	process  Process
	stdin    *bufio.Writer
	stdout   *bufio.Reader
	stderr   *bufio.Reader
}

func NewCommand(ctx context.Context, name string, args ...string) *Command {
	return newCommand(spawnProcess(ctx, name, args...))
}

func newCommand(process Process) *Command {
	return &Command{process: process}
}

func (c *Command) Call(protocol Protocol, input string) (string, error) {
	if c.process == nil {
		return "", xerrors.Errorf("cannot call closed command")
	}
	if !c.started {
		c.startErr = c.start()
		c.started = true
	}
	if c.startErr != nil {
		return "", c.startErr
	}
	return c.call(protocol, input)
}

func (c *Command) call(protocol Protocol, input string) (string, error) {
	err := c.transmit(protocol, input)
	if err != nil {
		return "", err
	}
	return c.recv(protocol)
}

func (c *Command) transmit(w ProtocolWriter, input string) error {
	err := w.WriteTo(c.stdin, input)
	if err != nil {
		return err
	}
	err = c.stdin.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Command) recv(r ProtocolReader) (string, error) {
	return r.ReadFrom(c.stdout)
}

func (c *Command) start() error {
	if c.process.Stdin() != nil {
		c.stdin = bufio.NewWriter(c.process.Stdin())
	}
	if c.process.Stdout() != nil {
		c.stdout = bufio.NewReader(c.process.Stdout())
	}
	if c.process.Stderr() != nil {
		c.stderr = bufio.NewReader(c.process.Stderr())
	}
	return c.process.Start()
}

func (c *Command) Close() error {
	if c.process == nil {
		return xerrors.Errorf("closing of closed command")
	}
	err := c.process.Close()
	c.process = nil
	return err
}

type realProcess struct {
	command    *exec.Cmd
	initErr    error
	stdinPipe  io.Writer
	stdoutPipe io.Reader
	stderrPipe io.Reader
}

func spawnProcess(ctx context.Context, name string, args ...string) *realProcess {
	p := new(realProcess)
	p.command = exec.CommandContext(ctx, name, args...)
	p.initErr = p.init()
	return p
}

func (p *realProcess) init() error {
	err := p.initStdin()
	if err != nil {
		return err
	}
	err = p.initStdout()
	if err != nil {
		return err
	}
	err = p.initStderr()
	if err != nil {
		return err
	}
	return nil
}

func (p *realProcess) initStdin() (err error) {
	p.stdinPipe, err = p.command.StdinPipe()
	return err
}

func (p *realProcess) initStdout() (err error) {
	p.stdoutPipe, err = p.command.StdoutPipe()
	return err
}

func (p *realProcess) initStderr() (err error) {
	p.stderrPipe, err = p.command.StderrPipe()
	return err
}

func (p *realProcess) Start() error {
	if p.initErr != nil {
		return p.initErr
	}
	return p.command.Start()
}

func (p *realProcess) Close() error {
	if p.command.Process != nil {
		return p.command.Process.Kill()
	}
	return nil
}

func (p *realProcess) Stdin() io.Writer {
	return p.stdinPipe
}

func (p *realProcess) Stdout() io.Reader {
	return p.stdoutPipe
}

func (p *realProcess) Stderr() io.Reader {
	return p.stderrPipe
}

type Timer struct {
	ProtocolWriter
	ProtocolReader
	Timeout time.Duration
	timer   *time.Timer
}

func (t *Timer) WriteTo(w *bufio.Writer, input string) error {
	if t.timer == nil {
		t.timer = time.NewTimer(0)
	}
	if !t.timer.Stop() && len(t.timer.C) > 0 {
		<-t.timer.C
	}
	t.timer.Reset(t.Timeout)

	errc := make(chan error, 1)
	go func() {
		err := t.ProtocolWriter.WriteTo(w, input)
		errc <- err
	}()
	select {
	case err := <-errc:
		return err
	case <-t.timer.C:
		return xerrors.Errorf("module timeout")
	}
}

func (t *Timer) ReadFrom(r *bufio.Reader) (string, error) {
	type returnArgs struct {
		output string
		err    error
	}

	ret := make(chan returnArgs, 1)
	go func() {
		output, err := t.ProtocolReader.ReadFrom(r)
		ret <- returnArgs{output, err}
	}()
	select {
	case args := <-ret:
		return args.output, args.err
	case <-t.timer.C:
		return "", xerrors.Errorf("module timeout")
	}
}

type TextProtocol struct{}

func (TextProtocol) ReadFrom(r *bufio.Reader) (string, error) {
	rd := TextReader{Reader: r}
	buf, err := rd.BetweenFunc(isTextDelimiter)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (TextProtocol) WriteTo(w *bufio.Writer, input string) error {
	_, err := w.WriteString(input)
	return err
}

type PathProtocol struct {
	ProtocolWriter
	ProtocolReader
	Path string
}

func (p PathProtocol) ReadFrom(r *bufio.Reader) (string, error) {
	return p.ProtocolReader.ReadFrom(r)
}

func (p PathProtocol) WriteTo(w *bufio.Writer, input string) error {
	_, err := w.WriteString(p.Path)
	if err != nil {
		return err
	}
	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}
	return p.ProtocolWriter.WriteTo(w, input)
}

type TextReader struct {
	Reader *bufio.Reader
}

func (r *TextReader) BetweenFunc(isDelimiter func([]byte) bool) ([]byte, error) {
	var buf []byte
	for {
		line, err := r.ReadLine()
		if err != nil {
			return buf, err
		}
		if isDelimiter(line) {
			break
		}
		// 모든 CRLF는 LF로 변환됩니다.
		buf = append(buf, line...)
		buf = append(buf, '\n')
	}
	return buf, nil
}

func (r *TextReader) ReadLine() ([]byte, error) {
	line, err := r.ScanFunc(isLF)
	if err != nil {
		return nil, err
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line, nil
}

func (r *TextReader) ScanFunc(isInvalid func(rune) bool) ([]byte, error) {
	var buf []byte
	var enc [utf8.UTFMax]byte
	for {
		r0, _, err := r.Reader.ReadRune()
		if err == io.EOF && len(buf) > 0 {
			return buf, nil
		}
		if err != nil {
			return buf, err
		}
		if isInvalid(r0) {
			return buf, nil
		}
		size := utf8.EncodeRune(enc[:], r0)
		buf = append(buf, enc[:size]...)
	}
}

func (r *TextReader) Int() (int, error) {
	s, err := r.String()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (r *TextReader) Int64() (int64, error) {
	s, err := r.String()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (r *TextReader) String() (string, error) {
	buf, err := r.ScanFunc(isSeparatingCharacter)
	return string(buf), err
}

func isTextDelimiter(buf []byte) bool {
	return len(buf) >= 2 && buf[0] == '/' && buf[1] == '/'
}

func isLF(r rune) bool {
	return r == '\n'
}

func isSeparatingCharacter(r rune) bool {
	return r == ' ' || r == '\n'
}
