package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
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

func (c Class) AsBool(common.State) bool {
	return true
}

func (c *Class) GetTypeName() string {
	if c.prototype == c {
		return c.Object.GetTypeName()
	}

	return c.prototype.GetTypeName()
}

func (c *Class) GetPrototype() *Class {
	if c.prototype == c {
		return c
	}

	return c.prototype
}

func (c *Class) GetAttribute(name string) (common.Type, error) {
	if c.isType() {
		if name == ops.AttributesName {
			return nil, util.AttributeNotFoundError(c.GetTypeName(), name)
		}

		return c.Object.GetAttribute(name)
	}

	if attr, err := c.prototype.GetAttribute(name); err == nil {
		return attr, nil
	}

	return c.Object.GetAttribute(name)
}

func (c *Class) SetAttribute(name string, value common.Type) error {
	if c.isType() {
		if name == ops.AttributesName {
			return util.AttributeNotFoundError(c.GetTypeName(), name)
		}

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

	if !c.prototype.HasAttribute(name) {
		return c.Object.HasAttribute(name)
	}

	return false
}

func (c *Class) InitAttributes() {
	if c.Object.initAttributes != nil {
		c.Attributes = c.Object.initAttributes()
		c.Object.initAttributes = nil
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
		if constructor, ok := attributes[ops.ConstructorName]; ok {
			switch handler := constructor.(type) {
			case common.CallableType:
				object.callHandler = handler.Call
			}
		}

		if _, ok := attributes[ops.DocAttributeName]; !ok {
			if len(doc) > 0 {
				attributes[ops.DocAttributeName] = NewStringInstance(doc)
			} else {
				attributes[ops.DocAttributeName] = NewNilInstance()
			}
		}

		attributes[ops.PackageAttributeName] = package_
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
	if i.HasAttribute(ops.StringOperatorName) {
		result, err := CallByName(state, i, ops.StringOperatorName, nil, nil, true)
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

func (i ClassInstance) AsBool(state common.State) bool {
	if i.HasAttribute(ops.BoolOperatorName) {
		result, err := CallByName(state, i, ops.BoolOperatorName, nil, nil, true)
		if err != nil {
			// TODO: return error instead of panic
			panic(err)
		}

		return result.AsBool(state)
	}

	panic(ops.BoolOperatorName + " not found")

	// if attribute, err := i.GetAttribute(ops.BoolOperatorName); err == nil {
	// 	switch __bool__ := attribute.(type) {
	// 	case *FunctionInstance:
	// 		args := []common.Type{i}
	// 		kwargs := map[string]common.Type{__bool__.Parameters[0].Name: i}
	// 		if err := CheckFunctionArguments(state, __bool__, &args, &kwargs); err == nil {
	// 			result, err := __bool__.Call(state, &args, &kwargs)
	// 			if err == nil {
	// 				switch boolResult := result.(type) {
	// 				case BoolInstance:
	// 					return boolResult.AsBool(state)
	// 				}
	// 			} else {
	// 				// TODO: return error
	// 			}
	// 		} else {
	// 			// TODO: return error
	// 		}
	// 	}
	// }
	//
	// // TODO: return error
	// return false
}

func (i ClassInstance) Copy() *ClassInstance {
	instance := ClassInstance{
		CommonInstance: i.CommonInstance.Copy(),
		class:          i.class,
	}
	instance.Address = fmt.Sprintf("%p", &instance)
	return &instance
}
