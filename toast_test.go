package toast

import (
	"fmt"
	"testing"
)

func initTest() {
	fmt.Println("Doing some init job")
}

func cleanupTest() {
	fmt.Println("Doing some cleanup job")
}

func testWrap(t *testing.T) {
	t.Log("This is a wrapped test")
}

func TestToast(t *testing.T) {
	Assert(t, mkfmt(1, 2, 3) == "%v %v %v")
	AssertOrPanic(true)
	ShouldPanic(t, func() { AssertOrPanic(2 == 3) })
	Wrap(t, initTest, testWrap, cleanupTest)
	CheckErr(t, nil)
	CheckErr(t, fmt.Errorf("This is a random error"))
	t.Log("Ending test")
}
