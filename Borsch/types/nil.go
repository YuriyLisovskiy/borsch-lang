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
			typeName:    common.NilTypeName,
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t NilInstance) String(common.Context) string {
	return "нуль"
}

func (t NilInstance) Representation(ctx common.Context) string {
	return t.String(ctx)
}

func (t NilInstance) AsBool(common.Context) bool {
	return false
}

func (t NilInstance) GetTypeName() string {
	return t.GetPrototype().GetTypeName()
}

func (t NilInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (t NilInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (NilInstance) GetPrototype() *Class {
	return Nil
}

func compareNils(_ common.Context, _ common.Type, other common.Type) (int, error) {
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
					func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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