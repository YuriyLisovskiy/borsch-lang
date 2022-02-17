package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type NilInstance struct {
	BuiltinInstance
}

func NewNilInstance() NilInstance {
	return NilInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance: ClassInstance{
				class:      NilClass,
				attributes: map[string]common.Object{},
				address:    "",
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

func compareNils(_ common.State, _ common.Operator, _ common.Object, other common.Object) (int, error) {
	switch other.(type) {
	case NilInstance:
		return 0, nil
	default:
		// -2 is something other than -1, 0 or 1 and means 'not equals'
		return -2, nil
	}
}

func nilMethod_Constructor() common.Object {
	return NewFunctionInstance(
		common.ConstructorName,
		[]FunctionParameter{
			{
				Type:       NilClass,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(_ common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
			return (*args)[0], nil
		},
		[]FunctionReturnType{
			{
				Type:       NilClass,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func newNilClass() *Class {
	return &Class{
		Name:    common.NilTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Object) {
			*attrs = MergeAttributes(
				map[string]common.Object{
					common.ConstructorName: nilMethod_Constructor(),
				},
				MakeLogicalOperators(NilClass),
				MakeComparisonOperators(NilClass, compareNils),
				MakeCommonOperators(NilClass),
			)
		},
		GetEmptyInstance: func() (common.Object, error) {
			return NewNilInstance(), nil
		},
	}
}
