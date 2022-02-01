package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

type Class struct {
	CommonObject

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
		CommonObject: CommonObject{
			Object:    *newClassObject(name, package_, initAttributes, doc),
			prototype: TypeClass,
		},
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
		CommonObject: CommonObject{
			Object:    *newClassObject(typeName, package_, initAttributes, doc),
			prototype: TypeClass,
		},
		GetEmptyInstance: getEmptyInstance,
	}
}

func (c Class) String(common.State) string {
	return fmt.Sprintf("<клас '%s'>", c.GetTypeName())
}

func (c Class) Representation(state common.State) string {
	return c.String(state)
}

func (c Class) AsBool(common.State) bool {
	return true
}

func (c *Class) InitAttributes() {
	if c.Object.initAttributes != nil {
		c.Attributes = c.Object.initAttributes()
		c.Object.initAttributes = nil
	}
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
	CommonObject
	class   *Class
	Address string
}

func NewClassInstance(class *Class, attributes map[string]common.Type) *ClassInstance {
	instance := &ClassInstance{
		CommonObject: CommonObject{
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

func (i ClassInstance) String(state common.State) string {
	if attribute, err := i.GetAttribute("__рядок__"); err == nil {
		switch __str__ := attribute.(type) {
		case *FunctionInstance:
			args := []common.Type{i}
			kwargs := map[string]common.Type{__str__.Arguments[0].Name: i}
			if err := CheckFunctionArguments(state, __str__, &args, &kwargs); err == nil {
				result, err := __str__.Call(state, &args, &kwargs)
				if err == nil {
					return result.String(state)
				} else {
					// TODO: return error
				}
			} else {
				// TODO: return error
			}
		}
	}

	return fmt.Sprintf("<об'єкт %s з адресою %s>", i.GetTypeName(), i.Address)
}

// Representation TODO: поміняти __рядок__ на __представлення__
func (i ClassInstance) Representation(state common.State) string {
	return i.String(state)
}

func (i ClassInstance) AsBool(state common.State) bool {
	if attribute, err := i.GetAttribute(ops.BoolOperatorName); err == nil {
		switch __bool__ := attribute.(type) {
		case *FunctionInstance:
			args := []common.Type{i}
			kwargs := map[string]common.Type{__bool__.Arguments[0].Name: i}
			if err := CheckFunctionArguments(state, __bool__, &args, &kwargs); err == nil {
				result, err := __bool__.Call(state, &args, &kwargs)
				if err == nil {
					switch boolResult := result.(type) {
					case BoolInstance:
						return boolResult.AsBool(state)
					}
				} else {
					// TODO: return error
				}
			} else {
				// TODO: return error
			}
		}
	}

	// TODO: return error
	return false
}

func (i ClassInstance) Copy() *ClassInstance {
	instance := ClassInstance{
		CommonObject: i.CommonObject.Copy(),
		class:        i.class,
	}
	instance.Address = fmt.Sprintf("%p", &instance)
	return &instance
}
