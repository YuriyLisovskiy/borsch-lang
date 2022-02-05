package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageInstance struct {
	ClassInstance
	Name   string
	Parent *PackageInstance

	ctx common.Context
}

func NewPackageInstance(
	ctx common.Context,
	name string,
	parent *PackageInstance,
	attributes map[string]common.Value,
) *PackageInstance {
	instance := &PackageInstance{
		ClassInstance: *NewClassInstance(Package, attributes),
		Name:          name,
		Parent:        parent,
		ctx:           ctx,
	}

	instance.address = fmt.Sprintf("%p", instance)
	return instance
}

func (p *PackageInstance) String(common.State) (string, error) {
	return fmt.Sprintf("<пакет '%s'>", p.Name), nil
}

func (p *PackageInstance) Representation(state common.State) (string, error) {
	return p.String(state)
}

func (p *PackageInstance) GetContext() common.Context {
	return p.ctx
}

func (p *PackageInstance) SetContext(ctx common.Context) {
	p.ctx = ctx
}

func (p *PackageInstance) SetAttributes(attrs map[string]common.Value) {
	p.attributes = attrs
	if p.attributes == nil {
		p.attributes = map[string]common.Value{}
	}
}

func comparePackages(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
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
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = MergeAttributes(
			MakeLogicalOperators(Package),
			MakeComparisonOperators(Package, comparePackages),
			MakeCommonOperators(Package),
		)
	}

	return NewClass(common.PackageTypeName, nil, BuiltinPackage, initAttributes, nil)
}
