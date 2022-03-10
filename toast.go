package toast

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	FailFast = true
)

const (
	assertFailMsg = "assertion failed"
)

func mkfmt(i ...interface{}) string {
	fmt := make([]string, len(i))
	for i := range fmt {
		fmt[i] = "%v"
	}
	return strings.Join(fmt, " ")
}

func msg(or string, i ...interface{}) string {
	var prefix string
	_, filename, line, ok := runtime.Caller(2)

	if ok {
		prefix = fmt.Sprintf("\r    %s:%d:", filepath.Base(filename), line)
		i = append([]interface{}{prefix}, i...)
		if len(i) > 1 {
			return fmt.Sprintf(mkfmt(i...), i...)
		}
		return strings.Join([]string{prefix, or}, " ")
	}

	if len(i) > 0 {
		return fmt.Sprintf(mkfmt(i...), i...)
	}
	return or
}

func AssertOrPanic(condition bool, i ...interface{}) {
	if !condition {
		panic(msg(assertFailMsg, i...))
	}
}

type T struct {
	*testing.T
	FailFast bool
	mock     bool // to be able to test that structure without test failure
}

func FromT(t *testing.T) *T {
	return &T{t, FailFast, false}
}

func (t *T) log(s string) {
	f := t.T.Error
	if t.mock {
		f = t.T.Log
	}
	f(s)
}

func (t *T) Error(i ...interface{}) {
	t.log(msg("", i...))
	if t.FailFast {
		os.Exit(1)
	}
}

func (t *T) CheckErr(err error) {
	if err != nil {
		t.log(msg("", err))
		if t.FailFast {
			os.Exit(1)
		}
	}
}

func (t *T) ExpectErr(err, expect error) {
	if !errors.As(err, &expect) {
		t.log(msg("unexpected error", fmt.Errorf("expecting %v got %v", expect, err)))
		if t.FailFast {
			os.Exit(1)
		}
	}
}

func (t *T) ShouldPanic(f func(), i ...interface{}) {
	defer func() { recover() }()
	f()
	t.log(msg("should have panicked", i...))
	if t.FailFast {
		os.Exit(1)
	}
}

func (t *T) Wrap(init, test, cleanup func(*testing.T)) {
	if init != nil {
		init(t.T)
	}

	if cleanup != nil {
		defer func() { cleanup(t.T) }()
	}

	test(t.T)
}

func (t *T) Assert(condition bool, i ...interface{}) {
	if !condition {
		t.log(msg(assertFailMsg, i...))
		if t.FailFast {
			os.Exit(1)
		}
	}
}
