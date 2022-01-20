package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Object struct {
	typeName    string
	Attributes  map[string]common.Type
	callHandler func(common.Context, *[]common.Type, *map[string]common.Type) (common.Type, error)
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

func (o Object) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		dict, err := getAttributes(o.Attributes)
		if err != nil {
			return nil, err
		}

		return dict, nil
	}

	if o.Attributes != nil {
		if val, ok := o.Attributes[name]; ok {
			return val, nil
		}
	}

	return nil, util.AttributeNotFoundError(o.GetTypeName(), name)
}

func (o Object) SetAttribute(name string, value common.Type) error {
	if name == ops.AttributesName {
		return util.RuntimeError(
			fmt.Sprintf(
				"неможливо записати значення у атрибут '%s', що призначений лише для читання",
				name,
			),
		)
	}

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

func (o *Object) Call(ctx common.Context, args *[]common.Type, kwargs *map[string]common.Type) (
	common.Type,
	error,
) {
	if o.callHandler != nil {
		return o.callHandler(ctx, args, kwargs)
	}

	return nil, util.ObjectIsNotCallable(o.GetTypeName(), o.GetTypeName())
}

func (o Object) Copy() Object {
	object := Object{
		typeName:   o.typeName,
		Attributes: map[string]common.Type{},
	}
	for k, v := range o.Attributes {
		object.Attributes[k] = v
	}

	return object
}
