package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ObjectBase struct {
	typeName       string
	Attributes     map[string]common.Value
	callHandler    func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error)
	initAttributes AttributesInitializer
}

func NewObjectBase(
	typeName string,
	attributes map[string]common.Value,
	callHandler func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error),
) *ObjectBase {
	return &ObjectBase{
		typeName:       typeName,
		Attributes:     attributes,
		callHandler:    callHandler,
		initAttributes: nil,
	}
}

func (o ObjectBase) GetTypeName() string {
	return o.typeName
}

func (o ObjectBase) GetAttribute(name string) (common.Value, error) {
	if name == common.AttributesName {
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

func (o ObjectBase) SetAttribute(name string, value common.Value) error {
	if name == common.AttributesName {
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
		if val.(ObjectInstance).GetPrototype() == value.(ObjectInstance).GetPrototype() {
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

func (o ObjectBase) HasAttribute(name string) bool {
	if name == common.AttributesName {
		return true
	}

	_, ok := o.Attributes[name]
	return ok
}

func (o *ObjectBase) Call(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (
	common.Value,
	error,
) {
	if o.callHandler != nil {
		return o.callHandler(state, args, kwargs)
	}

	return nil, util.ObjectIsNotCallable(o.GetTypeName(), o.GetTypeName())
}

func (o ObjectBase) Copy() ObjectBase {
	object := ObjectBase{
		typeName:    o.typeName,
		Attributes:  map[string]common.Value{},
		callHandler: o.callHandler,
	}
	for k, v := range o.Attributes {
		object.Attributes[k] = v
	}

	return object
}

func (o ObjectBase) makeAttributes() (DictionaryInstance, error) {
	dict := NewDictionaryInstance()
	for key, val := range o.Attributes {
		err := dict.SetElement(NewStringInstance(key), val)
		if err != nil {
			return DictionaryInstance{}, err
		}
	}

	return dict, nil
}
