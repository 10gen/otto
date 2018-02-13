package otto

import (
	"fmt"
	"testing"
)

type NoopNativeMemUsageChecker struct {
}

func (nnmu NoopNativeMemUsageChecker) NativeMemUsage(goNativeValue interface{}) (uint64, bool) {
	return 0, false
}

func TestMemCheck(t *testing.T) {
	// This is the sum of property names allocated at each new (empty) scope
	var functionOverhead uint64 = 91
	/*
		└──────Property "arguments"
		└─────────Property "length"
		└─────────done: "length" (6) size 8 (total = 14)
		└─────────Property "callee"
		└────────────Property "name"
		└────────────done: "name" (4) size 0 (total = 4)
		└────────────Property "length"
		└────────────done: "length" (6) size 8 (total = 14)
		└────────────Property "prototype"
		└───────────────Property "constructor"
		└───────────────done: "constructor" (11) size 0 (total = 11)
		└────────────done: "prototype" (9) size 19 (total = 28)
		└─────────done: "callee" (6) size 54 (total = 60)
		└──────done: "arguments" (9) size 82 (total = 91)
	*/

	for _, tc := range []struct {
		description      string
		script           string
		expectedSizeDiff uint64
	}{
		{
			"number",
			`x = []
			checkMem()
			x.push(0)
			checkMem()`,
			NumberSize + 1,
		},
		{
			"boolean",
			`x = []
			checkMem()
			x.push(true)
			checkMem()`,
			BoolSize + 1,
		},
		{
			"null",
			`x = []
			checkMem()
			x.push(null)
			checkMem()`,
			EmptySize + 1,
		},
		{
			"undefined",
			`x = []
			checkMem()
			x.push(undefined)
			checkMem()`,
			EmptySize + 1,
		},
		{
			"string",
			`x = []
			checkMem()
			x.push("12345")
			checkMem()`,
			5 + 1,
		},
		{
			"string with multi-byte characters",
			`x = []
			checkMem()
			x.push("\u2318")
			checkMem()`,
			3 + 1, // "x" property name, plus single char with 3-byte width
		},
		{
			"nested objects",
			`y = []
			checkMem()
			y.push({"a":10, "b":"1234", "c":{}})
			checkMem()`,
			EmptySize +
				(1 + NumberSize) + // "a" and number
				(1 + 4) + // "b" and string
				(1 + EmptySize) + //  "c" (an object)
				(1), //"0" attribute for y
		},
		{
			"array of numbers",
			`y = []
			var i = 0;
			checkMem();
			y.push([]);
			for(i=0;i<20;i++){
				y[0].push(i);
			};
			checkMem()`,
			// Array overhead,
			// size of property names,
			// size of property values,
			// plus "length" attribute for the array.
			// plus the length of the "0" property of the outermost array.
			EmptySize + (10 + 10*2) + 20*NumberSize + 6 + NumberSize + 1,
		},
		{
			"overhead of a single new scope",
			`checkMem();
			(function(){
				checkMem();
			})();`, // over
			functionOverhead,
		},
		{
			"overhead of each scope is equivalent regardless of depth",
			`checkMem();
			(function(){
				(function(){
					(function(){
						(function(){
							(function(){
								(function(){
									checkMem();
								})();
							})();
						})();
					})();
				})();
			})();`,
			functionOverhead * 6,
		},
		{"values attached to lexical scope in a function",
			`checkMem();
			(function(){
				var zzzx = 10;
				checkMem();
			})();`,
			// function overhead plus the number value of the "zzzx" property and its string name
			functionOverhead + NumberSize + 4,
		},
		{
			"cyclical data structure",
			// cyclical data structure does not recurse infinitely
			// and does not artificially inflate mem count. The only change in mem
			// between the two checks is for the new property names for "y" and "x".
			`var zzza = {}
			 var zzzb = {}
			 checkMem();
			 zzza.y = zzzb
			 zzzb.x = zzza
			 checkMem()`,
			2, // "x" and "y" property names
		},
	} {
		t.Run(fmt.Sprintf(tc.description), func(t *testing.T) {
			memChecks := []uint64{}
			vm := New()
			vm.Set("checkMem", func(call FunctionCall) Value {

				mem, err := vm.MemUsage(NewMemUsageContext(100, NoopNativeMemUsageChecker{}))
				if err != nil {
					t.Fatal(err)
				}
				memChecks = append(memChecks, mem)
				return UndefinedValue()
			})
			_, err := vm.Run(tc.script)
			if err != nil {
				t.Fatal(err)
			}
			if len(memChecks) < 2 {
				t.Fatalf("expected at least two entries in mem check function, but got %d", len(memChecks))
			}

			memDiff := memChecks[len(memChecks)-1] - memChecks[0]
			if memDiff != tc.expectedSizeDiff {
				t.Fatalf("expected memory change to equal %d but got %d instead", tc.expectedSizeDiff, memDiff)
			}
		})
	}

	t.Log("Check max depth error condition")
	vm := New()
	_, err := vm.Run(`
		var x = {"a": {"b": {"c":{"d":1}}}}
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vm.MemUsage(NewMemUsageContext(3, NoopNativeMemUsageChecker{}))
	if err != ErrMaxDepth {
		t.Fatalf("expected mem check to hit depth limit error, but got nil %v", err)
	}

	_, err = vm.MemUsage(NewMemUsageContext(4, NoopNativeMemUsageChecker{}))
	if err != nil {
		t.Fatalf("expected to NOT hit mem check hit depth limit error, but got %v", err)
	}

}
