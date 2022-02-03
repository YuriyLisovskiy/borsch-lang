package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageInstance struct {
	CommonInstance
	IsBuiltin bool
	Name      string
	Parent    *PackageInstance

	ctx common.Context
}

func NewPackageInstance(
	ctx common.Context,
	isBuiltin bool,
	name string,
	parent *PackageInstance,
	attributes map[string]common.Type,
) *PackageInstance {
	return &PackageInstance{
		ctx:       ctx,
		IsBuiltin: isBuiltin,
		Name:      name,
		Parent:    parent,
		CommonInstance: CommonInstance{
			Object: Object{
				typeName:    common.PackageTypeName,
				Attributes:  attributes,
				callHandler: nil,
			},
			prototype: Package,
		},
	}
}

func (p PackageInstance) String(common.State) (string, error) {
	return fmt.Sprintf("<пакет '%s'>", p.Name), nil
}

func (p PackageInstance) Representation(state common.State) (string, error) {
	return p.String(state)
}

func (p PackageInstance) AsBool(common.State) (bool, error) {
	return true, nil
}

func (p PackageInstance) SetAttribute(name string, value common.Type) error {
	if p.IsBuiltin {
		if p.HasAttribute(name) {
			return util.AttributeIsReadOnlyError(p.GetTypeName(), name)
		}

		return util.AttributeNotFoundError(p.GetTypeName(), name)
	}

	return p.CommonInstance.SetAttribute(name, value)
}

func (p *PackageInstance) GetContext() common.Context {
	return p.ctx
}

func (p *PackageInstance) SetContext(ctx common.Context) {
	p.ctx = ctx
}

func comparePackages(_ common.State, op common.Operator, self common.Type, other common.Type) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case *PackageInstance:
		return -2, util.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, util.OperatorNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func NewPackageClass() *Class {
	initAttributes := func() map[string]common.Type {
		return MergeAttributes(
			MakeLogicalOperators(Package),
			MakeComparisonOperators(Package, comparePackages),
			MakeCommonOperators(Package),
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
