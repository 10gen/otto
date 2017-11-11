package otto

import (
	"testing"
)

func TestCreateNativeFunction(t *testing.T) {
	tt(t, func() {
		vm := New()

		nativeFunc, err := vm.CreateNativeFunction("myFunc", "somerandomfile.js", 23, func(c FunctionCall) Value {
			return toValue_string(c.ArgumentList[0].String() + " world")
		})
		is(err, nil)

		nativeObj := nativeFunc._object().value.(_nativeFunctionObject)
		is(nativeObj.name, "myFunc")
		is(nativeObj.file, "somerandomfile.js")
		is(nativeObj.line, 23)

		is(vm.Set("callTheFunc", nativeFunc), nil)

		ret, err := vm.Run(`callTheFunc("hello")`)
		is(err, nil)
		is(ret, "hello world")

		_, err = vm.CreateNativeFunction("myFunc", "somerandomfile.js", 23, nil)
		is(err, "!=", nil)
	})
}
