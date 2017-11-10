package otto

// A NativeClass enables a way to dynamically create classes in Go
// to be represented in JavaScript via properties.
// The struct contains the function for the JavaScript class
// and a function to directly instantiate an instance of it.
type NativeClass struct {

	// Function is the function representing the class.
	Function Value

	// InstanceOf is a "blessed" function that allows any
	// value passed to it to be set as the underlying value
	// representing the class.
	InstanceOf func(value interface{}) Value
}

// CreateNativeClass creates a native class of the given name. Required
// is a constructor that will be called upon calling the function (with or without 'new').
// Optional arguments are the properties that should exist on the class and function
// respectively.
func (self Otto) CreateNativeClass(
	className string,
	ctor func(call FunctionCall) interface{},
	classProps []Property,
	funcProps []Property,
) NativeClass {
	classProto := &_object{
		runtime:       self.runtime,
		class:         className,
		objectClass:   _classObject,
		prototype:     self.runtime.global.ObjectPrototype,
		extensible:    true,
		value:         nil,
		property:      map[string]_property{},
		propertyOrder: []string{},
	}

	for _, prop := range classProps {
		classProto.propertyOrder = append(classProto.propertyOrder, prop.Name)
		classProto.property[prop.Name] = _property{
			value: prop.Value,
			mode:  0,
		}
	}

	classFunc := &_object{
		runtime:     self.runtime,
		class:       "Function",
		objectClass: _classObject,
		prototype:   self.runtime.global.FunctionPrototype,
		extensible:  true,
		value: _nativeFunctionObject{
			name: className,
			call: func(call FunctionCall) Value {
				obj := self.runtime.newObject()
				obj.class = className
				obj.value = _goNativeValue{ctor(call)}
				obj.prototype = classProto
				return toValue_object(obj)
			},
			construct: func(self *_object, argumentList []Value) Value {
				obj := self.runtime.newObject()
				obj.class = className

				call := FunctionCall{
					runtime:      self.runtime,
					eval:         false,
					This:         toValue_object(obj),
					ArgumentList: argumentList,
					Otto:         self.runtime.otto,
				}
				obj.value = _goNativeValue{ctor(call)}
				obj.prototype = classProto
				return toValue_object(obj)
			},
		},
		property: map[string]_property{
			"prototype": _property{
				mode: 0,
				value: Value{
					kind:  valueObject,
					value: classProto,
				},
			},
		},
		propertyOrder: []string{},
	}
	for _, prop := range funcProps {
		classFunc.propertyOrder = append(classFunc.propertyOrder, prop.Name)
		classFunc.property[prop.Name] = _property{
			value: prop.Value,
			mode:  0,
		}
	}

	blessedInstanceOf := func(blessedValue interface{}) Value {
		obj := self.runtime.newObject()
		obj.class = className
		obj.value = _goNativeValue{blessedValue}
		obj.prototype = classProto
		return toValue_object(obj)
	}
	return NativeClass{
		Function:   toValue_object(classFunc),
		InstanceOf: blessedInstanceOf,
	}
}

// _goNativeValue is a wrapper type to signify that this is
// an instance of a native class.
type _goNativeValue struct {
	value interface{}
}
