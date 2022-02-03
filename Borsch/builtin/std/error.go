package std

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

type ErrorInstance struct {
	types.BuiltinInstance
	message string
}

func NewErrorInstance(message string) *ErrorInstance {
	return &ErrorInstance{
		BuiltinInstance: types.BuiltinInstance{
			CommonInstance: *types.NewCommonInstance(types.NewObject(common.ErrorTypeName, nil, nil), ErrorClass),
		},
		message: message,
	}
}

func (t ErrorInstance) String(common.State) (string, error) {
	return t.message, nil
}

func (t ErrorInstance) Representation(common.State) (string, error) {
	return fmt.Sprintf("%s(\"%s\")", t.GetTypeName(), t.message), nil
}

func (t ErrorInstance) AsBool(common.State) (bool, error) {
	return false, nil
}

func compareErrors(_ common.State, self common.Type, other common.Type) (int, error) {
	if _, ok := other.(*ErrorInstance); ok {
		if self == other {
			return 0, nil
		}
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newErrorClass() *types.Class {
	initAttributes := func() map[string]common.Type {
		return types.MergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				common.ConstructorName: types.NewFunctionInstance(
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
					func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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
					},
					[]types.FunctionReturnType{
						{
							Type:       types.Nil,
							IsNullable: false,
						},
					},
					true,
					nil,
					"",
				),
				"повідомлення": types.NewFunctionInstance(
					"повідомлення",
					[]types.FunctionParameter{
						{
							Type:       ErrorClass,
							Name:       "я",
							IsVariadic: false,
							IsNullable: false,
						},
					},
					func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
						msg, err := (*args)[0].String(state)
						if err != nil {
							return nil, err
						}

						return types.NewStringInstance(msg), nil
					},
					[]types.FunctionReturnType{
						{
							Type:       types.String,
							IsNullable: false,
						},
					},
					true,
					nil,
					"",
				),
			},
			types.MakeLogicalOperators(ErrorClass),
			types.MergeAttributes(
				map[string]common.Type{
					common.EqualsOp.Name(): types.NewComparisonOperator(
						// TODO: add doc
						common.EqualsOp, ErrorClass, "", compareErrors, func(res int) bool {
							return res == 0
						},
					),
					common.NotEqualsOp.Name(): types.NewComparisonOperator(
						// TODO: add doc
						common.NotEqualsOp, ErrorClass, "", compareErrors, func(res int) bool {
							return res != 0
						},
					),
				},
			),
			types.MakeCommonOperators(ErrorClass),
		)
	}

	return types.NewBuiltinClass(
		common.ErrorTypeName,
		types.BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewErrorInstance(""), nil
		},
	)
}
