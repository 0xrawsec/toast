package toast

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	FailNow = true

	format = fmt.Sprintf
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

type Toaster interface {
	FailNow()
	Log(args ...interface{})
	Error(args ...interface{})
}

type T struct {
	T       Toaster
	FailNow bool
	mock    bool // to be able to test that structure without test failure
}

func FromT(t *testing.T) *T {
	return &T{t, FailNow, false}
}

func FromB(t *testing.B) *T {
	return &T{t, FailNow, false}
}

func From(t Toaster) *T {
	return &T{t, FailNow, false}
}

func (t *T) log(s string) {
	t.T.Log(s)
}

func (t *T) logErr(s string) {
	f := t.T.Error
	if t.mock {
		f = t.T.Log
	}
	f(s)
}

func (t *T) Error(i ...interface{}) {
	t.logErr(msg("", i...))
	if t.FailNow {
		t.T.FailNow()
	}
}

func (t *T) CheckErr(err error) {
	if err != nil {
		t.logErr(msg("", err))
		if t.FailNow {
			t.T.FailNow()
		}
	}
}

func (t *T) ExpectErr(err, expect error) {
	if !errors.Is(err, expect) {
		t.logErr(msg("unexpected error", fmt.Errorf("expecting %v got %v", expect, err)))
		if t.FailNow {
			t.T.FailNow()
		}
	}
}

func (t *T) ShouldPanic(f func(), i ...interface{}) {
	defer func() { recover() }()
	f()
	t.logErr(msg("should have panicked", i...))
	if t.FailNow {
		t.T.FailNow()
	}
}

func (t *T) Wrap(init, test, cleanup func(t Toaster)) {
	if init != nil {
		init(t.T)
	}

	if cleanup != nil {
		defer func() { cleanup(t.T) }()
	}

	test(t.T)
}

func (t *T) TimeIt(name string, f func()) {
	timer := time.Now()
	f()
	t.log(msg("", "Time", format("%s:", name), time.Since(timer).String()))
}

func (t *T) Assert(condition bool, i ...interface{}) {
	if !condition {
		t.logErr(msg(assertFailMsg, i...))
		if t.FailNow {
			t.T.FailNow()
		}
	}
}
