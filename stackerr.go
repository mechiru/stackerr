package stackerr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

var (
	errStack                             = &errorStack{}
	_        error                       = errStack
	_        interface{ Unwrap() error } = errStack
)

type errorStack struct {
	msg   string
	err   error
	stack []byte
}

func New(msg string) error {
	return &errorStack{
		msg:   msg,
		stack: trimStack(debug.Stack(), 2),
	}
}

func Errorf(format string, args ...any) error {
	var (
		e   = fmt.Errorf(format, args...)
		err = &errorStack{msg: e.(interface{ Error() string }).Error()}
	)

	if e, ok := e.(interface{ Unwrap() error }); ok {
		err.err = e.Unwrap()
	}

	if !errors.As(err.err, &errStack) {
		err.stack = trimStack(debug.Stack(), 2)
	}

	return err
}

func (e *errorStack) Error() string {
	var msg string
	switch {
	case e.msg != "":
		msg = e.msg
	case e.err != nil:
		msg = e.err.Error()
	}
	return strings.TrimRight(msg, "\n")
}

func (e *errorStack) Unwrap() error {
	return e.err
}

func (e *errorStack) Format(s fmt.State, v rune) {
	switch v {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Error()) //nolint:errcheck

			var (
				stack = e.stack
				err   = e.Unwrap()
			)
			for {
				if len(stack) > 0 {
					s.Write([]byte{'\n'})
					s.Write(bytes.TrimRight(stack, "\n"))
					return
				}

				if err == nil {
					return
				}

				if e, ok := err.(*errorStack); ok {
					stack = e.stack
					err = e.Unwrap()
				} else if e, ok := err.(interface{ Unwrap() error }); ok {
					err = e.Unwrap()
				} else {
					return
				}
			}
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error()) //nolint:errcheck
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func trimStack(buf []byte, skip int) []byte {
	if skip <= 0 {
		return buf
	}

	first := bytes.IndexByte(buf, '\n')
	if first < 0 {
		return buf
	}
	if len(buf)-1 == first {
		return buf[:first]
	}

	n := first + 1
	for i, skip := 0, skip*2; i < skip; i++ {
		idx := bytes.IndexByte(buf[n:], '\n')
		if idx < 0 {
			return buf[:first]
		}
		n += idx + 1
	}

	second := buf[n:]

	stack := make([]byte, 0, first+1+len(second))
	stack = append(stack, buf[:first+1]...)
	stack = append(stack, second...)
	return stack
}
