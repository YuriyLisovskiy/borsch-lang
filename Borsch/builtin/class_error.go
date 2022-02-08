package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type ErrorInstance struct {
	types.ClassInstance
	message string
}

func NewErrorInstance(message string) *ErrorInstance {
	return &ErrorInstance{
		ClassInstance: *types.NewClassInstance(ErrorClass, nil),
		message:       message,
	}
}

func (t ErrorInstance) String(common.State) (string, error) {
	return t.message, nil
}

func (t ErrorInstance) Representation(common.State) (string, error) {
	return fmt.Sprintf("%s(\"%s\")", t.GetTypeName(), t.message), nil
}

func (t ErrorInstance) AsBool(common.State) (bool, error) {
	return true, nil
}

func compareErrors(_ common.State, _ common.Operator, self common.Value, other common.Value) (int, error) {
	if _, ok := other.(*ErrorInstance); ok {
		if self == other {
			return 0, nil
		}
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func errorEvalConstructor(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	rawParts := (*args)[1:]
	self := (*args)[0].(*ErrorInstance)
	for _, rawPart := range rawParts {
		part, err := rawPart.String(state)
		if err != nil {
			return nil, err
		}

		self.message += part
	}

	(*args)[0] = self
	return types.NewNilInstance(), nil
}

func errorEvalMessage(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	msg, err := (*args)[0].String(state)
	if err != nil {
		return nil, err
	}

	return types.NewStringInstance(msg), nil
}

func makeErrorOperator_Constructor() common.Value {
	return types.NewFunctionInstance(
		common.ConstructorName,
		[]types.FunctionParameter{
			{
				Type:       ErrorClass,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       types.String,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		errorEvalConstructor,
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: false,
			},
		},
		true,
		nil,
		"",
	)
}

func makeErrorMethod_Message(name string) common.Value {
	return types.NewFunctionInstance(
		name,
		[]types.FunctionParameter{
			{
				Type:       ErrorClass,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		errorEvalMessage,
		[]types.FunctionReturnType{
			{
				Type:       types.String,
				IsNullable: false,
			},
		},
		true,
		nil,
		"",
	)
}

func newErrorClass() *types.Class {
	return &types.Class{
		Name:    common.ErrorTypeName,
		IsFinal: false,
		Bases:   []*types.Class{},
		Parent:  types.BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = types.MergeAttributes(
				map[string]common.Value{
					// TODO: add doc
					common.ConstructorName: makeErrorOperator_Constructor(),
					"повідомлення":         makeErrorMethod_Message("повідомлення"),
					common.EqualsOp.Name(): types.MakeComparisonOperator(
						// TODO: add doc
						common.EqualsOp, ErrorClass, "", compareErrors, func(res int) bool {
							return res == 0
						},
					),
					common.NotEqualsOp.Name(): types.MakeComparisonOperator(
						// TODO: add doc
						common.NotEqualsOp, ErrorClass, "", compareErrors, func(res int) bool {
							return res != 0
						},
					),
				},
				types.MakeLogicalOperators(ErrorClass),
				types.MakeCommonOperators(ErrorClass),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewErrorInstance(""), nil
		},
	}
}
