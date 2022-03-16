package types

import (
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

type NewFunc func(cls *Class, args Tuple) (Object, error)

type ConstructFunc func(self Object, args Tuple) error

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
	Dict map[string]Object

	IsFinal bool

	ObjectClass *Class
	Bases       Tuple

	New       NewFunc
	Construct ConstructFunc
	Flags     uint
}

var TypeClass = &Class{
	Name: "тип",
	Doc:  "тип(об_єкт) -> тип об'єкта\nтип(назва, бази, атрибути) -> новий тип",
	Dict: map[string]Object{},
}

var ObjectClass = &Class{
	Name: "об_єкт",
	Doc:  "Базовий тип",
	Dict: map[string]Object{},
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

func NewClass(Name string, Doc string) *Class {
	t := &Class{
		ObjectClass: TypeClass,
		Name:        Name,
		Doc:         Doc,
		Dict:        map[string]Object{},
	}
	TypeDelayReady(t)
	return t
}

func (c *Class) NewClass(name string, doc string, newF NewFunc, constructF ConstructFunc) *Class {
	if newF == nil {
		newF = c.New
	}

	if constructF == nil {
		constructF = c.Construct
	}

	return &Class{
		Name:        name,
		Doc:         doc,
		Dict:        map[string]Object{},
		IsFinal:     false,
		ObjectClass: c,
		Bases:       Tuple{c},
		New:         newF,
		Construct:   constructF,
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
		c.Bases = Tuple{}
	}

	// Initialize tp_dict
	dict := c.Dict
	if dict == nil {
		dict = map[string]Object{}
		c.Dict = dict
	}

	if _, ok := c.Dict[builtin.DocAttributeName]; ok {
		if c.Doc != "" {
			c.Dict[builtin.DocAttributeName] = String(c.Doc)
		} else {
			c.Dict[builtin.DocAttributeName] = Nil
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
	return c.ObjectClass
}

func (c *Class) Allocate() *Class {
	return &Class{
		Dict:        map[string]Object{},
		ObjectClass: c,
	}
}

func (c *Class) __str__() (Object, error) {
	if res, ok, err := c.CallMethod(nil, builtin.StringOperatorName, Tuple{c}); ok {
		return res, err
	}

	return c.__represent__()
}

func (c *Class) __represent__() (Object, error) {
	if res, ok, err := c.CallMethod(nil, builtin.RepresentOperatorName, Tuple{c}); ok {
		return res, err
	}

	if c.Name == "" {
		// FIXME: not a good way to tell objects from classes!
		return String(fmt.Sprintf("<об'єкт %s з адресою %p>", c.Class().Name, c)), nil
	}

	return String(fmt.Sprintf("<клас '%s'>", c.Name)), nil
}

// Lookup returns a borrowed reference, and doesn't set an exception,
// returning nil instead.
func (c *Class) Lookup(name string) Object {
	for _, baseObj := range c.Bases {
		base := baseObj.(*Class)
		if res, ok := base.Dict[name]; ok {
			return res
		}
	}

	return nil
}

// NativeGetAttrOrNil gets an attribute from the type of Go type.
func (c *Class) NativeGetAttrOrNil(name string) Object {
	// Look in type Dict
	if res, ok := c.Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return c.Lookup(name)
}

func (c *Class) GetAttrOrNil(name string) Object {
	// Look in instance dictionary first
	if res, ok := c.Dict[name]; ok {
		return res
	}

	// Then look in type Dict
	if res, ok := c.Class().Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return c.Lookup(name)
}

func (c *Class) CallMethod(state State, name string, args Tuple) (Object, bool, error) {
	fn := c.GetAttrOrNil(name)
	if fn == nil {
		return nil, false, nil
	}

	res, err := Call(state, fn, args)
	return res, true, err
}

func (c *Class) HasBase(cls *Class) bool {
	for _, base := range c.Bases {
		if cls == base {
			return true
		}
	}

	return false
}

func ObjectNew(t *Class, args Tuple) (Object, error) {
	// Check arguments to new only for object
	if t == ObjectClass && excessArgs(args) {
		return nil, ErrorNewf(TypeError, "об_єкт() не приймає аргументів")
	}

	return t.Allocate(), nil
}

// TypeNew creates a new type.
func TypeNew(cls *Class, args Tuple) (Object, error) {
	// Special case: тип(x) should return x.Type
	if cls != nil && len(args) == 1 {
		return args[0].Class(), nil
	}

	if len(args) != 3 {
		return nil, ErrorNewf(TypeError, "тип() приймає 1 або 3 аргументи")
	}

	// Check arguments: (name, bases, attributes)
	var nameObj, basesObj, attributesObj Object
	err := ParseExactArgs(
		args, "тип:sld",
		&nameObj,
		&basesObj,
		&attributesObj,
	)
	if err != nil {
		return nil, err
	}

	name := nameObj.(String)
	bases := basesObj.(Tuple)
	attributes := attributesObj.(Dict)

	// Adjust for empty tuple bases
	if len(bases) == 0 {
		bases = Tuple{Object(ObjectClass)}
	}

	for _, newBase := range bases {
		if base, ok := newBase.(*Class); ok {
			if base.Flags&TPFLAGS_BASETYPE == 0 {
				return nil, ErrorNewf(TypeError, "type '%s' is not an acceptable base type", base.Name)
			}
		} else {
			str, err := Str(newBase)
			if err != nil {
				return nil, err
			}

			return nil, ErrorNewf(TypeError, "object '%s' is not a type", str)
		}
	}

	dict := attributes.Copy()

	// Allocate the type object
	newType := cls.Allocate()
	newType.New = ObjectNew
	newType.Construct = ObjectConstruct

	// Keep name and slots alive in the extended type object
	et := newType
	et.Name = string(name)

	// Initialize Flags
	newType.Flags = TPFLAGS_BASETYPE

	// Set Bases
	newType.Bases = bases
	bases = nil

	// Initialize tp_dict from passed-in dict
	newType.Dict = dict

	// The __doc__ accessor will first look for Doc;
	// if that fails, it will still look into __dict__.
	if doc, ok := dict[builtin.DocAttributeName]; ok {
		if Doc, ok := doc.(String); ok {
			newType.Doc = string(Doc)
		}
	}

	// Initialize the rest
	err = newType.Ready()
	if err != nil {
		return nil, err
	}

	return newType, nil
}

func ObjectConstruct(self Object, args Tuple) error {
	t := self.Class()

	// Check args for object()
	if t == ObjectClass && excessArgs(args) {
		return ErrorNewf(TypeError, "об_єкт.%s() не приймає аргументів", builtin.ConstructorName)
	}

	// Call the '__конструктор__' method if it exists.
	if _, ok := self.(*Class); ok {
		init := t.GetAttrOrNil(builtin.ConstructorName)
		if init != nil {
			newArgs := make(Tuple, len(args)+1)
			newArgs[0] = self
			copy(newArgs[1:], args)
			_, err := Call(nil, init, newArgs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TypeConstruct(self Object, args Tuple) error {
	if len(args) != 1 && len(args) != 3 {
		return ErrorNewf(TypeError, "тип.%s() приймає 1 або 3 аргументи", builtin.ConstructorName)
	}

	// Call object.__init__(self) now.
	// XXX Could call super(type, cls).__init__() but what's the point?
	return ObjectConstruct(self, nil)
}

func TypeCall(self Object, name string, args Tuple) (Object, bool, error) {
	t, ok := self.(*Class)
	if !ok {
		return nil, false, nil
	}

	return t.CallMethod(nil, name, args)
}

// TypeCall0 calls TypeCall with 0 arguments.
func TypeCall0(self Object, name string) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self})
}

// TypeCall1 calls TypeCall with 1 argument.
func TypeCall1(self Object, name string, arg Object) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self, arg})
}

// TypeCall2 calls TypeCall with 2 arguments.
func TypeCall2(self Object, name string, arg1 Object, arg2 Object) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self, arg1, arg2})
}

// Return true if any arguments supplied.
func excessArgs(args Tuple) bool {
	return len(args) != 0
}
