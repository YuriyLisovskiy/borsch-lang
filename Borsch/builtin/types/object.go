package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Object struct {
	typeHash   uint64
	typeName   string
	Attributes map[string]Type
}

func newBuiltinObject(typeHash uint64, attributes map[string]Type) *Object {
	return &Object{
		typeHash:   typeHash,
		typeName:   GetTypeName(typeHash),
		Attributes: attributes,
	}
}

func newObject(typeName string, attributes map[string]Type) *Object {
	return &Object{
		typeHash:   hashObject(typeName),
		typeName:   typeName,
		Attributes: attributes,
	}
}

func (o Object) makeAttributes() (*DictionaryType, error) {
	dict := NewDictionaryType()
	for key, val := range o.Attributes {
		err := dict.SetElement(NewStringType(key), val)
		if err != nil {
			return nil, err
		}
	}

	return dict, nil
}

func (o Object) GetTypeHash() uint64 {
	return o.typeHash
}

func (o Object) GetTypeName() string {
	return o.typeName
}

func (o Object) GetAttribute(name string) (Type, error) {
	if name == "__атрибути__" {
		return o.makeAttributes()
	}

	if val, ok := o.Attributes[name]; ok {
		return val, nil
	}

	return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
}

func (o Object) SetAttribute(name string, value Type) error {
	if val, ok := o.Attributes[name]; ok {
		if val.GetTypeHash() == value.GetTypeHash() {
			o.Attributes[name] = value
			return nil
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
				value.GetTypeName(), name, val.GetTypeName(),
			),
		)
	}

	o.Attributes[name] = value
	return nil
}

func (o Object) HasAttribute(name string) bool {
	_, ok := o.Attributes[name]
	return ok
}
