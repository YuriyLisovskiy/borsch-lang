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

type CommonObject struct {
	Object
	prototype *Class
}

func (o CommonObject) GetTypeName() string {
	return o.GetPrototype().GetTypeName()
}

func (o CommonObject) GetPrototype() *Class {
	return o.prototype
}

func (o CommonObject) GetAttribute(name string) (common.Type, error) {
	if attribute, err := o.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return o.GetPrototype().GetAttribute(name)
}

func (o CommonObject) HasAttribute(name string) bool {
	return o.Object.HasAttribute(name) || o.GetPrototype().HasAttribute(name)
}

func (o CommonObject) Copy() CommonObject {
	return CommonObject{
		Object:    o.Object.Copy(),
		prototype: o.prototype,
	}
}

type BuiltinObject struct {
	CommonObject
}

func (o BuiltinObject) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
	}

	return o.CommonObject.GetAttribute(name)
}

func (o BuiltinObject) SetAttribute(name string, _ common.Type) error {
	if name == ops.AttributesName {
		return util.AttributeNotFoundError(o.GetTypeName(), name)
	}

	if o.Object.HasAttribute(name) || o.GetPrototype().HasAttribute(name) {
		return util.AttributeIsReadOnlyError(o.GetTypeName(), name)
	}

	return util.AttributeNotFoundError(o.GetTypeName(), name)
}

var (
	Bool       *Class = nil
	Dictionary *Class = nil
	Function   *Class = nil
	Integer    *Class = nil
	List       *Class = nil
	Nil        *Class = nil
	Package    *Class = nil
	Real       *Class = nil
	String     *Class = nil
	TypeClass  *Class = nil
	Any        *Class = nil
)

var BuiltinPackage *PackageInstance

func Init() {
	BuiltinPackage = NewPackageInstance(nil, true, "вбудований", nil, map[string]common.Type{})

	Bool = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Nil = newNilClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()
	TypeClass = newTypeClass()

	Bool.InitAttributes()
	Dictionary.InitAttributes()
	Function.InitAttributes()
	Integer.InitAttributes()
	List.InitAttributes()
	Nil.InitAttributes()
	Package.InitAttributes()
	Real.InitAttributes()
	String.InitAttributes()
	TypeClass.InitAttributes()
}
