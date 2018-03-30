package otto

import (
	"fmt"
	"testing"
)

type ErrNativeTester struct {
	wrapped error
}

func (ent *ErrNativeTester) Error() string {
	return ent.wrapped.Error()
}

func TestNativeErrorClass(t *testing.T) {
	tt(t, func() {
		vm := New()

		ctor := func(call FunctionCall) error {
			arg := call.ArgumentList[0].String()
			return &ErrNativeTester{fmt.Errorf(arg)}
		}

		var blessedCtor func(value interface{}) Value
		makeThisThing, err := vm.ToValue(func(call FunctionCall) Value {
			return blessedCtor("never_change")
		})
		is(err, nil)

		cls := vm.CreateNativeErrorClass(
			"TestError",
			ctor,
			func(error, string) {},
			func(error) string { return "" },
			[]Property{},
			[]Property{
				{
					Name:  "makeThisThing",
					Value: makeThisThing,
				},
			},
		)
		blessedCtor = cls.InstanceOf

		_, err = vm.Run(`TestError("yum")`)
		is(err, "!=", nil)

		vm.Set("TestError", cls.Function)

		ret, err := vm.Run(`TestError("yum").message`)
		is(err, nil)
		is(ret.value, "yum")

		ret, err = vm.Run(`TestError("yum").name`)
		is(err, nil)
		is(ret.value, "TestError")

		ret, err = vm.Run(`TestError("yum") instanceof Error`)
		is(err, nil)

		ret, err = vm.Run(`TestError("yum") instanceof TestError`)
		is(err, nil)
		is(ret.bool(), true)
	})
}
