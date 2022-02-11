package toast

import "testing"

func TestToast(t *testing.T) {
	Assert(t, mkfmt(1, 2, 3) == "%v %v %v")
	AssertOrPanic(true)
	ShouldPanic(t, func() { AssertOrPanic(2 == 3) })
}
