package stackerr_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"

	"github.com/mechiru/stackerr"
)

func TestNew(t *testing.T) {
	t.Log(stackerr.New("error"))
}

func TestErrorf(t *testing.T) {
	e1 := stackerr.Errorf("error1")
	e2 := stackerr.Errorf("error2: %w", e1)
	e3 := stackerr.Errorf("error3: %w", e2)
	e4 := errors.New("error4")
	e5 := stackerr.Errorf("error5: %w", e4)

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
	splited := strings.Split(s, "------- start stack trace -------")
	if len(splited) == 0 {
		return 0
	}
	return len(splited) - 1
}
