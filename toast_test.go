package toast

import (
	"fmt"
	"io"
	"testing"
)

var (
	i = 0
)

func initTest(t *testing.T) {
	if i != 0 {
		t.Error("Init failed")
	}
	t.Log("Doing some init job")
	i++
}

func cleanupTest(t *testing.T) {
	if i != 2 {
		t.Error("Cleanup test failed")
	}
	t.Log("Doing some cleanup job")
}

func testWrap(t *testing.T) {
	if i != 1 {
		t.Error("Wrapped test failed")
	}
	t.Log("This is a wrapped test")
	i++
}

func TestToast(t *testing.T) {
	tt := FromT(t)
	// making tests
	tt.FailNow = false
	tt.mock = true

	tt.ExpectErr(fmt.Errorf("random error"), nil)
	tt.ExpectErr(fmt.Errorf("encountered error %w", io.ErrClosedPipe), io.ErrClosedPipe)

	tt.Wrap(initTest, testWrap, cleanupTest)

	// ok
	tt.Assert(mkfmt(1, 2, 3) == "%v %v %v")
	// fail
	tt.Assert("one" == "two")

	tt.ShouldPanic(func() { AssertOrPanic(2 == 3) })
	tt.ShouldPanic(func() { AssertOrPanic(true) })

	// should not print message
	tt.CheckErr(nil)
	// should print message
	tt.CheckErr(fmt.Errorf("This is a random error"))
}
