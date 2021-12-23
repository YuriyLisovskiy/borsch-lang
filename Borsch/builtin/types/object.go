package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Object struct {
	typeName    string
	Attributes  map[string]Type
	callHandler func(*[]Type, *map[string]Type) (Type, error)
}

func (o Object) makeAttributes() (DictionaryInstance, error) {
	dict := NewDictionaryInstance()
	for key, val := range o.Attributes {
		err := dict.SetElement(NewStringInstance(key), val)
		if err != nil {
			return DictionaryInstance{}, err
		}
	}

	return dict, nil
}

func (o Object) GetTypeName() string {
	return o.typeName
}

func (o Object) GetAttribute(name string) (Type, error) {
	if o.Attributes != nil {
		if name == "__атрибути__" {
			dict, err := o.makeAttributes()
			if err != nil {
				return nil, err
			}

			return dict, nil
		}

		if val, ok := o.Attributes[name]; ok {
			return val, nil
		}
	}

	return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
}

func (o Object) SetAttribute(name string, value Type) error {
	if o.Attributes == nil {
		return util.AttributeNotFoundError(o.GetTypeName(), name)
	}

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

func (o *Object) Call(args *[]Type, kwargs *map[string]Type) (Type, error) {
	if o.callHandler != nil {
		return o.callHandler(args, kwargs)
	}

	return nil, util.ObjectIsNotCallable(o.GetTypeName(), o.GetTypeName())
}

func (o Object) Copy() Object {
	object := Object{
		typeName: o.typeName,
	}
	for k, v := range o.Attributes {
		object.Attributes[k] = v
	}

	return object
}
