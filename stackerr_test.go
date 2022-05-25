package stackerr

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	t.Log(New("error"))
}

func TestErrorf(t *testing.T) {
	e1 := Errorf("error1")
	e2 := Errorf("error2: %w", e1)
	e3 := Errorf("error3: %w", e2)
	e4 := errors.New("error4")
	e5 := Errorf("error5: %w", e4)

	assert.ErrorIs(t, e2, e1)
	assert.ErrorIs(t, e3, e1)
	assert.ErrorIs(t, e3, e2)
	assert.ErrorIs(t, e5, e4)

	assert.DeepEqual(t, errors.Unwrap(e2), e1, cmpopts.EquateErrors())
	assert.DeepEqual(t, errors.Unwrap(e3), e2, cmpopts.EquateErrors())
	assert.DeepEqual(t, errors.Unwrap(e5), e4, cmpopts.EquateErrors())

	assert.Equal(t, countStackTraceBlock(e1), 1)
	assert.Equal(t, countStackTraceBlock(e2), 1)
	assert.Equal(t, countStackTraceBlock(e3), 1)
	assert.Equal(t, countStackTraceBlock(e4), 0)
	assert.Equal(t, countStackTraceBlock(e5), 1)

	t.Log(e1)
	t.Log(e2)
	t.Log(e3)
	t.Log(e4)
	t.Log(e5)
}

func countStackTraceBlock(err error) int {
	s := err.Error()
	splited := strings.Split(s, "goroutine")
	if len(splited) == 0 {
		return 0
	}
	return len(splited) - 1
}

func TestTrimStack(t *testing.T) {
	type in struct {
		buf  string
		skip int
	}

	for _, c := range []struct {
		in   *in
		want string
	}{
		{
			&in{"head\nbody", 0},
			"head\nbody",
		},
		{
			&in{"head\n", 1},
			"head",
		},
		{
			&in{"head\nfunc\npath", 1},
			"head",
		},
		{
			&in{"head\nfunc\npath\nfunc2\npath2", 1},
			"head\nfunc2\npath2",
		},
		{
			&in{"head\nfunc\npath\nfunc2\npath2\nfunc3\npath3", 2},
			"head\nfunc3\npath3",
		},
	} {
		got := string(trimStack([]byte(c.in.buf), c.in.skip))
		assert.DeepEqual(t, c.want, got)
	}
}
