package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Class struct {
	name       string
	attributes map[string]common.Value

	class *Class
	bases []*Class

	parent common.Value

	attrInitializer  func(*map[string]common.Value)
	GetEmptyInstance func() (common.Value, error)
}

func NewClass(
	name string,
	bases []*Class,
	parent common.Value,
	attrInitializer func(*map[string]common.Value),
	getEmptyInstanceFunc func() (common.Value, error),
) *Class {
	class := &Class{
		name:             name,
		attributes:       map[string]common.Value{},
		class:            TypeClass,
		bases:            bases,
		parent:           parent,
		attrInitializer:  attrInitializer,
		GetEmptyInstance: getEmptyInstanceFunc,
	}
	if len(class.bases) == 0 {
		// TODO: set object as a base class
	}

	if class.GetEmptyInstance == nil {
		class.GetEmptyInstance = func() (common.Value, error) {
			return NewClassInstance(class, map[string]common.Value{}), nil
		}
	}

	return class
}

func (c *Class) GetName() string {
	return c.name
}

func (c *Class) GetTypeName() string {
	if c.isType() {
		return c.name
	}

	return c.GetClass().GetTypeName()
}

func (c *Class) GetClass() *Class {
	if c.class == nil {
		panic("Class: class is nil")
	}

	return c.class
}

func (c *Class) GetAddress() string {
	return fmt.Sprintf("%p", c)
}

func (c *Class) String(common.State) (string, error) {
	return fmt.Sprintf("<клас '%s'>", c.GetName()), nil
}

func (c *Class) Representation(state common.State) (string, error) {
	return c.String(state)
}

func (c *Class) AsBool(common.State) (bool, error) {
	return true, nil
}

func (c *Class) GetOperator(name string) (common.Value, error) {
	if c.isType() {
		if c.attributes != nil {
			if val, ok := c.attributes[name]; ok {
				return val, nil
			}
		}

		return nil, util.AttributeNotFoundError(c.GetTypeName(), name)
	}

	return c.GetClass().GetAttribute(name)
}

func (c *Class) GetAttribute(name string) (common.Value, error) {
	if c.attributes != nil {
		if val, ok := c.attributes[name]; ok {
			return val, nil
		}
	}

	basesLastIdx := len(c.bases) - 1
	for i := basesLastIdx; i >= 0; i-- {
		if attr, err := c.bases[i].GetAttribute(name); err == nil {
			return attr, nil
		}
	}

	if !c.isType() {
		if attr, err := c.GetClass().GetAttribute(name); err == nil {
			return attr, nil
		}
	}

	return nil, util.AttributeNotFoundError(c.GetName(), name)
}

func (c *Class) SetAttribute(name string, newValue common.Value) error {
	if c.isType() {
		if c.HasAttribute(name) {
			return util.AttributeIsReadOnlyError(c.GetTypeName(), name)
		}

		return util.AttributeNotFoundError(c.GetTypeName(), name)
	}

	if oldValue, ok := c.attributes[name]; ok {
		oldValueClass := oldValue.(ObjectInstance).GetClass()
		newValueClass := newValue.(ObjectInstance).GetClass()
		if oldValueClass == newValueClass || newValueClass.HasBase(oldValueClass) {
			c.attributes[name] = newValue
			return nil
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
				newValue.GetTypeName(), name, oldValue.GetTypeName(),
			),
		)
	}

	c.attributes[name] = newValue
	return nil
}

func (c *Class) HasAttribute(name string) bool {
	if _, ok := c.attributes[name]; !ok {
		if !c.isType() {
			return c.GetClass().HasAttribute(name)
		}

		return false
	}

	return true
}

func (c *Class) SetAttributes(attrs map[string]common.Value) {
	c.attributes = attrs
	if c.attributes == nil {
		c.attributes = map[string]common.Value{}
	}
}

func (c *Class) InitAttributes() {
	if c.attrInitializer != nil {
		c.attrInitializer(&c.attributes)
		c.attrInitializer = nil
	}

	if _, ok := c.attributes[common.ConstructorName]; !ok {
		// TODO: add doc
		c.attributes[common.ConstructorName] = getDefaultConstructor(c, "")
	}
}

func (c *Class) EqualsTo(other common.Value) bool {
	cls, ok := other.(*Class)
	return ok && cls == c
}

func (c *Class) HasBase(cls *Class) bool {
	for _, base := range c.bases {
		if cls == base {
			return true
		}
	}

	return false
}

// Call executes common.ConstructorName operator if it exists in attributes.
func (c *Class) Call(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
	operator, err := c.GetOperator(common.ConstructorName)
	if err != nil {
		return nil, util.ObjectIsNotCallable(c.GetName(), c.GetTypeName())
	}

	return CallAttribute(state, c, operator, common.ConstructorName, args, kwargs, true)
}

// isType checks if address of current class is equal to TypeClass.
func (c *Class) isType() bool {
	return c.class == c
}
