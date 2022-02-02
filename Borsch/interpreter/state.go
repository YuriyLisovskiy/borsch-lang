package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type StateImpl struct {
	parser         common.Parser
	interpreter    common.Interpreter
	context        common.Context
	currentPackage common.Type
}

func NewState(
	parser common.Parser,
	interpreter common.Interpreter,
	context common.Context,
	currentPackage common.Type,
) *StateImpl {
	return &StateImpl{
		parser:         parser,
		interpreter:    interpreter,
		context:        context,
		currentPackage: currentPackage,
	}
}

func (s *StateImpl) GetParser() common.Parser {
	if s.parser != nil {
		return s.parser
	}

	panic("state: parser is nil")
}

func (s *StateImpl) GetInterpreter() common.Interpreter {
	if s.interpreter != nil {
		return s.interpreter
	}

	panic("state: interpreter is nil")
}

func (s *StateImpl) GetContext() common.Context {
	if s.context != nil {
		return s.context
	}

	panic("state: context is nil")
}

func (s *StateImpl) GetCurrentPackage() common.Type {
	if s.currentPackage != nil {
		return s.currentPackage
	}

	panic("state: current package is nil")
}

func (s *StateImpl) GetCurrentPackageOrNil() common.Type {
	return s.currentPackage
}

func (s *StateImpl) WithContext(ctx common.Context) common.State {
	return &StateImpl{
		parser:         s.parser,
		interpreter:    s.interpreter,
		context:        ctx,
		currentPackage: s.currentPackage,
	}
}

func (s *StateImpl) WithPackage(pkg common.Type) common.State {
	return &StateImpl{
		parser:         s.parser,
		interpreter:    s.interpreter,
		context:        s.context,
		currentPackage: pkg,
	}
}
