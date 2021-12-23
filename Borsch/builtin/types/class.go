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
	getEmptyInstance func() (Type, error),
) *Class {
	return &Class{
		Object:           *newClassObject(name, package_, attributes, doc),
		typeHash:         hashObject(name),
		GetEmptyInstance: getEmptyInstance,
	}
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
			attributes[ops.DocAttributeName] = NilInstance{}
		}
	}

	attributes[ops.PackageAttributeName] = package_
	object.Attributes = attributes
	return object
}

// type ClassInstance struct {
// 	Object
// 	class *class
//
// 	address string
// }
//
// func NewClassInstance(class *class) *ClassInstance {
// 	instance := &ClassInstance{
// 		class: class,
// 	}
// 	instance.address = fmt.Sprintf("%p", instance)
// 	return instance
// }
//
// func (i ClassInstance) String() string {
// 	if attribute, err := i.GetAttribute("__рядок__"); err == nil {
// 		switch __str__ := attribute.(type) {
// 		case CallableType:
// 			result, _ := __str__.Call([]Type{i}, map[string]Type{"я": i})
// 			return result.String()
// 		}
// 	}
//
// 	return fmt.Sprintf("<об'єкт %s з адресою %s>", i.GetTypeName(), i.address)
// }
//
// // Representation TODO: поміняти __рядок__ на __представлення__
// func (i ClassInstance) Representation() string {
// 	return i.String()
// }
//
// func (i ClassInstance) GetTypeHash() uint64 {
// 	return i.class.GetTypeHash()
// }
//
// func (i ClassInstance) GetTypeName() string {
// 	return i.class.GetTypeName()
// }
//
// func (i ClassInstance) AsBool() bool {
// 	return i.class.AsBool()
// }
//
// func (i ClassInstance) GetAttribute(name string) (Type, error) {
// 	return i.class.GetAttribute(name)
// }
//
// func (i ClassInstance) SetAttribute(name string, value Type) (Type, error) {
// 	return i.class.SetAttribute(name, value)
// }
