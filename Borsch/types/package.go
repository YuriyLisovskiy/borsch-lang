package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageInstance struct {
	Object
	IsBuiltin bool
	Name      string
	Parent    string
}

func NewPackageInstance(isBuiltin bool, name, parent string, attributes map[string]Type) *PackageInstance {
	return &PackageInstance{
		IsBuiltin: isBuiltin,
		Name:      name,
		Parent:    parent,
		Object: Object{
			typeName:    GetTypeName(PackageTypeHash),
			Attributes:  attributes,
			callHandler: nil,
		},
	}
}

func (t PackageInstance) String() string {
	name := t.Name
	if t.IsBuiltin {
		name = "АТБ"
	}

	return fmt.Sprintf("<пакет '%s'>", name)
}

func (t PackageInstance) Representation() string {
	return t.String()
}

func (t PackageInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t PackageInstance) AsBool() bool {
	return true
}

func (t PackageInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (t PackageInstance) SetAttribute(name string, value Type) (Type, error) {
	if t.IsBuiltin {
		if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
			return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
		}

		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (PackageInstance) GetClass() *Class {
	return Package
}

func comparePackages(self Type, other Type) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case *PackageInstance, PackageInstance:
		return -2, util.RuntimeError(
			fmt.Sprintf(
				"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	default:
		return -2, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func NewPackageClass() *Class {
	attributes := mergeAttributes(
		// map[string]Type{
		// 	// TODO: add doc
		// 	ops.ConstructorName: NewFunctionInstance(
		// 		ops.ConstructorName,
		// 		[]FunctionArgument{
		// 			{
		// 				TypeHash:   PackageTypeHash,
		// 				Name:       "я",
		// 				IsVariadic: false,
		// 				IsNullable: false,
		// 			},
		// 			{
		// 				TypeHash:   StringTypeHash,
		// 				Name:       "шлях",
		// 				IsVariadic: false,
		// 				IsNullable: false,
		// 			},
		// 		},
		// 		func(args *[]Type, _ *map[string]Type) (Type, error) {
		// 			self, err := importPackage((*args)[1].(StringInstance).Value)
		// 			// self, err := handler((*args)[1:]...)
		// 			if err != nil {
		// 				return nil, err
		// 			}
		//
		// 			(*args)[0] = self
		// 			return NewNilInstance(), nil
		// 		},
		// 		[]FunctionReturnType{
		// 			{
		// 				TypeHash:   NilTypeHash,
		// 				IsNullable: false,
		// 			},
		// 		},
		// 		true,
		// 		nil,
		// 		"", // TODO: add doc
		// 	),
		// },
		makeLogicalOperators(PackageTypeHash),
		makeComparisonOperators(PackageTypeHash, comparePackages),
		makeCommonOperators(PackageTypeHash),
	)
	return NewBuiltinClass(
		PackageTypeHash,
		BuiltinPackage,
		attributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
