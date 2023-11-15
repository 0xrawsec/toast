package toast

import (
	"errors"
	"fmt"
	"os"
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

	var filename string
	var line int
	var ok bool

	for skip := 2; skip < 8; skip++ {
		_, filename, line, ok = runtime.Caller(skip)
		if strings.HasSuffix(filepath.Base(filename), "_test.go") {
			break
		}
	}

	if ok {
		prefix = fmt.Sprintf("    %s:%d:", filename, line)
		//prefix = fmt.Sprintf("    %s:%d:", filepath.Base(filename), line)
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
	Cleanup(f func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Setenv(key, value string)
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

type T struct {
	Toaster
	FailNowFl bool
	mock      bool // to be able to test that structure without test failure
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

func (t *T) logErr(s string) {
	if t.mock {
		t.Log(s)
	} else {
		t.Error(s)
	}
}

func (t *T) Log(args ...interface{}) {
	fmt.Fprintln(os.Stderr, msg("", args...))
}

func (t *T) Logf(f string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, msg("", fmt.Sprintf(f, args...)))
}

func (t *T) Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, msg("", args...))
	t.Fail()
}

func (t *T) Errorf(f string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, msg("", fmt.Sprintf(f, args...)))
	t.Fail()
}

func (t *T) CheckErr(err error) {
	if err != nil {
		t.logErr(fmt.Sprint(err))
		if t.FailNowFl {
			t.FailNow()
		}
	}
}

func (t *T) ExpectErr(err, expect error) {
	if !errors.Is(err, expect) {
		t.logErr(fmt.Sprintf("expecting %v got %v", expect, err))
		if t.FailNowFl {
			t.FailNow()
		}
	}
}

func (t *T) ShouldPanic(f func(), i ...interface{}) {
	defer func() { recover() }()
	f()
	t.logErr(fmt.Sprint(i...))
	if t.FailNowFl {
		t.FailNow()
	}
}

func (t *T) Wrap(init, test, cleanup func(t Toaster)) {
	if init != nil {
		init(t)
	}

	if cleanup != nil {
		defer func() { cleanup(t) }()
	}

	test(t)
}

func (t *T) TimeIt(name string, f func()) {
	timer := time.Now()
	f()
	t.Log("Time", format("%s:", name), time.Since(timer).String())
}

func (t *T) Assert(condition bool, i ...interface{}) {
	if !condition {
		t.logErr(assertFailMsg)
		if t.FailNowFl {
			t.FailNow()
		}
	}
}
