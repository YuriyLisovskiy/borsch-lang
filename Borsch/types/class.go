package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Class struct {
	Object

	prototype        *Class
	attrInitializer  AttributesInitializer
	GetEmptyInstance func() (common.Type, error)
}

func NewClass(
	name string,
	package_ *PackageInstance,
	initAttributes func() map[string]common.Type,
	doc string,
) *Class {
	class := &Class{
		Object:    *newClassObject(name, package_, initAttributes, doc),
		prototype: TypeClass,
	}
	class.GetEmptyInstance = func() (common.Type, error) {
		// TODO: set default attributes
		return NewClassInstance(class, map[string]common.Type{}), nil
	}
	return class
}

func NewBuiltinClass(
	typeName string,
	package_ *PackageInstance,
	initAttributes func() map[string]common.Type,
	doc string,
	getEmptyInstance func() (common.Type, error),
) *Class {
	return &Class{
		Object:           *newClassObject(typeName, package_, initAttributes, doc),
		prototype:        TypeClass,
		GetEmptyInstance: getEmptyInstance,
	}
}

func (c Class) String(common.State) (string, error) {
	return fmt.Sprintf("<клас '%s'>", c.GetTypeName()), nil
}

func (c Class) Representation(state common.State) (string, error) {
	return c.String(state)
}

func (c Class) AsBool(common.State) (bool, error) {
	return true, nil
}

// func (c *Class) GetTypeName() string {
// 	if c.prototype == c {
// 		return c.Object.GetTypeName()
// 	}
//
// 	return c.prototype.GetTypeName()
// }

func (c *Class) GetPrototype() *Class {
	if c.prototype == nil {
		panic("Class: prototype is nil")
	}

	if c.prototype == c {
		return c
	}

	return c.prototype
}

func (c *Class) GetOperator(name string) (common.Type, error) {
	if c.isType() {
		return c.Object.GetAttribute(name)
	}

	return c.prototype.GetAttribute(name)
}

func (c *Class) GetAttribute(name string) (common.Type, error) {
	if c.isType() {
		return c.Object.GetAttribute(name)
	}

	if attr, err := c.Object.GetAttribute(name); err == nil {
		return attr, nil
	}

	if attr, err := c.prototype.GetAttribute(name); err == nil {
		return attr, nil
	}

	return nil, util.AttributeNotFoundError(c.GetTypeName(), name)
}

func (c *Class) SetAttribute(name string, value common.Type) error {
	if c.isType() {
		if c.HasAttribute(name) {
			return util.AttributeIsReadOnlyError(c.GetTypeName(), name)
		}

		return util.AttributeNotFoundError(c.GetTypeName(), name)
	}

	return c.Object.SetAttribute(name, value)
}

func (c *Class) HasAttribute(name string) bool {
	if c.isType() {
		return c.Object.HasAttribute(name)
	}

	if !c.Object.HasAttribute(name) {
		return c.prototype.HasAttribute(name)
	}

	return true
}

func (c *Class) InitAttributes() {
	if c.Object.initAttributes != nil {
		c.Attributes = c.Object.initAttributes()
		c.Object.initAttributes = nil
	}
}

func (c *Class) EqualsTo(other common.Type) bool {
	switch right := other.(type) {
	case *Class:
		return c == right
	default:
		return false
	}
}

func (c *Class) isType() bool {
	return c.prototype == c
}

func newClassObject(
	typeName string,
	package_ *PackageInstance,
	attrInitializer AttributesInitializer,
	doc string,
) *Object {
	object := &Object{
		typeName:    typeName,
		Attributes:  nil,
		callHandler: nil,
	}

	object.initAttributes = func() map[string]common.Type {
		attributes := attrInitializer()
		if constructor, ok := attributes[common.ConstructorName]; ok {
			switch handler := constructor.(type) {
			case common.CallableType:
				object.callHandler = handler.Call
			}
		}

		if _, ok := attributes[common.DocAttributeName]; !ok {
			if len(doc) > 0 {
				attributes[common.DocAttributeName] = NewStringInstance(doc)
			} else {
				attributes[common.DocAttributeName] = NewNilInstance()
			}
		}

		attributes[common.PackageAttributeName] = package_
		return attributes
	}

	if object.callHandler == nil {
		// TODO: set handler which returns class instance!
	}

	return object
}

type ClassInstance struct {
	CommonInstance
	class   *Class
	Address string
}

func NewClassInstance(class *Class, attributes map[string]common.Type) *ClassInstance {
	instance := &ClassInstance{
		CommonInstance: CommonInstance{
			Object: Object{
				typeName:    class.GetTypeName(),
				Attributes:  attributes,
				callHandler: nil,
			},
			prototype: class,
		},
		class: class,
	}
	instance.Address = fmt.Sprintf("%p", instance)
	return instance
}

func (i ClassInstance) String(state common.State) (string, error) {
	if i.HasAttribute(common.StringOperatorName) {
		result, err := CallByName(state, i, common.StringOperatorName, nil, nil, true)
		if err != nil {
			return "", err
		}

		return result.String(state)
	}

	return fmt.Sprintf("<об'єкт %s з адресою %s>", i.GetTypeName(), i.Address), nil
}

// Representation TODO: поміняти __рядок__ на __представлення__
func (i ClassInstance) Representation(state common.State) (string, error) {
	return i.String(state)
}

func (i ClassInstance) AsBool(state common.State) (bool, error) {
	if !i.HasAttribute(common.BoolOperatorName) {
		return true, nil
	}

	boolOperator, _ := i.GetOperator(common.BoolOperatorName)
	result, err := CallAttribute(state, i, boolOperator, common.BoolOperatorName, nil, nil, true)
	if err != nil {
		return false, err
	}

	return result.AsBool(state)
}

func (i ClassInstance) Copy() *ClassInstance {
	instance := ClassInstance{
		CommonInstance: i.CommonInstance.Copy(),
		class:          i.class,
	}
	instance.Address = fmt.Sprintf("%p", &instance)
	return &instance
}
