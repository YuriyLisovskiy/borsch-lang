package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ObjectInstance interface {
	GetPrototype() *Class
}

type AttributesInitializer func() map[string]common.Type

type CommonInstance struct {
	Object
	prototype *Class
}

func (o CommonInstance) GetTypeName() string {
	return o.GetPrototype().GetTypeName()
}

func (o CommonInstance) GetPrototype() *Class {
	return o.prototype
}

func (o CommonInstance) GetAttribute(name string) (common.Type, error) {
	if attribute, err := o.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	if proto := o.GetPrototype(); proto != nil {
		return proto.GetAttribute(name)
	}

	return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
}

func (o CommonInstance) HasAttribute(name string) bool {
	if o.Object.HasAttribute(name) {
		return true
	}

	if proto := o.GetPrototype(); proto != nil {
		return proto.HasAttribute(name)
	}

	return false
}

func (o CommonInstance) Copy() CommonInstance {
	return CommonInstance{
		Object:    o.Object.Copy(),
		prototype: o.prototype,
	}
}

type BuiltinInstance struct {
	CommonInstance
}

func (o BuiltinInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
	}

	return o.CommonInstance.GetAttribute(name)
}

func (o BuiltinInstance) SetAttribute(name string, _ common.Type) error {
	if name == ops.AttributesName {
		return util.AttributeNotFoundError(o.GetTypeName(), name)
	}

	if o.HasAttribute(name) {
		return util.AttributeIsReadOnlyError(o.GetTypeName(), name)
	}

	return util.AttributeNotFoundError(o.GetTypeName(), name)
}

var (
	Any        *Class = nil
	TypeClass  *Class = nil
	Nil        *Class = nil
	Bool       *Class = nil
	Dictionary *Class = nil
	Function   *Class = nil
	Integer    *Class = nil
	List       *Class = nil
	Package    *Class = nil
	Real       *Class = nil
	String     *Class = nil
)

var BuiltinPackage *PackageInstance

func Init() {
	BuiltinPackage = NewPackageInstance(nil, true, "вбудований", nil, map[string]common.Type{})

	TypeClass = newTypeClass()
	Nil = newNilClass()
	Bool = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()

	TypeClass.InitAttributes()
	Nil.InitAttributes()
	Bool.InitAttributes()
	Dictionary.InitAttributes()
	Function.InitAttributes()
	Integer.InitAttributes()
	List.InitAttributes()
	Package.InitAttributes()
	Real.InitAttributes()
	String.InitAttributes()
}
