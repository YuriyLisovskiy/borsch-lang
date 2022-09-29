package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

var (
	ObjectClass = ClassNew("об_єкт", []*Class{}, map[string]Object{}, false, nil, nil)

	TypeClass = ObjectClass.ClassNew("тип", map[string]Object{}, true, nil, nil)
)

type (
	NewFunc       func(ctx Context, cls *Class, args Tuple) (Object, error)
	ConstructFunc func(ctx Context, self Object, args Tuple) error
)

type Class struct {
	Name      string
	Dict      StringDict
	Operators map[common.OperatorHash]*Method
	Bases     []*Class
	IsFinal   bool

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
		Name:      name,
		Dict:      attributes,
		Operators: map[common.OperatorHash]*Method{},
		Bases:     bases,
		IsFinal:   isFinal,

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
		Name:      name,
		Dict:      attributes,
		Operators: map[common.OperatorHash]*Method{},
		Bases:     []*Class{value},
		IsFinal:   isFinal,

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

	// initInstance(cls, &cls.Dict, TypeClass)
	return cls
}

func (value *Class) allocate(attributes map[string]Object) *Class {
	instance := &Class{
		Name:    "",
		Dict:    attributes,
		Bases:   nil,
		IsFinal: false,

		ClassType: value,
	}

	initAttributes(value, instance)

	// for _, base := range value.Bases {
	//
	// }
	//
	// for name, attr := range value.Dict {
	// 	if wrapped, ok := wrapMethod(instance, attr); ok {
	// 		instance.Dict[name] = wrapped
	// 	}
	// }

	return instance
}

func initAttributes(cls, instance *Class) {
	for _, base := range cls.Bases {
		initAttributes(base, instance)
	}

	for name, attr := range cls.Dict {
		if wrapped, ok := wrapMethod(instance, attr); ok {
			instance.Dict[name] = wrapped
		}
	}
}

func (value *Class) Class() *Class {
	if value.IsInstance() {
		return value.ClassType
	}

	return TypeClass
}

// AddAttributes is used to create and set attributes for classes
// after built-in package is initialized.
//
// Do not use this method for class instances!
func (value *Class) AddAttributes(attributes map[string]Object) {
	value.Dict = attributes
}

func (value *Class) call(ctx Context, args Tuple) (Object, error) {
	if value.IsInstance() {
		if attr := value.GetOperatorOrNil(common.CallOp); attr != nil {
			return Call(ctx, attr, append(Tuple{value}, args...))
		}

		return nil, NewTypeErrorf("обʼєкт ʼ%sʼ не може бути викликаний", value.Class().Name)
	}

	// Create instance of the class.
	instance, err := value.New(ctx, value, args)
	if err != nil {
		return nil, err
	}

	if value.Construct != nil {
		err = value.Construct(ctx, instance, args)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (value *Class) getAttribute(_ Context, name string) (Object, error) {
	// TODO: call __отримати_атрибут__ method if exists

	// The (nil, nil) result forces the caller to return the default error.
	return value.GetAttributeOrNil(name), nil
}

func (value *Class) setAttribute(_ Context, name string, newValue Object) error {
	return setAttributeTo(value, &value.Dict, value.GetAttributeOrNil(name), name, newValue)
}

func (value *Class) deleteAttribute(_ Context, name string) (Object, error) {
	if attr, ok := value.Dict[name]; ok {
		delete(value.Dict, name)
		return attr, nil
	}

	// The (nil, nil) result forces the caller to return the default error.
	return nil, nil
}

func (value *Class) represent(ctx Context) (Object, error) {
	if value.IsInstance() {
		if attr := value.GetOperatorOrNil(common.RepresentationOp); attr != nil {
			return Call(ctx, attr, Tuple{value})
		}

		return String(fmt.Sprintf("<об'єкт '%s' з адресою %p>", value.Class().Name, value)), nil
	}

	return String(fmt.Sprintf("<клас '%s'>", value.Name)), nil
}

func (value *Class) string(ctx Context) (Object, error) {
	if value.IsInstance() {
		if attr := value.GetOperatorOrNil(common.StringOp); attr != nil {
			return Call(ctx, attr, Tuple{value})
		}
	}

	return value.represent(ctx)
}

func (value *Class) toInt(ctx Context) (Object, error) {
	if value.IsInstance() {
		if attr := value.GetOperatorOrNil(common.IntOp); attr != nil {
			result, err := Call(ctx, attr, Tuple{value})
			if err != nil {
				return nil, err
			}

			if _, ok := result.(Int); ok {
				return result, nil
			}

			return nil, NewTypeErrorf("%s повернув не цілий тип, а '%s'", common.IntOp.Name(), result.Class().Name)
		}
	}

	// Marks that method could not be executed due to incorrect arguments.
	// Caller should return the default error message in this case.
	return nil, nil
}

func (value *Class) pow(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.PowOp)
}

func (value *Class) mod(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.ModuloOp)
}

