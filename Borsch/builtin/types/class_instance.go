package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type ClassInstance struct {
	class      *Type
	attributes map[string]common.Value
	address    string
}

func NewClassInstance(class *Type, attributes map[string]common.Value) *ClassInstance {
	if attributes == nil {
		attributes = map[string]common.Value{}
	}

	instance := &ClassInstance{
		class:      class,
		attributes: attributes,
		address:    "",
	}

	instance.address = fmt.Sprintf("%p", instance)
	return instance
}

func (i ClassInstance) GetClass() *Type {
	return i.class
}

func (i ClassInstance) GetAddress() string {
	return i.address
}

func (i ClassInstance) String(state common.State) (string, error) {
	if operator, err := i.GetOperator(common.StringOperator); err == nil {
		result, err := CallAttribute(state, i, operator, common.StringOperator, nil, nil, true)
		if err != nil {
			return "", err
		}

		return result.String(state)
	}

	return fmt.Sprintf("<об'єкт %s з адресою %s>", i.GetTypeName(), i.GetAddress()), nil
}

func (i ClassInstance) Representation(state common.State) (string, error) {
	if operator, err := i.GetOperator(common.RepresentOperator); err == nil {
		result, err := CallAttribute(state, i, operator, common.RepresentOperator, nil, nil, true)
		if err != nil {
			return "", err
		}

		return result.String(state)
	}

	return fmt.Sprintf("<об'єкт %s з адресою %s>", i.GetTypeName(), i.GetAddress()), nil
}

func (i ClassInstance) GetTypeName() string {
	return i.GetClass().GetName()
}

func (i ClassInstance) AsBool(common.State) (bool, error) {
	return true, nil
}

func (i ClassInstance) GetOperator(name string) (common.Value, error) {
	if attr, err := i.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.OperatorNotFoundError(i.GetTypeName(), name)
}

func (i ClassInstance) GetAttribute(name string) (common.Value, error) {
	if val, ok := i.attributes[name]; ok {
		return val, nil
	}

	if attr, err := i.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.AttributeNotFoundError(i.GetTypeName(), name)
}

func (i ClassInstance) SetAttribute(name string, newValue common.Value) error {
	if oldValue, ok := i.attributes[name]; ok {
		oldValueClass := oldValue.(ObjectInstance).GetClass()
		newValueClass := newValue.(ObjectInstance).GetClass()
		if oldValueClass == newValueClass || newValueClass.HasBase(oldValueClass) {
			i.attributes[name] = newValue
			return nil
		}

		return errors.New(
			fmt.Sprintf(
				"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
				newValue.GetTypeName(), name, oldValue.GetTypeName(),
			),
		)
	}

	i.attributes[name] = newValue
	return nil
}

func (i ClassInstance) HasAttribute(name string) bool {
	if _, ok := i.attributes[name]; ok {
		return true
	}

	return i.GetClass().HasAttribute(name)
}

func (i ClassInstance) Call(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (
	common.Value,
	error,
) {
	operator, err := i.GetOperator(common.CallOperatorName)
	if err != nil {
		return nil, utilities.ObjectIsNotCallable("", i.GetTypeName())
	}

	return CallAttribute(state, i, operator, common.CallOperatorName, args, kwargs, true)
}

func (i ClassInstance) Copy() *ClassInstance {
	instance := &ClassInstance{
		class:      i.class,
		attributes: map[string]common.Value{},
		address:    "",
	}

	for k, v := range i.attributes {
		instance.attributes[k] = v
	}

	instance.address = fmt.Sprintf("%p", instance)
	return instance
}
