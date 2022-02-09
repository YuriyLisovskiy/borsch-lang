package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type Class struct {
	Name       string
	attributes map[string]common.Value

	IsFinal bool

	Class *Class
	Bases []*Class

	Parent common.Value

	AttrInitializer  func(*map[string]common.Value)
	GetEmptyInstance func() (common.Value, error)
}

func (c *Class) Setup() {
	c.Class = TypeClass
	if c.GetEmptyInstance == nil {
		c.GetEmptyInstance = func() (common.Value, error) {
			return NewClassInstance(c, map[string]common.Value{}), nil
		}
	}

	if len(c.Bases) == 0 {
		// TODO: set object as a base Class
		c.Bases = []*Class{}
	}

	c.initializeAttributes()
	if c.attributes == nil {
		c.attributes = map[string]common.Value{}
	}
}

func (c *Class) IsValid() bool {
	if len(c.Name) == 0 {
		return false
	}

	if c.attributes == nil {
		return false
	}

	if c.Class == nil {
		return false
	}

	if c.Parent == nil {
		return false
	}

	if c.GetEmptyInstance == nil {
		return false
	}

	return true
}

func (c *Class) GetName() string {
	return c.Name
}

func (c *Class) GetTypeName() string {
	if c.isType() {
		return c.Name
	}

	return c.GetClass().GetTypeName()
}

func (c *Class) IsFinalClass() bool {
	return c.IsFinal
}

func (c *Class) GetClass() *Class {
	if c.Class == nil {
		panic("class is nil")
	}

	return c.Class
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

		return nil, utilities.AttributeNotFoundError(c.GetTypeName(), name)
	}

	return c.GetClass().GetAttribute(name)
}

// GetAttribute uses getAttribute and in case of failure, searches for
// an attribute in TypeClass.
func (c *Class) GetAttribute(name string) (common.Value, error) {
	if val, err := c.getAttribute(name); err == nil {
		return val, nil
	}

	if !c.isType() {
		if attr, err := c.GetClass().GetAttribute(name); err == nil {
			return attr, nil
		}
	}

	return nil, utilities.AttributeNotFoundError(c.GetName(), name)
}

func (c *Class) SetAttribute(name string, newValue common.Value) error {
	if c.isType() {
		if c.HasAttribute(name) {
			return utilities.AttributeIsReadOnlyError(c.GetTypeName(), name)
		}

		return utilities.AttributeNotFoundError(c.GetTypeName(), name)
	}

	if oldValue, ok := c.attributes[name]; ok {
		oldValueClass := oldValue.(ObjectInstance).GetClass()
		newValueClass := newValue.(ObjectInstance).GetClass()
		if oldValueClass == newValueClass || newValueClass.HasBase(oldValueClass) {
			c.attributes[name] = newValue
			return nil
		}

		return utilities.RuntimeError(
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

func (c *Class) EqualsTo(other common.Value) bool {
	cls, ok := other.(*Class)
	return ok && cls == c
}

func (c *Class) HasBase(cls *Class) bool {
	for _, base := range c.Bases {
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
		return nil, utilities.ObjectIsNotCallable(c.GetName(), c.GetTypeName())
	}

	return CallAttribute(state, c, operator, common.ConstructorName, args, kwargs, true)
}

// getAttribute searches for attribute only in current attributes and
// in Bases.
func (c *Class) getAttribute(name string) (common.Value, error) {
	if c.attributes != nil {
		if val, ok := c.attributes[name]; ok {
			return val, nil
		}
	}

	basesLastIdx := len(c.Bases) - 1
	for i := basesLastIdx; i >= 0; i-- {
		if attr, err := c.Bases[i].getAttribute(name); err == nil {
			return attr, nil
		}
	}

	return nil, utilities.AttributeNotFoundError(c.GetName(), name)
}

// isType checks if address of current Class is equal to TypeClass.
func (c *Class) isType() bool {
	return c.Class == c
}

func (c *Class) initializeAttributes() {
	if c.AttrInitializer != nil {
		c.AttrInitializer(&c.attributes)
		c.AttrInitializer = nil
	}

	if _, ok := c.attributes[common.ConstructorName]; !ok {
		// TODO: add doc
		c.attributes[common.ConstructorName] = makeDefaultConstructor(c, "")
	}
}
