package otto

func (rt *_runtime) newErrorObject(name string, message Value, nativeErr error, stackFramesToPop int) *_object {
	self := rt.newClassObject("Error")
	if message.IsDefined() {
		msg := message.string()
		self.defineProperty("message", toValue_string(msg), 0111, false)
		self.value = newError(rt, name, stackFramesToPop, nativeErr, msg)
	} else {
		self.value = newError(rt, name, stackFramesToPop, nativeErr)
	}

	self.defineOwnProperty("stack", _property{
		value: _propertyGetSet{
			rt.newNativeFunction("get", "internal", 0, func(FunctionCall) Value {
				return toValue_string(self.value.(_error).formatWithStack())
			}),
			&_nilGetSetObject,
		},
		mode: modeConfigureMask & modeOnMask,
	}, false)

	return self
}
