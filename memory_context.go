package otto

import "errors"

type visitTracker struct {
	objsVisited    map[*_object]bool
	stashesVisited map[_stash]bool
}

func (vt visitTracker) IsObjVisited(obj *_object) bool {
	_, ok := vt.objsVisited[obj]
	return ok
}

func (vt visitTracker) VisitObj(obj *_object) {
	vt.objsVisited[obj] = true
}

func (vt visitTracker) IsStashVisited(stash _stash) bool {
	_, ok := vt.stashesVisited[stash]
	return ok
}

func (vt visitTracker) VisitStash(stash _stash) {
	vt.stashesVisited[stash] = true
}

type depthTracker struct {
	curDepth int
	maxDepth int
}

func (dt depthTracker) Depth() int {
	return dt.curDepth
}

func (dt *depthTracker) Descend() error {
	if dt.curDepth == dt.maxDepth {
		return ErrMaxDepth
	}
	dt.curDepth++
	return nil
}

func (dt *depthTracker) Ascend() {
	if dt.curDepth == 0 {
		panic("can't ascend with depth 0")
	}
	dt.curDepth--
}

type NativeMemUsageChecker interface {
	NativeMemUsage(goNativeValue interface{}) (uint64, bool)
}

type MemUsageContext struct {
	visitTracker
	*depthTracker
	NativeMemUsageChecker
}

func NewMemUsageContext(maxDepth int, nativeChecker NativeMemUsageChecker) *MemUsageContext {
	return &MemUsageContext{
		visitTracker:          visitTracker{objsVisited: map[*_object]bool{}, stashesVisited: map[_stash]bool{}},
		depthTracker:          &depthTracker{curDepth: 0, maxDepth: maxDepth},
		NativeMemUsageChecker: nativeChecker,
	}
}

var (
	ErrMaxDepth = errors.New("reached max depth")
)
