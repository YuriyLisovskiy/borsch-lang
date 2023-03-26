package types

func init() {
	ErrorClass = ObjectClass.ClassNew("Помилка", map[string]Object{}, false, nil, ErrorConstruct)

	AssertionErrorClass = ErrorClass.ClassNew("ПомилкаПрипущення", map[string]Object{}, false, AssertionErrorNew, nil)

	AttributeErrorClass = ErrorClass.ClassNew("ПомилкаАтрибута", map[string]Object{}, false, AttributeErrorNew, nil)

	IdentifierErrorClass = ErrorClass.ClassNew(
		"ПомилкаІдентифікатора",
		map[string]Object{},
		false,
		IdentifierErrorNew,
		nil,
	)

	IndexOutOfRangeErrorClass = ErrorClass.ClassNew(
		"ПомилкаІндексу",
		map[string]Object{},
		false,
		IndexOutOfRangeErrorNew,
		nil,
	)

	RuntimeErrorClass = ErrorClass.ClassNew("ПомилкаВиконання", map[string]Object{}, false, RuntimeErrorNew, nil)

	TypeErrorClass = ErrorClass.ClassNew("ПомилкаТипу", map[string]Object{}, false, TypeErrorNew, nil)

	ValueErrorClass = ErrorClass.ClassNew("ПомилкаЗначення", map[string]Object{}, false, ValueErrorNew, nil)

	ZeroDivisionErrorClass = ErrorClass.ClassNew(
		"ПомилкаДіленняНаНуль",
		map[string]Object{},
		false,
		ZeroDivisionErrorNew,
		nil,
	)

	overflowErrorGo = NewErrorf("ціле число занадто велике, щоб перетворити його в Go int")
}
