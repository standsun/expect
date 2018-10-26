package expect

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	ErrGroupEmpty     = errors.New("expect group is empty")
	ErrGroupListEmpty = errors.New("expect group list is empty")
	ErrCommandIllegal = errors.New("command is illegal")
)

type Expect struct {
	f *os.File
	l []*Group
}

func New(name string, arg ...string) (*Expect, error) {
	c := exec.Command(name, arg...)

	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	return &Expect{
		f: f,
	}, nil
}

func (e *Expect) AddGroup(g *Group) error {
	if g == nil {
		return ErrGroupEmpty
	}
	e.l = append(e.l, g)

	return nil
}

func (e *Expect) Run() error {
	if len(e.l) < 1 {
		return ErrGroupListEmpty
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, e.f); err != nil {
				log.Printf("error resizing f: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

	for _, g := range e.l {
		if err := e.expect(g); err != nil {
			return err
		}
	}

	e.interact()

	return nil
}

func (e *Expect) expect(g *Group) error {
	ch := make(chan error, 1)
	timeout := time.After(g.timeout)

	go func(e *Expect, g *Group) {
		for {
			output := make([]byte, 1024)
			num, err := e.f.Read(output)
			if err != nil {
				ch <- err
				return
			}
			output = output[0:num]

			if input := g.Search(string(output)); input != "" {
				input := []byte(input + "\r\n")

				e.f.Write(input)

				os.Stdout.Write(output)
				if g.show {
					os.Stdout.Write(input)
				} else {
					os.Stdout.Write([]byte("•••••••\r\n")) // 终端不显示输入，比如：密码
				}
				ch <- nil
				return
			}
		}
	}(e, g)

	select {
	case err := <-ch:
		return err
	case <-timeout:
		return errors.New("timeout")
	}
}

func (e *Expect) Close() {
	e.f.Close()
}

func (e *Expect) interact() {
	go func() {
		io.Copy(e.f, os.Stdin)
	}()
	io.Copy(os.Stdout, e.f)
}
