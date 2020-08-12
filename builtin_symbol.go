package otto

import (
	"fmt"
)

func builtinSymbol(call FunctionCall) Value {
	return toValue_object(builtinNewSymbolNative(call.runtime, call.ArgumentList))
}

func builtinNewSymbol(self *_object, argumentList []Value) Value {
	return toValue_object(builtinNewSymbolNative(self.runtime, argumentList))
}

func builtinNewSymbolNative(runtime *_runtime, argumentList []Value) *_object {
	var description interface{}
	if len(argumentList) > 0 {
		description = argumentList[0].value
	}

	return runtime.newSymbol(description)
}

func builtinSymbol_for(call FunctionCall) Value {
	if len(call.ArgumentList) < 1 {
		panic(call.runtime.panicTypeError("Symbol.for takes one argument -- for(data)"))
	}

	description := call.Argument(0)
	if symbol, ok := call.runtime.symbols[description]; ok {
		return symbol
	}

	symbol := builtinNewSymbolNative(call.runtime, call.ArgumentList)
	call.runtime.symbols[description] = toValue_object(symbol)
	return toValue_object(symbol)
}

func builtinSymbol_keyFor(call FunctionCall) Value {
	if len(call.ArgumentList) < 1 {
		panic(call.runtime.panicTypeError("Symbol.keyFor takes one argument -- keyFor(data)"))
	}

	lookingForSymbol := call.Argument(0)
	for symbolVal, symbol := range call.runtime.symbols {
		if symbol == lookingForSymbol {
			return toValue(symbolVal)
		}
	}
	return UndefinedValue()
}

func builtinSymbol_toString(call FunctionCall) Value {
	object := call.thisClassObject("Symbol") // Should throw a TypeError unless Symbol
	switch sym := object.value.(type) {
	case _symbolObject:
		return toValue_string(fmt.Sprintf("Symbol(%s)", sym.description))
	}

	panic(call.runtime.panicTypeError("Symbol.toString()"))
}

func builtinSymbol_valueOf(call FunctionCall) Value {
	object := call.thisClassObject("Symbol") // Should throw a TypeError unless Symbol
	return toValue_object(object)
}

