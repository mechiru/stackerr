package stackerr

import (
	"errors"
	"fmt"
	"runtime"
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
	// Copy from debug.Stack()
	stack := make([]byte, 1024)
	for {
		n := runtime.Stack(stack, false)
		if n < len(stack) {
			stack = stack[:n]
			break
		}
		stack = make([]byte, 2*len(stack))
	}

	return &errorStack{
		msg:   msg,
		stack: stack,
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

	if !errors.As(e, &errStack) {
		// Copy from debug.Stack()
		buf := make([]byte, 1024)
		for {
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				buf = buf[:n]
				break
			}
			buf = make([]byte, 2*len(buf))
		}

		err.stack = buf
	}

	return err
}

func (e *errorStack) Error() string {
	var sb strings.Builder

	switch {
	case e.msg != "":
		sb.WriteString(e.msg)
	case e.err != nil:
		sb.WriteString(e.err.Error())
	}

	if len(e.stack) > 0 {
		sb.WriteString("\n------- start stack trace -------\n")
		sb.Write(e.stack)
		sb.WriteString("------- end stack trace -------")
	}

	return sb.String()
}

func (e *errorStack) Unwrap() error {
	return e.err
}
