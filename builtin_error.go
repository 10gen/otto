package otto

import (
	"fmt"
)

func builtinError(call FunctionCall) Value {
	return toValue_object(call.runtime.newError("Error", call.Argument(0), nil, 1))
}

func builtinNewError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newError("Error", valueOfArrayIndex(argumentList, 0), nil, 0))
}

func builtinError_toString(call FunctionCall) Value {
	thisObject := call.thisObject()
	if thisObject == nil {
		panic(call.runtime.panicTypeError())
	}

	name := "Error"
	nameValue := thisObject.get("name")
	if nameValue.IsDefined() {
		name = nameValue.string()
	}

	message := ""
	messageValue := thisObject.get("message")
	if messageValue.IsDefined() {
		message = messageValue.string()
	}

	if len(name) == 0 {
		return toValue_string(message)
	}

	if len(message) == 0 {
		return toValue_string(name)
	}

	return toValue_string(fmt.Sprintf("%s: %s", name, message))
}

func (runtime *_runtime) getNativeErrorFunction(className string) *_object {
	proto := &_object{
		runtime:     runtime,
		class:       className,
		objectClass: _classObject,
		prototype:   runtime.global.ErrorPrototype,
		extensible:  true,
		value:       nil,
		property: map[string]_property{
			"name": _property{
				mode: 0101,
				value: Value{
					kind:  valueString,
					value: className,
				},
			},
		},
		propertyOrder: []string{
			"name",
		},
	}

	errFunc := &_object{
		runtime:     runtime,
		class:       "Function",
		objectClass: _classObject,
		prototype:   runtime.global.FunctionPrototype,
		extensible:  true,
		value: _nativeFunctionObject{
			name: className,
			call: func(call FunctionCall) Value {
				panic(&_exception{
					value: newError(runtime, "Error", 0, nil, "cannot call "+className),
				})
			},
			construct: func(self *_object, argumentList []Value) Value {
				panic(&_exception{
					value: newError(runtime, "Error", 0, nil, "cannot construct a "+className),
				})
			},
		},
		property: map[string]_property{
			"length": _property{
				mode: 0,
				value: Value{
					kind:  valueNumber,
					value: 1,
				},
			},
			"prototype": _property{
				mode: 0,
				value: Value{
					kind:  valueObject,
					value: proto,
				},
			},
		},
		propertyOrder: []string{
			"length",
			"prototype",
		},
	}

	return errFunc
}

func (runtime *_runtime) newEvalError(message Value) *_object {
	self := runtime.newErrorObject("EvalError", message, nil, 0)
	self.prototype = runtime.global.EvalErrorPrototype
	return self
}

func builtinEvalError(call FunctionCall) Value {
	return toValue_object(call.runtime.newEvalError(call.Argument(0)))
}

func builtinNewEvalError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newEvalError(valueOfArrayIndex(argumentList, 0)))
}

func (runtime *_runtime) newTypeError(message Value) *_object {
	self := runtime.newErrorObject("TypeError", message, nil, 0)
	self.prototype = runtime.global.TypeErrorPrototype
	return self
}

func builtinTypeError(call FunctionCall) Value {
	return toValue_object(call.runtime.newTypeError(call.Argument(0)))
}

func builtinNewTypeError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newTypeError(valueOfArrayIndex(argumentList, 0)))
}

func (runtime *_runtime) newRangeError(message Value) *_object {
	self := runtime.newErrorObject("RangeError", message, nil, 0)
	self.prototype = runtime.global.RangeErrorPrototype
	return self
}

func builtinRangeError(call FunctionCall) Value {
	return toValue_object(call.runtime.newRangeError(call.Argument(0)))
}

func builtinNewRangeError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newRangeError(valueOfArrayIndex(argumentList, 0)))
}

func (runtime *_runtime) newURIError(message Value) *_object {
	self := runtime.newErrorObject("URIError", message, nil, 0)
	self.prototype = runtime.global.URIErrorPrototype
	return self
}

func (runtime *_runtime) newReferenceError(message Value) *_object {
	self := runtime.newErrorObject("ReferenceError", message, nil, 0)
	self.prototype = runtime.global.ReferenceErrorPrototype
	return self
}

func builtinReferenceError(call FunctionCall) Value {
	return toValue_object(call.runtime.newReferenceError(call.Argument(0)))
}

func builtinNewReferenceError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newReferenceError(valueOfArrayIndex(argumentList, 0)))
}

func (runtime *_runtime) newSyntaxError(message Value) *_object {
	self := runtime.newErrorObject("SyntaxError", message, nil, 0)
	self.prototype = runtime.global.SyntaxErrorPrototype
	return self
}

func builtinSyntaxError(call FunctionCall) Value {
	return toValue_object(call.runtime.newSyntaxError(call.Argument(0)))
}

func builtinNewSyntaxError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newSyntaxError(valueOfArrayIndex(argumentList, 0)))
}

func builtinURIError(call FunctionCall) Value {
	return toValue_object(call.runtime.newURIError(call.Argument(0)))
}

func builtinNewURIError(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newURIError(valueOfArrayIndex(argumentList, 0)))
}
