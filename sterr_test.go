package sterr

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func TestIs(t *testing.T) {
	testCases := []struct {
		desc string
		a, b error
		res  bool
	}{
		{
			desc: "same simple",
			a:    New("a"),
			b:    New("a"),
			res:  true,
		},
		{
			desc: "same complex",
			a:    New("a").Wrap(New("b")),
			b:    New("a").Wrap(New("b")),
			res:  true,
		},
		{
			desc: "different simple",
			a:    New("a"),
			b:    New("b"),
		},
		{
			desc: "different complex",
			a:    New("a").Wrap(New("c")),
			b:    New("a").Wrap(New("b")),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if errors.Is(tC.a, tC.b) != tC.res {
				t.Fail()
			}
		})
	}
}

func TestTrace(t *testing.T) {
	err := New("hello2").Trace(1)

	trc := ReadTrace(err)
	if len(trc) != 1 {
		t.Error(trc)
	}

}

func Test(t *testing.T) {
	test()
	t.Fail()
}

func test() {
	fmt.Println(runtime.Caller(0))
	fmt.Println(runtime.Caller(1))
}
