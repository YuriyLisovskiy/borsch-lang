package types

import (
	"errors"
	"fmt"

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
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t NilInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (NilInstance) GetClass() *Class {
	return Nil
}

func compareNils(self common.Type, other common.Type) (int, error) {
	switch right := other.(type) {
	case NilInstance:
		return 0, nil
	default:
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
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
