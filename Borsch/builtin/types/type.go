package types

import (
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type NewFunc func(cls *Type, args Tuple, kwargs StringDict) (Object, error)

type ConstructFunc func(self Object, args Tuple, kwargs StringDict) error

const (
	// TPFLAGS_BASETYPE is set if the type allows subclassing.
	TPFLAGS_BASETYPE uint = 1 << 10

	// TPFLAGS_READY is set if the type is 'ready' -- fully initialized.
	TPFLAGS_READY uint = 1 << 12

	// TPFLAGS_READYING is set while the type is being 'readied', to prevent
	// recursive ready calls.
	TPFLAGS_READYING uint = 1 << 13
)

type Type struct {
	Name string
	Doc  string
	Dict StringDict

	IsFinal bool

	ObjectType *Type
	Bases      Tuple

	Parent common.Value

	Initializer func(*map[string]common.Value)
	New         NewFunc
	Construct   ConstructFunc

	Flags uint
}

var TypeType = &Type{
	Name: "тип",
	Doc:  "тип(об_єкт) -> тип об'єкта\nтип(назва, бази, словник) -> новий тип",
	Dict: StringDict{},
}

var ObjectType = &Type{
	Name: "об_єкт",
	Doc:  "Базовий тип",
	Dict: StringDict{},
}

func init() {
	// Initialised like this to avoid initialisation loops
	TypeType.New = TypeNew
	TypeType.Construct = TypeInit
	TypeType.ObjectType = TypeType
	ObjectType.New = ObjectNew
	ObjectType.Construct = ObjectConstruct
	ObjectType.ObjectType = TypeType
	err := TypeType.Ready()
	if err != nil {
		log.Fatal(err)
	}

	err = ObjectType.Ready()
	if err != nil {
		log.Fatal(err)
	}
}

func NewType(Name string, Doc string) *Type {
	t := &Type{
		ObjectType: TypeType,
		Name:       Name,
		Doc:        Doc,
		Dict:       StringDict{},
	}
	TypeDelayReady(t)
	return t
}

func (t *Type) NewType(name string, doc string, newF NewFunc, constructF ConstructFunc) *Type {
	if newF == nil {
		newF = t.New
	}

	if constructF == nil {
		constructF = t.Construct
	}

	return &Type{
		Name:        name,
		Doc:         doc,
		Dict:        StringDict{},
		IsFinal:     false,
		ObjectType:  t,
		Bases:       Tuple{t},
		Parent:      nil,
		Initializer: nil,
		New:         newF,
		Construct:   constructF,
	}
}

// delayedReady holds types waiting to be intialised
var delayedReady []*Type

// TypeDelayReady stores the list of types to initialise
//
// Call MakeReady when all initialised
func TypeDelayReady(t *Type) {
	delayedReady = append(delayedReady, t)
}

// TypeMakeReady readies all the types
func TypeMakeReady() (err error) {
	for _, t := range delayedReady {
		err = t.Ready()
		if err != nil {
			return fmt.Errorf("Error initialising go type %s: %v", t.Name, err)
		}
	}
	delayedReady = nil
	return nil
}

func init() {
	err := TypeMakeReady()
	if err != nil {
		log.Fatal(err)
	}
}

func (t *Type) Type() *Type {
	return t.ObjectType
}

func (t *Type) Error() string {
	return t.Name
}

func (t *Type) GetDict() StringDict {
	return t.Dict
}

func (t *Type) Allocate() *Type {
	return &Type{
		Dict:       StringDict{},
		ObjectType: t,
	}
}

func (t *Type) IsValid() bool {
	if len(t.Name) == 0 {
		return false
	}

	if t.Dict == nil {
		return false
	}

	if t.ObjectType == nil {
		return false
	}

	if t.Parent == nil {
		return false
	}

	if t.Construct == nil {
		return false
	}

	return true
}

func (t *Type) __call__(args Tuple, kwargs StringDict) (Object, error) {
	if t.New == nil {
		return nil, ErrorNewf(TypeError, "cannot create '%s' instances", t.Name)
	}

	obj, err := t.New(t, args, kwargs)
	if err != nil {
		return nil, err
	}

	// When the call was тип(щось), don't call __конструктор__ on the result.
	if t == TypeType && len(args) == 1 && len(kwargs) == 0 {
		return obj, nil
	}

	// If the returned object is not an instance of type,
	// it won't be initialized.
	if !obj.Type().HasBase(t) {
		return obj, nil
	}

	objType := obj.Type()
	if objType.Construct != nil {
		err = objType.Construct(obj, args, kwargs)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Lookup returns a borrowed reference, and doesn't set an exception,
// returning nil instead.
func (t *Type) Lookup(name string) Object {
	for _, baseObj := range t.Bases {
		base := baseObj.(*Type)
		if res, ok := base.Dict[name]; ok {
			return res
		}
	}

	return nil
}

func (t *Type) GetAttrOrNil(name string) Object {
	// Look in instance dictionary first
	if res, ok := t.Dict[name]; ok {
		return res
	}

	// Then look in type Dict
	if res, ok := t.Type().Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return t.Lookup(name)
}

func (t *Type) HasBase(cls *Type) bool {
	for _, base := range t.Bases {
		if cls == base {
			return true
		}
	}

	return false
}

func (t *Type) CallMethod(name string, args Tuple, kwargs StringDict) (Object, bool, error) {
	fn := t.GetAttrOrNil(name)
	if fn == nil {
		return nil, false, nil
	}

	res, err := Call(fn, args, kwargs)
	return res, true, err
}

func ObjectConstruct(self Object, args Tuple, kwargs StringDict) error {
	t := self.Type()

	// Check args for object()
	if t == ObjectType && excessArgs(args, kwargs) {
		return ErrorNewf(TypeError, "об_єкт.%s() не приймає аргументів", common.ConstructorName)
	}

	// Call the '__конструктор__' method if it exists.
	if _, ok := self.(*Type); ok {
		init := t.GetAttrOrNil(common.ConstructorName)
		if init != nil {
			newArgs := make(Tuple, len(args)+1)
			newArgs[0] = self
			copy(newArgs[1:], args)
			_, err := Call(init, newArgs, kwargs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ObjectNew(t *Type, args Tuple, kwargs StringDict) (Object, error) {
	// Check arguments to new only for object
	if t == ObjectType && excessArgs(args, kwargs) {
		return nil, ErrorNewf(TypeError, "об_єкт() не приймає аргументів")
	}

	return t.Allocate(), nil
}

// TypeNew creates a new type.
func TypeNew(cls *Type, args Tuple, kwargs StringDict) (Object, error) {
	// Special case: тип(x) should return x.Type
	if cls != nil && len(args) == 1 && len(kwargs) == 0 {
		return args[0].Type(), nil
	}

	if len(args)+len(kwargs) != 3 {
		return nil, ErrorNewf(TypeError, "тип() приймає 1 або 3 аргументи")
	}

	// Check arguments: (name, bases, attributes)
	var nameObj, basesObj, origDictObj Object
	err := ParseTupleAndKeywords(
		args, kwargs, "UOO:type", []string{"назва", "бази", "атрибути"},
		&nameObj,
		&basesObj,
		&origDictObj,
	)
	if err != nil {
		return nil, err
	}

	name := nameObj.(String)
	bases := basesObj.(Tuple)
	origDict := origDictObj.(StringDict)

	// Adjust for empty tuple bases
	if len(bases) == 0 {
		bases = Tuple{Object(ObjectType)}
	}

	for _, newBase := range bases {
		if base, ok := newBase.(*Type); ok {
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

	dict := origDict.Copy()

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
	if doc, ok := dict[common.DocAttributeName]; ok {
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

func TypeInit(self Object, args Tuple, kwargs StringDict) error {
	if len(kwargs) != 0 {
		return ErrorNewf(TypeError, "тип.%s() не приймає аргументи у вигляді словника", common.ConstructorName)
	}

	if len(args) != 1 && len(args) != 3 {
		return ErrorNewf(TypeError, "тип.%s() приймає 1 або 3 аргументи", common.ConstructorName)
	}

	// Call object.__init__(self) now.
	// XXX Could call super(type, cls).__init__() but what's the point?
	return ObjectConstruct(self, nil, nil)
}

func TypeCall(self Object, name string, args Tuple, kwargs StringDict) (Object, bool, error) {
	t, ok := self.(*Type)
	if !ok {
		return nil, false, nil
	}

	return t.CallMethod(name, args, kwargs)
}

// TypeCall0 calls TypeCall with 0 arguments.
func TypeCall0(self Object, name string) (Object, bool, error) {
	return TypeCall(self, name, Tuple{self}, nil)
}

// Ready the type for use.
//
// Returns an error on problems.
func (t *Type) Ready() error {
	if t.Flags&TPFLAGS_READY != 0 {
		if t.Dict == nil {
			return ErrorNewf(SystemError, "Type.Ready is Ready but Dict is nil")
		}

		return nil
	}

	if t.Flags&TPFLAGS_READYING != 0 {
		return ErrorNewf(SystemError, "Type.Ready already readying")
	}

	t.Flags |= TPFLAGS_READYING

	// Now the only way base can still be nil is if type is
	// ObjectType.

	// Initialize Bases
	if t.Bases == nil {
		t.Bases = Tuple{}
	}

	// Initialize tp_dict
	dict := t.Dict
	if dict == nil {
		dict = NewStringDict()
		t.Dict = dict
	}

	if _, ok := t.Dict[common.DocAttributeName]; ok {
		if t.Doc != "" {
			t.Dict[common.DocAttributeName] = String(t.Doc)
		} else {
			t.Dict[common.DocAttributeName] = Nil
		}
	}

	// Link into each base class's list of subclasses
	bases := t.Bases
	for i := range bases {
		b, ok := bases[i].(*Type)
		if ok {
			addSubclass(b, t)
		}
	}

	// All done -- set the ready flag
	if t.Dict == nil {
		panic("Type.Ready Dict is nil")
	}

	t.Flags = (t.Flags &^ TPFLAGS_READYING) | TPFLAGS_READY
	return nil
}

// Return true if any arguments supplied.
func excessArgs(args Tuple, kwargs StringDict) bool {
	return len(args) != 0 || len(kwargs) != 0
}

func addSubclass(base, tp *Type) {
	// TODO: addSubclass
}

func (t *Type) __equal__(other Object) (Object, error) {
	if otherT, ok := other.(*Type); ok && t == otherT {
		return True, nil
	}

	return False, nil
}

func (t *Type) __not_equal__(other Object) (Object, error) {
	if otherT, ok := other.(*Type); ok && t == otherT {
		return False, nil
	}

	return True, nil
}

func (t *Type) __str__() (Object, error) {
	if res, ok, err := t.CallMethod(common.StringOperator, Tuple{t}, nil); ok {
		return res, err
	}

	return t.__represent__()
}

func (t *Type) __represent__() (Object, error) {
	if res, ok, err := t.CallMethod(common.RepresentOperator, Tuple{t}, nil); ok {
		return res, err
	}

	if t.Name == "" {
		// FIXME: not a good way to tell objects from classes!
		return String(fmt.Sprintf("<об'єкт %s з адресою %p>", t.Type().Name, t)), nil
	}

	return String(fmt.Sprintf("<клас '%s'>", t.Name)), nil
}

// Make sure it satisfies the interface
var _ Object = (*Type)(nil)
var _ I__call__ = (*Type)(nil)
var _ IGetDict = (*Type)(nil)
var _ I__represent__ = (*Type)(nil)
var _ I__str__ = (*Type)(nil)
