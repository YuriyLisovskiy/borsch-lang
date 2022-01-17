package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type NilInstance struct {
	Object
}

func NewNilInstance() NilInstance {
	return NilInstance{
		Object: Object{
			typeName:    GetTypeName(NilTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t NilInstance) String() string {
	return "нуль"
}

func (t NilInstance) Representation() string {
	return t.String()
}

func (t NilInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t NilInstance) AsBool() bool {
	return false
}

func (t NilInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}
	
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (t NilInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (NilInstance) GetClass() *Class {
	return Nil
}

func compareNils(_ common.Type, other common.Type) (int, error) {
	switch other.(type) {
	case NilInstance:
		return 0, nil
	default:
		// -2 is something other than -1, 0 or 1 and means 'not equals'
		return -2, nil
	}
}

func newNilClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			// TODO: add doc
			ops.ConstructorName: NewFunctionInstance(
				ops.ConstructorName,
				[]FunctionArgument{
					{
						TypeHash:   NilTypeHash,
						Name:       "я",
						IsVariadic: false,
						IsNullable: false,
					},
				},
				func(_ interface{}, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
					return (*args)[0], nil
				},
				[]FunctionReturnType{
					{
						TypeHash:   NilTypeHash,
						IsNullable: false,
					},
				},
				true,
				nil,
				"",
			),
		},
		makeLogicalOperators(NilTypeHash),
		makeComparisonOperators(NilTypeHash, compareNils),
		makeCommonOperators(NilTypeHash),
	)
	return NewBuiltinClass(
		NilTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewNilInstance(), nil
		},
	)
}
