package toast

import (
	"fmt"
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
		t.Error(msg(assertFailMsg, i...))
		if FailFast {
			t.FailNow()
		}
	}
}

func ShouldPanic(t *testing.T, f func(), i ...interface{}) {
	defer func() { recover() }()
	f()
	t.Error(msg(assertFailMsg, i...))
	if FailFast {
		t.FailNow()
	}
}