func (value *Class) add(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.AddOp)
}

func (value *Class) sub(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.SubOp)
}

func (value *Class) mul(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.MulOp)
}

func (value *Class) div(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.DivOp)
}

func (value *Class) negate(ctx Context) (Object, error) {
	return callUnaryOperator(ctx, value, common.UnaryMinus)
}

func (value *Class) positive(ctx Context) (Object, error) {
	return callUnaryOperator(ctx, value, common.UnaryPlus)
}

func (value *Class) invert(ctx Context) (Object, error) {
	return callUnaryOperator(ctx, value, common.UnaryBitwiseNotOp)
}

func (value *Class) shiftLeft(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.BitwiseLeftShiftOp)
}

func (value *Class) shiftRight(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.BitwiseRightShiftOp)
}

func (value *Class) bitwiseAnd(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.BitwiseAndOp)
}

func (value *Class) bitwiseXor(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.BitwiseXorOp)
}

func (value *Class) bitwiseOr(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.BitwiseOrOp)
}

func (value *Class) equals(ctx Context, other Object) (Object, error) {
	if value.IsInstance() {
		return callBinaryOperator(ctx, value, other, common.EqualsOp)
	}

	if o, ok := other.(*Class); ok && !o.IsInstance() {
		return goBoolToBoolObject(value == o), nil
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для ==: '%s' та '%s'",
		value.Class().Name,
		other.Class().Name,
	)
}

func (value *Class) notEquals(ctx Context, other Object) (Object, error) {
	if value.IsInstance() {
		return callBinaryOperator(ctx, value, other, common.NotEqualsOp)
	}

	if o, ok := other.(*Class); ok && !o.IsInstance() {
		return goBoolToBoolObject(value != o), nil
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для !=: '%s' та '%s'",
		value.Class().Name,
		other.Class().Name,
	)
}

func (value *Class) greater(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.GreaterOp)
}

func (value *Class) greaterOrEquals(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.GreaterOrEqualsOp)
}

func (value *Class) less(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.LessOp)
}

func (value *Class) lessOrEquals(ctx Context, other Object) (Object, error) {
	return callBinaryOperator(ctx, value, other, common.LessOrEqualsOp)
}

// Lookup returns an attribute from one of the base class,
// and doesn't set an exception, but returns nil instead.
func (value *Class) Lookup(name string) Object {
	var bases []*Class
	if value.IsInstance() {
		bases = value.Class().Bases
	} else {
		bases = value.Bases
	}

	for _, base := range bases {
		if res := base.GetAttributeOrNil(name); res != nil {
			return res
		}
	}

	return nil
}

// NativeGetAttributeOrNil returns an attribute from the class.
func (value *Class) NativeGetAttributeOrNil(name string) Object {
	// Look in type Dict
	if res, ok := value.Class().Dict[name]; ok {
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

func (value *Class) GetOperatorOrNil(op common.OperatorHash) Object {
	if value.IsInstance() {
		if res, ok := value.Class().Operators[op]; ok {
			return res
		}
	}

	if res, ok := value.Operators[op]; ok {
		return res
	}

	return nil
}

func (value *Class) IsInstance() bool {
	return value.ClassType != nil
}

func (value *Class) IsBaseOf(cls *Class) bool {
	for _, c := range cls.Bases {
		if value == c || value.IsBaseOf(c) {
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

	return cls.allocate(map[string]Object{}), nil
}

func ObjectConstruct(ctx Context, instance Object, args Tuple) error {
	t := instance.Class()

	// Check args for object()
	if t == ObjectClass && excessArgs(args) {
		return NewErrorf("об_єкт.%s() не приймає аргументів", builtin.ConstructorName)
	}

	// Call the '__конструктор__' method if it exists.
	if constructor := t.GetOperatorOrNil(common.ConstructorOp); constructor != nil {
		_, err := Call(ctx, constructor, append([]Object{instance}, args...))
		if err != nil {
			return err
		}
	}

	return nil
}

func TypeNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	// Special case: тип(x) should return x.Class()
	if cls != nil && len(args) == 1 {
		return args[0].Class(), nil
	}

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

func callBinaryOperator(ctx Context, a, b Object, opHash common.OperatorHash) (Object, error) {
	if op, ok := a.Class().Operators[opHash]; ok {
		return Call(ctx, op, []Object{a, b})
	}

	// Marks that method could not be executed due to incorrect arguments.
	// Caller should return the default error message in this case.
	return nil, nil
}

func callUnaryOperator(ctx Context, a Object, opHash common.OperatorHash) (Object, error) {
	if op, ok := a.Class().Operators[opHash]; ok {
		return Call(ctx, op, []Object{a})
	}

	// Marks that method could not be executed due to incorrect arguments.
	// Caller should return the default error message in this case.
	return nil, nil
}
