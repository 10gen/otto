package otto

type _symbolObject struct {
	description interface{}
}

func (runtime *_runtime) newSymbolObject(description interface{}) *_object {
	self := runtime.newObject()
	self.class = "Symbol"

	symbol := _symbolObject{
		description: description,
	}
	self.value = symbol
	self.defineProperty("description", toValue(description), 0000, false)

	if _, ok := runtime.symbols[description]; !ok {
		runtime.symbols[description] = toValue_object(self)
	}

	return self
}
