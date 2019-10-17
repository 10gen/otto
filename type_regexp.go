package otto

import (
	"fmt"
	"unicode/utf8"

	"github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/regexp"
	"github.com/robertkrimen/otto/regexp/pcre"
	"github.com/robertkrimen/otto/regexp/re2"

	"github.com/dlclark/regexp2"
)

type _regExpObject struct {
	regularExpression regexp.Regexp
	global            bool
	ignoreCase        bool
	multiline         bool
	source            string
	flags             string
}

func (runtime *_runtime) newRegExpObject(pattern string, flags string) *_object {
	self := runtime.newObject()
	self.class = "RegExp"

	global := false
	ignoreCase := false
	multiline := false
	re2flags := ""
	var pcreFlags regexp2.RegexOptions = regexp2.ECMAScript

	// TODO Maybe clean up the panicking here... TypeError, SyntaxError, ?

	for _, chr := range flags {
		switch chr {
		case 'g':
			if global {
				panic(runtime.panicSyntaxError("newRegExpObject: %s %s", pattern, flags))
			}
			global = true
		case 'm':
			if multiline {
				panic(runtime.panicSyntaxError("newRegExpObject: %s %s", pattern, flags))
			}
			multiline = true
			re2flags += "m"
			pcreFlags |= regexp2.Multiline
		case 'i':
			if ignoreCase {
				panic(runtime.panicSyntaxError("newRegExpObject: %s %s", pattern, flags))
			}
			ignoreCase = true
			re2flags += "i"
			pcreFlags |= regexp2.IgnoreCase
		}
	}

	transformedPattern, transformErr := parser.TransformRegExp(pattern)
	if transformedPattern == "" && transformErr != nil {
		panic(runtime.panicTypeError("Invalid regular expression: %s", transformErr.Error()))
	}

	var regularExpression regexp.Regexp
	var regexpErr error
	if transformErr != nil {
		regularExpression, regexpErr = pcre.New(transformedPattern, pcreFlags)
	} else {
		if len(re2flags) > 0 {
			transformedPattern = fmt.Sprintf("(?%s:%s)", re2flags, transformedPattern)
		}
		regularExpression, regexpErr = re2.New(transformedPattern)
	}
	if regexpErr != nil {
		panic(runtime.panicSyntaxError("Invalid regular expression: %s", regexpErr.Error()))
	}

	self.value = _regExpObject{
		regularExpression: regularExpression,
		global:            global,
		ignoreCase:        ignoreCase,
		multiline:         multiline,
		source:            pattern,
		flags:             flags,
	}
	self.defineProperty("global", toValue_bool(global), 0, false)
	self.defineProperty("ignoreCase", toValue_bool(ignoreCase), 0, false)
	self.defineProperty("multiline", toValue_bool(multiline), 0, false)
	self.defineProperty("lastIndex", toValue_int(0), 0100, false)
	self.defineProperty("source", toValue_string(pattern), 0, false)
	return self
}

func (self *_object) regExpValue() _regExpObject {
	value, _ := self.value.(_regExpObject)
	return value
}

func execRegExp(this *_object, target string) (match bool, result []int) {
	if this.class != "RegExp" {
		panic(this.runtime.panicTypeError("Calling RegExp.exec on a non-RegExp object"))
	}
	lastIndex := this.get("lastIndex").number().int64
	index := lastIndex
	var resultErr error
	global := this.get("global").bool()
	if !global {
		index = 0
	}
	if 0 > index || index > int64(len(target)) {
	} else {
		result, resultErr = this.regExpValue().regularExpression.FindStringSubmatchIndex(target[index:])
	}
	if resultErr != nil || result == nil {
		//this.defineProperty("lastIndex", toValue_(0), 0111, true)
		this.put("lastIndex", toValue_int(0), true)
		return // !match
	}
	match = true
	startIndex := index
	endIndex := int(lastIndex) + result[1]
	// We do this shift here because the .FindStringSubmatchIndex above
	// was done on a local subordinate slice of the string, not the whole string
	for index, _ := range result {
		result[index] += int(startIndex)
	}
	if global {
		//this.defineProperty("lastIndex", toValue_(endIndex), 0111, true)
		this.put("lastIndex", toValue_int(endIndex), true)
	}
	return // match
}

func execResultToArray(runtime *_runtime, target string, result []int) *_object {
	captureCount := len(result) / 2
	valueArray := make([]Value, captureCount)
	for index := 0; index < captureCount; index++ {
		offset := 2 * index
		if result[offset] != -1 {
			valueArray[index] = toValue_string(target[result[offset]:result[offset+1]])
		} else {
			valueArray[index] = Value{}
		}
	}
	matchIndex := result[0]
	if matchIndex != 0 {
		matchIndex = 0
		// Find the rune index in the string, not the byte index
		for index := 0; index < result[0]; {
			_, size := utf8.DecodeRuneInString(target[index:])
			matchIndex += 1
			index += size
		}
	}
	match := runtime.newArrayOf(valueArray)
	match.defineProperty("input", toValue_string(target), 0111, false)
	match.defineProperty("index", toValue_int(matchIndex), 0111, false)
	return match
}
