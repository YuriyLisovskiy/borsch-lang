package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
)

// Class TODO: Improve
type Class struct {
	Object

	typeHash         uint64
	GetEmptyInstance func() (Type, error)
}

func NewClass(
	name string,
	package_ *PackageInstance,
	attributes map[string]Type,
	doc string,
) *Class {
	class := &Class{
		Object:   *newClassObject(name, package_, attributes, doc),
		typeHash: hashObject(name),
	}
	class.GetEmptyInstance = func() (Type, error) {
		// TODO: set default attributes
		return NewClassInstance(class, map[string]Type{}), nil
	}

	return class
}

func NewBuiltinClass(
	typeHash uint64,
	package_ *PackageInstance,
	attributes map[string]Type,
	doc string,
	getEmptyInstance func() (Type, error),
) *Class {
	return &Class{
		Object:           *newClassObject(GetTypeName(typeHash), package_, attributes, doc),
		typeHash:         typeHash,
		GetEmptyInstance: getEmptyInstance,
	}
}

func (t Class) String() string {
	return fmt.Sprintf("<клас '%s'>", t.GetTypeName())
}

func (t Class) Representation() string {
	return t.String()
}

func (t Class) GetTypeHash() uint64 {
	return t.typeHash
}

func (t Class) AsBool() bool {
	return true
}

// SetAttribute TODO: якщо атрибут не існує, встановити.
//  Якщо атрибут існує, перевірити його тип і, якщо типи співпадають
//  встановити, інакше помилка.
func (t Class) SetAttribute(name string, value Type) (Type, error) {
	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t Class) GetClass() *Class {
	return TypeClass
}

func newClassObject(name string, package_ *PackageInstance, attributes map[string]Type, doc string) *Object {
	object := &Object{
		typeName:    name,
		Attributes:  nil,
		callHandler: nil,
	}

	if constructor, ok := attributes[ops.ConstructorName]; ok {
		switch handler := constructor.(type) {
		case CallableType:
			object.callHandler = handler.Call
		}
	}

	if object.callHandler == nil {
		// TODO: set handler which returns class instance!
	}

	if _, ok := attributes[ops.DocAttributeName]; !ok {
		if len(doc) > 0 {
			attributes[ops.DocAttributeName] = NewStringInstance(doc)
		} else {
			attributes[ops.DocAttributeName] = NewNilInstance()
		}
	}

	attributes[ops.PackageAttributeName] = package_
	object.Attributes = attributes
	return object
}

type ClassInstance struct {
	Object
	class   *Class
	Address string
}

func NewClassInstance(class *Class, attributes map[string]Type) *ClassInstance {
	instance := &ClassInstance{
		Object: Object{
			typeName:    class.GetTypeName(),
			Attributes:  attributes,
			callHandler: nil,
		},
		class: class,
	}
	instance.Address = fmt.Sprintf("%p", instance)
	return instance
}

func (i ClassInstance) String() string {
	if attribute, err := i.GetAttribute("__рядок__"); err == nil {
		switch __str__ := attribute.(type) {
		case *FunctionInstance:
			args := []Type{i}
			kwargs := map[string]Type{__str__.Arguments[0].Name: i}
			if err := CheckFunctionArguments(__str__, &args, &kwargs); err == nil {
				result, err := __str__.Call(&args, &kwargs)
				if err == nil {
					return result.String()
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
func (i ClassInstance) Representation() string {
	return i.String()
}

func (i ClassInstance) GetTypeHash() uint64 {
	return i.class.GetTypeHash()
}

func (i ClassInstance) AsBool() bool {
	if attribute, err := i.GetAttribute("__логічний__"); err == nil {
		switch __bool__ := attribute.(type) {
		case *FunctionInstance:
			args := []Type{i}
			kwargs := map[string]Type{"я": i}
			if err := CheckFunctionArguments(__bool__, &args, &kwargs); err == nil {
				result, err := __bool__.Call(&args, &kwargs)
				if err == nil {
					switch boolResult := result.(type) {
					case BoolInstance:
						return boolResult.AsBool()
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

func (i ClassInstance) SetAttribute(name string, value Type) (Type, error) {
	err := i.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i ClassInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := i.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return i.GetClass().GetAttribute(name)
}

func (i ClassInstance) GetClass() *Class {
	return i.class
}

func (i ClassInstance) Copy() *ClassInstance {
	instance := ClassInstance{
		Object: i.Object.Copy(),
		class:  i.class,
	}
	instance.Address = fmt.Sprintf("%p", &instance)
	return &instance
}
