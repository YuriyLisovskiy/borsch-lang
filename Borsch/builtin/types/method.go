package types

var MethodClass = ObjectClass.ClassNew("метод", map[string]Object{}, true, nil, nil)

type MethodFunc func(ctx Context, args Tuple, kwargs StringDict) (Object, error)

type Method struct {
	Name    string
	Package *Package

	Parameters  []MethodParameter
	ReturnTypes []MethodReturnType

	methodF MethodFunc
}

type MethodParameter struct {
	Class      *Class
	Name       string
	IsNullable bool
	IsVariadic bool
}

type MethodReturnType struct {
	Class      *Class
	IsNullable bool
}

func MethodNew(
	name string,
	pkg *Package,
	parameters []MethodParameter,
	returnTypes []MethodReturnType,
	methodF MethodFunc,
) *Method {
	if pkg == nil {
		panic("package is nil")
	}

	if methodF == nil {
		panic("methodF is nil")
	}

	return &Method{
		Name:        name,
		Package:     pkg,
		Parameters:  parameters,
		ReturnTypes: returnTypes,
		methodF:     methodF,
	}
}

func (value *Method) Class() *Class {
	return MethodClass
}

func (value *Method) call(args Tuple) (Object, error) {
	pLen := len(value.Parameters)
	aLen := len(args)
	if pLen != aLen {
		return nil, NewErrorf("кількість параметрів не дорівнює кількості аргументів, %d != %d", pLen, aLen)
	}

	// TODO: take into account variable parameters!

	kwargs := StringDict{}
	for i, arg := range args {
		parameter := value.Parameters[i]
		if err := checkArg(&parameter, arg); err != nil {
			return nil, err
		}

		kwargs[parameter.Name] = arg
	}

	ctx := value.Package.Context.Derive()
	ctx.PushScope(kwargs)
	result, err := value.methodF(ctx, args, kwargs)
	if err != nil {
		return nil, err
	}

	// TODO: check result
	return result, nil
}

func checkArg(parameter *MethodParameter, arg Object) error {
	if parameter.Class == AnyClass {
		return nil
	}

	if arg == Nil {
		if parameter.Class == NilClass || parameter.IsNullable {
			return nil
		}

		// TODO: return error
	}

	if parameter.Class == arg.Class() {
		return nil
	}

	// TODO: return error
	return nil
}

func checkReturnValue(cls *Class, returnValue Object) error {
	// TODO:
	return nil
}
