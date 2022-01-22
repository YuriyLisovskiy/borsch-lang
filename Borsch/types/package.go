package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageInstance struct {
	Object
	IsBuiltin bool
	Name      string
	Parent    *PackageInstance
}

func NewPackageInstance(
	isBuiltin bool,
	name string,
	parent *PackageInstance,
	attributes map[string]common.Type,
) *PackageInstance {
	return &PackageInstance{
		IsBuiltin: isBuiltin,
		Name:      name,
		Parent:    parent,
		Object: Object{
			typeName:    common.PackageTypeName,
			Attributes:  attributes,
			callHandler: nil,
		},
	}
}

func (t PackageInstance) String(common.Context) string {
	return fmt.Sprintf("<пакет '%s'>", t.Name)
}

func (t PackageInstance) Representation(ctx common.Context) string {
	return t.String(ctx)
}

func (t PackageInstance) AsBool(common.Context) bool {
	return true
}

func (t PackageInstance) GetTypeName() string {
	return t.GetPrototype().GetTypeName()
}

func (t PackageInstance) GetAttribute(name string) (common.Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (t PackageInstance) SetAttribute(name string, value common.Type) (common.Type, error) {
	if t.IsBuiltin {
		if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
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

func (PackageInstance) GetPrototype() *Class {
	return Package
}

func comparePackages(_ common.Context, self common.Type, other common.Type) (int, error) {
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
		map[string]common.Type{
			// TODO: add doc
			// ops.ConstructorName: NewFunctionInstance(
			// 	ops.ConstructorName,
			// 	[]FunctionArgument{
			// 		{
			// 			TypeHash:   PackageTypeHash,
			// 			Name:       "я",
			// 			IsVariadic: false,
			// 			IsNullable: false,
			// 		},
			// 		{
			// 			TypeHash:   StringTypeHash,
			// 			Name:       "шлях",
			// 			IsVariadic: false,
			// 			IsNullable: false,
			// 		},
			// 	},
			// 	func(ctx interface{}, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			// 		p, err := ImportPackage(baseScope, (*args)[0].(StringInstance).Value, ctx.(common.CallContext).Parser)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			//
			// 		return p, nil
			// 		self, err := handler((*args)[1:]...)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			//
			// 		(*args)[0] = self
			// 		return NewNilInstance(), nil
			// 	},
			// 	[]FunctionReturnType{
			// 		{
			// 			TypeHash:   NilTypeHash,
			// 			IsNullable: false,
			// 		},
			// 	},
			// 	true,
			// 	nil,
			// 	"", // TODO: add doc
			// ),
		},
		makeLogicalOperators(Package),
		makeComparisonOperators(Package, comparePackages),
		makeCommonOperators(Package),
	)
	return NewBuiltinClass(
		common.PackageTypeName,
		BuiltinPackage,
		attributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
