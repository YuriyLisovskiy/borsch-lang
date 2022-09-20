package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

var (
	ObjectClass = ClassNew("об_єкт", []*Class{}, map[string]Object{}, false, nil, nil)

	TypeClass = ObjectClass.ClassNew("тип", map[string]Object{}, true, nil, nil)
	AnyClass  = ClassNew("довільний", []*Class{}, map[string]Object{}, true, nil, nil)
)

type (
	NewFunc       func(ctx Context, cls *Class, args Tuple) (Object, error)
	ConstructFunc func(ctx Context, self Object, args Tuple) error
)

type Class struct {
	Name    string
	Dict    map[string]Object
	Bases   []*Class
	IsFinal bool

	New       NewFunc
	Construct ConstructFunc

	// If ClassType is not nil, it is an instance.
	ClassType *Class
}

func init() {
	TypeClass.New = TypeNew
	TypeClass.Construct = TypeConstruct
	ObjectClass.New = ObjectNew
	ObjectClass.Construct = ObjectConstruct
}

func ClassNew(
	name string,
	bases []*Class,
	attributes map[string]Object,
	isFinal bool,
	newF NewFunc,
	constructF ConstructFunc,
) *Class {
	return &Class{
		Name:    name,
		Dict:    attributes,
		Bases:   bases,
		IsFinal: isFinal,

		New:       newF,
		Construct: constructF,

		ClassType: nil,
	}
}

func (value *Class) ClassNew(
	name string,
	attributes map[string]Object,
	isFinal bool,
	newF NewFunc,
	constructF ConstructFunc,
) *Class {
	cls := &Class{
		Name:    name,
		Dict:    attributes,
		Bases:   []*Class{value},
		IsFinal: isFinal,

		ClassType: nil,
	}
	if newF == nil {
		newF = value.New
	}

	if constructF == nil {
		constructF = value.Construct
	}

	cls.New = newF
	cls.Construct = constructF
	return cls
}

func (value *Class) Allocate(attributes map[string]Object) *Class {
	instance := &Class{
		Name:    "",
		Dict:    attributes,
		Bases:   nil,
		IsFinal: false,

		ClassType: value,
	}

	for name, attr := range value.Dict {
		if m, ok := attr.(*Method); ok && m.Name != builtin.LambdaSignature {
			instance.Dict[name] = &MethodWrapper{
				Method:   m,
				Instance: instance,
			}
		}
	}

	return instance
}

func (value *Class) Class() *Class {
	if value.ClassType != nil {
		return value.ClassType
	}

	return TypeClass
}

func (value *Class) represent(Context) (Object, error) {
	if value.ClassType != nil {
		return String(fmt.Sprintf("<об'єкт %s з адресою %p>", value.Class().Name, value)), nil
	}

	return String(fmt.Sprintf("<клас '%s'>", value.Name)), nil
}

func (value *Class) string(ctx Context) (Object, error) {
	return value.represent(ctx)
}

// Lookup returns an attribute from one of the base class,
// and doesn't set an exception, but returns nil instead.
func (value *Class) Lookup(name string) Object {
	for _, base := range value.Bases {
		if res, ok := base.Dict[name]; ok {
			return res
		}
	}

	return nil
}

// NativeGetAttributeOrNil returns an attribute from the class.
func (value *Class) NativeGetAttributeOrNil(name string) Object {
	// Look in type Dict
	if res, ok := value.Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return value.Lookup(name)
}

// GetAttributeOrNil returns attribute from current object,
// it's class or from bases.
func (value *Class) GetAttributeOrNil(name string) Object {
	// Look in instance dictionary first
	if res, ok := value.Dict[name]; ok {
		return res
	}

	// Then look in type Dict
	if res, ok := value.Class().Dict[name]; ok {
		return res
	}

	// Now look through base classes etc
	return value.Lookup(name)
}

func (value *Class) DeleteAttributeOrNil(name string) Object {
	if attr, ok := value.Dict[name]; ok {
		delete(value.Dict, name)
		return attr
	}

	if attr, ok := value.Class().Dict[name]; ok {
		delete(value.Class().Dict, name)
		return attr
	}

	for _, base := range value.Bases {
		if attr, ok := base.Dict[name]; ok {
			delete(base.Dict, name)
			return attr
		}
	}

	return nil
}

func (value *Class) HasBase(cls *Class) bool {
	for _, base := range value.Bases {
		if base == cls {
			return true
		}
	}

	return false
}

func ObjectNew(_ Context, cls *Class, args Tuple) (Object, error) {
	// Check arguments to new only for object
	if cls == ObjectClass && excessArgs(args) {
		return nil, NewErrorf("об_єкт() не приймає аргументів")
	}

	return cls.Allocate(map[string]Object{}), nil
}

func ObjectConstruct(ctx Context, self Object, args Tuple) error {
	t := self.Class()

	// Check args for object()
	if t == ObjectClass && excessArgs(args) {
		return NewErrorf("об_єкт.%s() не приймає аргументів", builtin.ConstructorName)
	}

	// Call the '__конструктор__' method if it exists.
	if _, ok := self.(*Class); ok {
		init := t.GetAttributeOrNil(builtin.ConstructorName)
		if init != nil {
			newArgs := make(Tuple, len(args)+1)
			newArgs[0] = self
			copy(newArgs[1:], args)
			_, err := Call(ctx, init, newArgs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TypeNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	// Special case: тип(x) should return x.Class()
	if cls != nil && len(args) == 1 {
		return args[0].Class(), nil
	}

	// if len(args) != 3 {
	// 	return nil, NewErrorf("тип() приймає 1 або 3 аргументи")
	// }
	return nil, NewErrorf("тип() приймає 1 аргумент")
}

func TypeConstruct(ctx Context, self Object, args Tuple) error {
	if len(args) != 1 && len(args) != 3 {
		return NewErrorf("тип.%s() приймає 1 або 3 аргументи", builtin.ConstructorName)
	}

	// Call об_єкт.__конструктор__(я) now.
	return ObjectConstruct(ctx, self, nil)
}

// Return true if any arguments supplied.
func excessArgs(args Tuple) bool {
	return len(args) != 0
}
