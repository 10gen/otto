package otto

// _scope:
// entryFile
// entryIdx
// top?
// outer => nil

// _stash:
// lexical
// variable
//
// _thisStash (ObjectEnvironment)
// _fnStash
// _dclStash

// An ECMA-262 ExecutionContext
type _scope struct {
	lexical  _stash
	variable _stash
	this     *_object
	eval     bool // Replace this with kind?
	outer    *_scope
	depth    int

	frame _frame
}

func (self *_scope) MemUsage(ctx *MemUsageContext) (uint64, error) {
	total := uint64(0)
	if self.this != nil {
		scopesize, err := self.this.MemUsage(ctx)
		total += scopesize
		if err != nil {
			return total, err
		}
	}
	if self.lexical != nil {
		lexicalSize, err := self.lexical.MemUsage(ctx)
		total += lexicalSize
		if err != nil {
			return total, err
		}
	}
	if self.variable != nil {
		variableSize, err := self.variable.MemUsage(ctx)
		total += variableSize
		if err != nil {
			return total, err
		}
	}
	if self.outer != nil {
		if err := ctx.Descend(); err != nil {
			return total, err
		}
		inc, err := self.outer.MemUsage(ctx)
		ctx.Ascend()
		total += inc
		return total, err
	}
	return total, nil
}

func newScope(lexical _stash, variable _stash, this *_object) *_scope {
	return &_scope{
		lexical:  lexical,
		variable: variable,
		this:     this,
	}
}
