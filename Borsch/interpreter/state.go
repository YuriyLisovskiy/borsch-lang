package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type StateImpl struct {
	parent     State
	context    types.Context
	pkg        types.Object
	stacktrace *common.StackTrace
}

func NewInitialState(
	context types.Context,
	pkg types.Object,
	stacktrace *common.StackTrace,
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

func (s *StateImpl) Context() types.Context {
	if s.context != nil {
		return s.context
	}

	panic("state: context is nil")
}

func (s *StateImpl) Package() types.Object {
	if s.pkg != nil {
		return s.pkg
	}

	panic("state: package is nil")
}

func (s *StateImpl) StackTrace() *common.StackTrace {
	return s.stacktrace
}

func (s *StateImpl) PackageOrNil() types.Object {
	return s.pkg
}

func (s *StateImpl) WithContext(ctx types.Context) State {
	s.context = ctx
	return s
}

func (s *StateImpl) WithPackage(pkg types.Object) State {
	s.pkg = pkg
	return s
}

func (s *StateImpl) RuntimeError(message string, statement common.Statement) error {
	if statement != nil {
		s.Trace(statement, "")
	}

	return errors.New(fmt.Sprintf("Помилка виконання: %s", message))
}

func (s *StateImpl) Trace(statement common.Statement, place string) {
	s.stacktrace.Push(common.NewTraceRow(statement.Position(), statement.String(), place))
}

func (s *StateImpl) PopTrace() {
	if !s.stacktrace.IsEmpty() {
		s.stacktrace.Pop()
	}
}
