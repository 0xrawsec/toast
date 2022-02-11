package toast

import (
	"fmt"
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

func Assert(t *testing.T, condition bool, i ...interface{}) {
	if !condition {
		t.Log(msg(assertFailMsg, i...))
		if FailFast {
			t.FailNow()
		}
	}
}

func ShouldPanic(t *testing.T, f func(), i ...interface{}) {
	defer func() { recover() }()
	f()
	t.Log(msg(assertFailMsg, i...))
	if FailFast {
		t.FailNow()
	}
}

func Error(t *testing.T, i ...interface{}) {
	t.Log(msg("", i...))
	if FailFast {
		t.FailNow()
	}
}
