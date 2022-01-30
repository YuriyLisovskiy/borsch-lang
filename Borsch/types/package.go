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

func (p PackageInstance) String(common.Context) string {
	return fmt.Sprintf("<пакет '%s'>", p.Name)
}

func (p PackageInstance) Representation(ctx common.Context) string {
	return p.String(ctx)
}

func (p PackageInstance) AsBool(common.Context) bool {
	return true
}

func (p PackageInstance) GetTypeName() string {
	return p.GetPrototype().GetTypeName()
}

func (p PackageInstance) GetAttribute(name string) (common.Type, error) {
	if attribute, err := p.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return p.GetPrototype().GetAttribute(name)
}

func (p PackageInstance) SetAttribute(name string, value common.Type) (common.Type, error) {
	if p.IsBuiltin {
		if p.Object.HasAttribute(name) || p.GetPrototype().HasAttribute(name) {
			return nil, util.AttributeIsReadOnlyError(p.GetTypeName(), name)
		}

		return nil, util.AttributeNotFoundError(p.GetTypeName(), name)
	}

	err := p.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return p, nil
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
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
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
	}

	return NewBuiltinClass(
		common.PackageTypeName,
		BuiltinPackage,
		initAttributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
