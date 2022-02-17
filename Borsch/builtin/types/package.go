package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
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
	attributes map[string]common.Object,
) *PackageInstance {
	instance := &PackageInstance{
		ClassInstance: *NewClassInstance(PackageClass, attributes),
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

func (p *PackageInstance) SetAttributes(attrs map[string]common.Object) {
	p.attributes = attrs
	if p.attributes == nil {
		p.attributes = map[string]common.Object{}
	}
}

func comparePackages(_ common.State, op common.Operator, self common.Object, other common.Object) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case *PackageInstance:
		return -2, utilities.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, utilities.OperatorNotSupportedError(op, self, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func NewPackageClass() *Class {
	return &Class{
		Name:    common.PackageTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Object) {
			*attrs = MergeAttributes(
				MakeLogicalOperators(PackageClass),
				MakeComparisonOperators(PackageClass, comparePackages),
				MakeCommonOperators(PackageClass),
			)
		},
		GetEmptyInstance: func() (common.Object, error) {
			panic("unreachable")
		},
	}
}
