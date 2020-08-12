package otto

import "testing"

func TestSymbol(t *testing.T) {
	tt(t, func() {
		test, _ := test()

		// for..in loop
		test(`
			var obj = {
				a: "1",
				b: "2",
				c: "3"
			};
			var res = [];

			for (var i in obj) {
				res.push(i + ': ' + obj[i]);
			}

			res
		`, "a: 1,b: 2,c: 3")

		// new symbols should not equal
		test(`
			Symbol("foo") === Symbol("foo");
		`, false)

		// typeof symbol
		test(`
			typeof Symbol('foo');
		`, "symbol")

		// toString
		test(`
			Symbol('foo').toString()
		`, "Symbol(foo)")

		// Symbol.fors should equal
		test(`
			Symbol.for('bar') === Symbol.for('bar');
		`, true)

		// keyFor should return the description for a global symbol
		test(`
			var sym = Symbol.for('hi');
			Symbol.keyFor(sym);
		`, "hi")

		// keyFor should return undefined for a local symbol
		test(`
			var sym = Symbol('hi');
			Symbol.keyFor(sym);
		`, "undefined")

		// toStringTag should be a symbol
		test(`
			Symbol.toStringTag.toString();
		`, "Symbol(Symbol.toStringTag)")

		// toStringTag should replace the Object in toString
		test(`
			var obj = {};
			Object.defineProperty(obj, Symbol.toStringTag, { value: 'Something' })

			Object.prototype.toString.call(obj);
		`, "[object Something]")

		// iterator should be a symbol
		test(`
			Symbol.iterator.toString();
		`, "Symbol(Symbol.iterator)")

		// iterating through a list should work
		test(`
			var arr = ['a', 'b', 'c'];
			var eArr = arr[Symbol.iterator]();
			eArr.next().value + "," + eArr.next().value + "," + eArr.next().value
		`, "a,b,c")
	})
}
