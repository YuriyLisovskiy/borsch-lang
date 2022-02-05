package std

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

func newErrorClass() *types.Class {
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = types.MergeAttributes(
			map[string]common.Value{
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
					func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
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
					func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
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
				map[string]common.Value{
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

	return types.NewClass(
		common.ErrorTypeName,
		nil,
		types.BuiltinPackage,
		initAttributes,
		func() (common.Value, error) {
			return NewErrorInstance(""), nil
		},
	)
}
