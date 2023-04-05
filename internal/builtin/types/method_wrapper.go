package types

var MethodWrapperClass = ObjectClass.ClassNew("метод-обгортка", map[string]Object{}, true, nil, nil)

type MethodWrapper struct {
	Method   *Method
	Instance Object
}

func (value *MethodWrapper) Class() *Class {
	return MethodWrapperClass
}

func (value *MethodWrapper) call(ctx Context, args Tuple) (Object, error) {
	if value.Instance == nil {
		return nil, NewValueError("екземпляр класу не існує")
	}

	if value.Method == nil {
		return nil, NewValueErrorf("оригінальний метод класу %s не існує", value.Instance.Class().Name)
	}

	return value.Method.call(ctx, append([]Object{value.Instance}, args...))
}
