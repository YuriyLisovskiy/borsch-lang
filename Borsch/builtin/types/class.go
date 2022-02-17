package types

import (
	"errors"
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type NewFunc func(cls *Class, args List, state common.State) (common.Object, error)

type ConstructFunc func(self common.Object, args List, state common.State) error

const (
	// TPFLAGS_BASETYPE is set if the type allows subclassing.
	TPFLAGS_BASETYPE uint = 1 << 10

	// TPFLAGS_READY is set if the type is 'ready' -- fully initialized.
	TPFLAGS_READY uint = 1 << 12

	// TPFLAGS_READYING is set while the type is being 'readied', to prevent
	// recursive ready calls.
	TPFLAGS_READYING uint = 1 << 13
)

type Class struct {
	Name string
	Doc  string
	Dict map[string]common.Object

	IsFinal bool

	ObjectClass *Class
	Bases       *List

	New       NewFunc
	Construct ConstructFunc
	Flags     uint
}

var TypeClass = &Class{
	Name: "тип",
	Doc:  "тип(об_єкт) -> тип об'єкта\nтип(назва, бази, атрибути) -> новий тип",
	Dict: StringDict{},
}

var ObjectClass = &Class{
	Name: "об_єкт",
	Doc:  "Базовий тип",
	Dict: StringDict{},
}

func init() {
	// Initialised like this to avoid initialisation loops
	TypeClass.New = TypeNew
	TypeClass.Construct = TypeInit
	TypeClass.ObjectClass = TypeClass
	ObjectClass.New = ObjectNew
	ObjectClass.Construct = ObjectConstruct
	ObjectClass.ObjectClass = TypeClass
	err := TypeClass.Ready()
	if err != nil {
		log.Fatal(err)
	}

	err = ObjectClass.Ready()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Class) Ready() error {
	if c.Flags&TPFLAGS_READY != 0 {
		if c.Dict == nil {
			return ErrorNewf(SystemError, "Type.Ready is Ready but Dict is nil")
		}

		return nil
	}

	if c.Flags&TPFLAGS_READYING != 0 {
		return ErrorNewf(SystemError, "Type.Ready already readying")
	}

	c.Flags |= TPFLAGS_READYING

	// Now the only way base can still be nil is if type is
	// ObjectType.

	// Initialize Bases
	if c.Bases == nil {
		c.Bases = &List{}
	}

	// Initialize tp_dict
	dict := c.Dict
	if dict == nil {
		dict = NewStringDict()
		c.Dict = dict
	}

	if _, ok := c.Dict[common.DocAttributeName]; ok {
		if c.Doc != "" {
			c.Dict[common.DocAttributeName] = String(c.Doc)
		} else {
			c.Dict[common.DocAttributeName] = Nil
		}
	}

	// TODO: Link into each base class's list of subclasses
	// bases := c.Bases
	// for i := range bases {
	// }

	// All done -- set the ready flag
	if c.Dict == nil {
		panic("Type.Ready Dict is nil")
	}

	c.Flags = (c.Flags &^ TPFLAGS_READYING) | TPFLAGS_READY
	return nil
}

func (c *Class) Class() *Class {

}

func (c *Class) __str__() (common.Object, error) {
	if res, ok, err := c.CallMethod(common.StringOperator, Tuple{t}, nil); ok {
		return res, err
	}

	return c.__represent__()
}

func (c *Class) __represent__() (common.Object, error) {
	if res, ok, err := c.CallMethod(common.RepresentOperator, List{c}, nil); ok {
		return res, err
	}

	if c.Name == "" {
		// FIXME: not a good way to tell objects from classes!
		return String(fmt.Sprintf("<об'єкт %s з адресою %p>", c.Type().Name, c)), nil
	}

	return String(fmt.Sprintf("<клас '%s'>", c.Name)), nil
}

func (c *Class) CallMethod(name string, args Tuple, state common.State) (Object, bool, error) {
	fn := c.GetAttrOrNil(name)
	if fn == nil {
		return nil, false, nil
	}

	res, err := Call(state, fn, args, nil)
	return res, true, err
}

func (c *Class) GetOperator(name string) (common.Object, error) {
	if c.isType() {
		if c.Dict != nil {
			if val, ok := c.Dict[name]; ok {
				return val, nil
			}
		}

		return nil, utilities.AttributeNotFoundError(c.GetTypeName(), name)
	}

	return c.GetClass().GetAttribute(name)
}

// GetAttribute uses getAttribute and in case of failure, searches for
// an attribute in TypeClass.
func (c *Class) GetAttribute(name string) (common.Object, error) {
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

func (c *Class) SetAttribute(name string, newValue common.Object) error {
	if c.isType() {
		if c.HasAttribute(name) {
			return utilities.AttributeIsReadOnlyError(c.GetTypeName(), name)
		}

		return utilities.AttributeNotFoundError(c.GetTypeName(), name)
	}

	if oldValue, ok := c.Dict[name]; ok {
		oldValueClass := oldValue.(ObjectInstance).GetClass()
		newValueClass := newValue.(ObjectInstance).GetClass()
		if oldValueClass == newValueClass || newValueClass.HasBase(oldValueClass) {
			c.Dict[name] = newValue
			return nil
		}

		return errors.New(
			fmt.Sprintf(
				"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
				newValue.GetTypeName(), name, oldValue.GetTypeName(),
			),
		)
	}

	c.Dict[name] = newValue
	return nil
}

func (c *Class) HasAttribute(name string) bool {
	if _, ok := c.Dict[name]; !ok {
		if !c.isType() {
			return c.GetClass().HasAttribute(name)
		}

		return false
	}

	return true
}

func (c *Class) SetAttributes(attrs map[string]common.Object) {
	c.Dict = attrs
	if c.Dict == nil {
		c.Dict = map[string]common.Object{}
	}
}

func (c *Class) EqualsTo(other common.Object) bool {
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

// Call executes common.ConstructorName operator if it exists in Dict.
func (c *Class) Call(state common.State, args *[]common.Object, kwargs *map[string]common.Object) (
	common.Object,
	error,
) {
	operator, err := c.GetOperator(common.ConstructorName)
	if err != nil {
		return nil, utilities.ObjectIsNotCallable(c.GetName(), c.GetTypeName())
	}

	return CallAttribute(state, c, operator, common.ConstructorName, args, kwargs, true)
}

// getAttribute searches for attribute only in current Dict and
// in Bases.
func (c *Class) getAttribute(name string) (common.Object, error) {
	if c.Dict != nil {
		if val, ok := c.Dict[name]; ok {
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
		c.AttrInitializer(&c.Dict)
		c.AttrInitializer = nil
	}

	if _, ok := c.Dict[common.ConstructorName]; !ok {
		// TODO: add doc
		c.Dict[common.ConstructorName] = makeDefaultConstructor(c, "")
	}
}
