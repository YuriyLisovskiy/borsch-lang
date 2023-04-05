package interpreter

import (
	"errors"
	"fmt"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
	common2 "github.com/YuriyLisovskiy/borsch-lang/internal/common"
)

type StateImpl struct {
	parent     State
	context    types2.Context
	pkg        types2.Object
	stacktrace *common2.StackTrace
}

func NewInitialState(
	context types2.Context,
	pkg types2.Object,
	stacktrace *common2.StackTrace,
) State {
	return &StateImpl{
		parent:     nil,
		context:    context,
		pkg:        pkg,
		stacktrace: stacktrace,
	}
}

func (s *StateImpl) Parent() State {
	return s.parent
}

func (s *StateImpl) NewChild() State {
	return &StateImpl{
		parent:     s,
		context:    s.context,
		pkg:        s.pkg,
		stacktrace: s.stacktrace,
	}
}

func (s *StateImpl) Context() types2.Context {
	if s.context != nil {
		return s.context
	}

	panic("state: context is nil")
}

func (s *StateImpl) Package() types2.Object {
	if s.pkg != nil {
		return s.pkg
	}

	panic("state: package is nil")
}

func (s *StateImpl) StackTrace() *common2.StackTrace {
	return s.stacktrace
}

func (s *StateImpl) PackageOrNil() types2.Object {
	return s.pkg
}

func (s *StateImpl) WithContext(ctx types2.Context) State {
	s.context = ctx
	return s
}

func (s *StateImpl) WithPackage(pkg types2.Object) State {
	s.pkg = pkg
	return s
}

func (s *StateImpl) RuntimeError(message string, statement common2.Statement) error {
	if statement != nil {
		s.Trace(statement, "")
	}

	return errors.New(fmt.Sprintf("Помилка виконання: %s", message))
}

func (s *StateImpl) Trace(statement common2.Statement, place string) {
	s.stacktrace.Push(common2.NewTraceRow(statement.Position(), statement.String(), place))
}

func (s *StateImpl) PopTrace() {
	if !s.stacktrace.IsEmpty() {
		s.stacktrace.Pop()
	}
}
