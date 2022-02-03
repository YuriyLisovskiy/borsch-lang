package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ObjectInstance interface {
	GetPrototype() *Class
}

type AttributesInitializer func() map[string]common.Value

type CommonInstance struct {
	ObjectBase
	prototype *Class
}

func NewCommonInstance(object *ObjectBase, prototype *Class) *CommonInstance {
	return &CommonInstance{
		ObjectBase: *object,
		prototype:  prototype,
	}
}

func (o CommonInstance) GetTypeName() string {
	return o.ObjectBase.GetTypeName()
}

func (o CommonInstance) GetPrototype() *Class {
	if o.prototype == nil {
		panic("CommonInstance: prototype is nil")
	}

	return o.prototype
}

func (o CommonInstance) GetOperator(name string) (common.Value, error) {
	if common.IsOperator(name) {
		if attr, err := o.GetPrototype().GetAttribute(name); err == nil {
			return attr, nil
		}
	}

	return nil, util.OperatorNotFoundError(o.GetTypeName(), name)
}

func (o CommonInstance) GetAttribute(name string) (common.Value, error) {
	if attr, err := o.ObjectBase.GetAttribute(name); err == nil {
		return attr, nil
	}

	if attr, err := o.GetPrototype().GetAttribute(name); err == nil {
		return attr, nil
	}

	return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
}

func (o CommonInstance) HasAttribute(name string) bool {
	if o.ObjectBase.HasAttribute(name) {
		return true
	}

	return o.GetPrototype().HasAttribute(name)
}

func (o CommonInstance) Copy() CommonInstance {
	return CommonInstance{
		ObjectBase: o.ObjectBase.Copy(),
		prototype:  o.prototype,
	}
}

type BuiltinInstance struct {
	CommonInstance
}

func (o BuiltinInstance) GetAttribute(name string) (common.Value, error) {
	if name == common.AttributesName {
		return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
	}

	return o.CommonInstance.GetAttribute(name)
}

func (o BuiltinInstance) SetAttribute(name string, _ common.Value) error {
	if name == common.AttributesName {
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
	BuiltinPackage = NewPackageInstance(nil, true, "вбудований", nil, map[string]common.Value{})

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
