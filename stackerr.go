package stackerr

import (
	"bytes"
	"errors"
	"fmt"
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

	if !errors.As(e, &errStack) {
		err.stack = trimStack(debug.Stack(), 2)
	}

	return err
}

func (e *errorStack) Error() string {
	var (
		sb  strings.Builder
		msg string
	)

	switch {
	case e.msg != "":
		msg = e.msg
	case e.err != nil:
		msg = e.err.Error()
	}
	sb.WriteString(strings.TrimRight(msg, "\n"))

	if len(e.stack) > 0 {
		sb.WriteString("\n------- stack trace -------\n")
		sb.Write(e.stack)
	}

	return sb.String()
}

func (e *errorStack) Unwrap() error {
	return e.err
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
