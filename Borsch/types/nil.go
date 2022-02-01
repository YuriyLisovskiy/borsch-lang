package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

type NilInstance struct {
	BuiltinObject
}

func NewNilInstance() NilInstance {
	return NilInstance{
		BuiltinObject{
			CommonObject{
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

func (t NilInstance) String(common.State) string {
	return "нуль"
}

func (t NilInstance) Representation(state common.State) string {
	return t.String(state)
}

func (t NilInstance) AsBool(common.State) bool {
	return false
}

func compareNils(_ common.State, _ common.Type, other common.Type) (int, error) {
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
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: NewFunctionInstance(
					ops.ConstructorName,
					[]FunctionArgument{
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
			makeLogicalOperators(Nil),
			makeComparisonOperators(Nil, compareNils),
			makeCommonOperators(Nil),
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
