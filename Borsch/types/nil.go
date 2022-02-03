package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type NilInstance struct {
	BuiltinInstance
}

func NewNilInstance() NilInstance {
	return NilInstance{
		BuiltinInstance{
			CommonInstance{
				Object: Object{
					typeName:    common.NilTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: Nil,
			},
		},
	}
}

func (t NilInstance) String(common.State) (string, error) {
	return "нуль", nil
}

func (t NilInstance) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t NilInstance) AsBool(common.State) (bool, error) {
	return false, nil
}

func compareNils(_ common.State, _ common.Operator, _ common.Type, other common.Type) (int, error) {
	switch other.(type) {
	case NilInstance:
		return 0, nil
	default:
		// -2 is something other than -1, 0 or 1 and means 'not equals'
		return -2, nil
	}
}

func newNilClass() *Class {
	initAttributes := func() map[string]common.Type {
		return MergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				common.ConstructorName: NewFunctionInstance(
					common.ConstructorName,
					[]FunctionParameter{
						{
							Type:       Nil,
							Name:       "я",
							IsVariadic: false,
							IsNullable: false,
						},
					},
					func(_ common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
						return (*args)[0], nil
					},
					[]FunctionReturnType{
						{
							Type:       Nil,
							IsNullable: false,
						},
					},
					true,
					nil,
					"",
				),
			},
			MakeLogicalOperators(Nil),
			MakeComparisonOperators(Nil, compareNils),
			MakeCommonOperators(Nil),
		)
	}

	return NewBuiltinClass(
		common.NilTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewNilInstance(), nil
		},
	)
}
