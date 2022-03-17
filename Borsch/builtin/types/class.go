package types

import (
	"errors"
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type NewF func(state common.State, cls *Class, args Tuple, kwargs map[string]common.Value) (common.Value, error)
type ConstructF func(state common.State, self common.Value, args Tuple, kwargs map[string]common.Value) error

type Class struct {
	Name        string
	Doc         string
	Dict        map[string]common.Value
	IsFinal     bool
	ObjectClass *Class
	Bases       []*Class
	New         *FunctionInstance
	Construct   ConstructF
}

var TypeClass = &Class{
	Name: "тип",
	Doc:  "тип(об_єкт) -> тип об'єкта",
	Dict: map[string]common.Value{},
}

var ObjectClass = &Class{
	Name: "об_єкт",
	Doc:  "Базовий тип",
	Dict: map[string]common.Value{},
}

func init() {
	// Initialised like this to avoid initialisation loops
	TypeClass.New = TypeNew
	TypeClass.Construct = TypeConstruct
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

func NewClass(name string, doc string) *Class {
	t := &Class{
		ObjectClass: TypeClass,
		Name:        name,
		Doc:         doc,
		Dict:        map[string]common.Value{},
	}
	TypeDelayReady(t)
	return t
}

func (c *Class) Ready() error {
	return nil
}

func (c *Class) NewClass(name string, doc string, newFunc NewF, construct ConstructF) *Class {
	if newFunc == nil {
		newFunc = c.New
	}

	if construct == nil {
		construct = c.Construct
	}

	return &Class{
		Name:        name,
		Doc:         doc,
		Dict:        map[string]common.Value{},
		IsFinal:     false,
		ObjectClass: c,
		Bases:       []*Class{c},
		Construct:   construct,
	}
}

func (c *Class) GetClass() *Class {
	if c.ObjectClass == nil {
		panic("class is nil")
	}

	return c.ObjectClass
}

func (c *Class) Allocate() *Class {
	return &Class{
		Dict:        map[string]common.Value{},
		ObjectClass: c,
	}
}

// Lookup returns a borrowed reference, and doesn't set an exception,
// returning nil instead.
func (c *Class) Lookup(name string) common.Value {
	for _, baseObj := range c.Bases {
		if res, ok := baseObj.Dict[name]; ok {
			return res
		}
	}

	return nil
}

// NativeGetAttrOrNil gets an attribute from the type of Go type.
func (c *Class) NativeGetAttrOrNil(name string) common.Value {
	// Look in type Dict
	if res, ok := c.Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return c.Lookup(name)
}

func (c *Class) GetAttrOrNil(name string) common.Value {
	// Look in instance dictionary first
	if res, ok := c.Dict[name]; ok {
		return res
	}

	// Then look in type Dict
	if res, ok := c.GetClass().Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return c.Lookup(name)
}

func (c *Class) CallMethod(state common.State, name string, args Tuple, kwargs map[string]common.Value) (
	common.Value,
	bool,
	error,
) {
	fn := c.GetAttrOrNil(name)
	if fn == nil {
		return nil, false, nil
	}

	res, err := Call(state, fn.(*FunctionInstance), args, kwargs)
	return res, true, err
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

func (c *Class) SetAttributes(attrs map[string]common.Value) {
	c.Dict = attrs
	if c.Dict == nil {
		c.Dict = map[string]common.Value{}
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

	return CallAttribute(state, c, operator, common.ConstructorName, *args, *kwargs, true)
}

// getAttribute searches for attribute only in current attributes and
// in Bases.
func (c *Class) getAttribute(name string) (common.Value, error) {
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
	return c.ObjectClass == c
}

func ObjectNew(_ common.State, cls *Class, args Tuple, _ map[string]common.Value) (common.Value, error) {
	// Check arguments to new only for object
	if cls == ObjectClass && excessArgs(args) {
		return nil, ErrorNewf(TypeError, "об_єкт() не приймає аргументів")
	}

	return cls.Allocate(), nil
}

// TypeNew creates a new type.
var TypeNew = &FunctionInstance{
	ClassInstance: ClassInstance{
		class:      Function,
		attributes: nil,
		address:    "",
	},
	package_: nil,
	address:  "",
	Name:     "__новий__",
	Parameters: []FunctionParameter{
		{
			Type:       TypeClass,
			Name:       "я",
			IsVariadic: false,
			IsNullable: false,
		},
		{
			Type:       TypeClass,
			Name:       "я",
			IsVariadic: false,
			IsNullable: false,
		},
	},
	ReturnTypes: nil,
	IsMethod:    true,
	callFunc: func(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
		// Special case: тип(x) should return x.Type
		if cls != nil && len(args) == 1 {
			return args[0].(ObjectInstance).GetClass(), nil
		}

		return nil, ErrorNewf(TypeError, "тип() приймає 1 аргумент")
	},
}

func TypeNeww(_ common.State, cls *Class, args Tuple, _ map[string]common.Value) (common.Value, error) {
	// Special case: тип(x) should return x.Type
	if cls != nil && len(args) == 1 {
		return args[0].(ObjectInstance).GetClass(), nil
	}

	return nil, ErrorNewf(TypeError, "тип() приймає 1 аргумент")

	// if len(args) != 3 {
	// 	return nil, ErrorNewf(TypeError, "тип() приймає 1 або 3 аргументи")
	// }
	//
	// // Check arguments: (name, bases, attributes)
	// var nameObj, basesObj, attributesObj Object
	// err := ParseExactArgs(
	// 	args, "тип:sld",
	// 	&nameObj,
	// 	&basesObj,
	// 	&attributesObj,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	//
	// name := nameObj.(StringClass)
	// bases := basesObj.(Tuple)
	// attributes := attributesObj.(Dict)
	//
	// // Adjust for empty tuple bases
	// if len(bases) == 0 {
	// 	bases = Tuple{Object(ObjectClass)}
	// }
	//
	// for _, newBase := range bases {
	// 	if base, ok := newBase.(*Class); ok {
	// 		if base.Flags&TPFLAGS_BASETYPE == 0 {
	// 			return nil, ErrorNewf(TypeError, "type '%s' is not an acceptable base type", base.Name)
	// 		}
	// 	} else {
	// 		str, err := Str(newBase)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		return nil, ErrorNewf(TypeError, "object '%s' is not a type", str)
	// 	}
	// }
	//
	// dict := attributes.Copy()
	//
	// // Allocate the type object
	// newType := cls.Allocate()
	// newType.New = ObjectNew
	// newType.Construct = ObjectConstruct
	//
	// // Keep name and slots alive in the extended type object
	// et := newType
	// et.Name = string(name)
	//
	// // Initialize Flags
	// newType.Flags = TPFLAGS_BASETYPE
	//
	// // Set Bases
	// newType.Bases = bases
	// bases = nil
	//
	// // Initialize tp_dict from passed-in dict
	// newType.Dict = dict
	//
	// // The __doc__ accessor will first look for Doc;
	// // if that fails, it will still look into __dict__.
	// if doc, ok := dict[builtin.DocAttributeName]; ok {
	// 	if Doc, ok := doc.(StringClass); ok {
	// 		newType.Doc = string(Doc)
	// 	}
	// }
	//
	// // Initialize the rest
	// err = newType.Ready()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return newType, nil
}

func ObjectConstruct(state common.State, self common.Value, args Tuple, kwargs map[string]common.Value) error {
	t := self.(ObjectInstance).GetClass()

	// Check args for object()
	if t == ObjectClass && excessArgs(args) {
		return ErrorNewf(TypeError, "об_єкт.%s() не приймає аргументів", common.ConstructorName)
	}

	// Call the '__конструктор__' method if it exists.
	if _, ok := self.(*Class); ok {
		init := t.GetAttrOrNil(common.ConstructorName)
		if init != nil {
			newArgs := make(Tuple, len(args)+1)
			newArgs[0] = self
			copy(newArgs[1:], args)

			newKwargs := map[string]common.Value{}
			for k, v := range kwargs {
				newKwargs[k] = v
			}

			newKwargs["я"] = self
			if fn, ok := init.(I__call__); ok {
				_, err := fn.__call__(state, newArgs, newKwargs)
				if err != nil {
					return err
				}
			}

			return ErrorNewf(TypeError, "об'єкт '%s' неможливо викликати", init.GetTypeName())
		}
	}

	return nil
}

func TypeConstruct(state common.State, self common.Value, args Tuple, kwargs map[string]common.Value) error {
	if len(args) != 1 && len(args) != 3 {
		return ErrorNewf(TypeError, "тип.%s() приймає 1 або 3 аргументи", common.ConstructorName)
	}

	// Call object.__init__(self) now.
	// XXX Could call super(type, cls).__init__() but what's the point?
	return ObjectConstruct(state, self, nil, nil)
}

func TypeCall(
	state common.State,
	self common.Value,
	name string,
	args Tuple,
	kwargs map[string]common.Value,
) (common.Value, bool, error) {
	t, ok := self.(*Class)
	if !ok {
		return nil, false, nil
	}

	return t.CallMethod(state, name, args, kwargs)
}

// TypeCall0 calls TypeCall with 0 arguments.
func TypeCall0(state common.State, self common.Value, name string) (common.Value, bool, error) {
	return TypeCall(state, self, name, Tuple{self}, map[string]common.Value{"я": self})
}

// TypeCall1 calls TypeCall with 1 argument.
// func TypeCall1(state common.State, self common.Value, name string, arg common.Value) (common.Value, bool, error) {
// 	return TypeCall(self, name, Tuple{self, arg})
// }

// TypeCall2 calls TypeCall with 2 arguments.
// func TypeCall2(state common.State, self common.Value, name string, arg1 common.Value, arg2 common.Value) (
// 	common.Value,
// 	bool,
// 	error,
// ) {
// 	return TypeCall(self, name, Tuple{self, arg1, arg2})
// }

// Return true if any arguments supplied.
func excessArgs(args Tuple) bool {
	return len(args) != 0
}
