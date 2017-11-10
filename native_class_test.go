package otto

import (
	"testing"
)

func TestNativeClass(t *testing.T) {
	tt(t, func() {
		vm := New()

		ctor := func(call FunctionCall) interface{} {
			call.This.Object().Set("woof", 2)
			exp, _ := call.ArgumentList[0].Export()
			return exp
		}
		hello, err := vm.ToValue(func(call FunctionCall) Value {
			return toValue_string("hello world")
		})
		is(err, nil)
		toString, err := vm.ToValue(func(call FunctionCall) Value {
			exp, _ := call.This.Export()
			return toValue_string(exp.(string))
		})
		is(err, nil)

		var blessedCtor func(value interface{}) Value
		makeThisThing, err := vm.ToValue(func(call FunctionCall) Value {
			return blessedCtor("never_change")
		})
		is(err, nil)

		cls := vm.CreateNativeClass(
			"Carrot",
			ctor,
			[]Property{
				{
					Name:  "doHello",
					Value: hello,
				},
				{
					Name:  "staticValue",
					Value: toValue_int32(42),
				},
				{
					Name:  "toString",
					Value: toString,
				},
			},
			[]Property{
				{
					Name:  "makeThisThing",
					Value: makeThisThing,
				},
			},
		)
		blessedCtor = cls.InstanceOf

		_, err = vm.Run(`Carrot("yum")`)
		is(err, "!=", nil)

		vm.Set("Carrot", cls.Function)
		ret, err := vm.Run(`Carrot("yum")`)
		is(err, nil)
		is(ret._object().value.(_goNativeValue).value, "yum")

		ret, err = vm.Run(`Carrot("yum").doHello()`)
		is(err, nil)
		is(ret.String(), "hello world")

		ret, err = vm.Run(`Carrot("yum").staticValue`)
		is(err, nil)
		is(ret.number().int64, 42)

		ret, err = vm.Run(`(new Carrot("yum")).woof`)
		is(err, nil)
		is(ret.number().int64, 2)

		ret, err = vm.Run(`Carrot("yum").woof`)
		is(err, nil)
		is(ret.number().int64, 0)

		ret, err = vm.Run(`woof`)
		is(err, nil)
		is(ret.number().int64, 2)

		ret, err = vm.Run(`Carrot.makeThisThing("yummy")`)
		is(err, nil)
		is(ret._object().value.(_goNativeValue).value, "never_change")

		ret, err = vm.Run(`Carrot("yum").toString()`)
		is(err, nil)
		is(ret.String(), "yum")

		ret, err = vm.Run(`Carrot("yum") instanceof Carrot`)
		is(err, nil)
		is(ret.bool(), true)

		ret, err = vm.Run(`5 instanceof Carrot`)
		is(err, nil)
		is(ret.bool(), false)
	})
}
