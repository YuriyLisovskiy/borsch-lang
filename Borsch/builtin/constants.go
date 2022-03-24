package builtin

const BORSCH_LIB = "BORSCH_LIB"

const LANGUAGE_FILE_EXT = "борщ"

// Signatures
const (
	LambdaSignature = "<лямбда>"
)

// Built-in types
const (
	AnyTypeName        = "довільний"
	BoolTypeName       = "логічний"
	DictionaryTypeName = "словник"
	FunctionTypeName   = "функція"
	IntegerTypeName    = "цілий"
	ListTypeName       = "список"
	NilTypeName        = "нульовий"
	PackageTypeName    = "пакет"
	RealTypeName       = "дійсний"
	StringTypeName     = "рядок"
	TypeTypeName       = "тип"
	ErrorTypeName      = "Помилка"
)

// Special attributes
const (
	PackageAttributeName  = "__пакет__"
	AttributesName        = "__атрибути__"
	DocAttributeName      = "__документ__"
	ExportedAttributeName = "__експортовані__"
)

// Special operators
const (
	ConstructorName            = "__конструктор__"
	CallOperatorName           = "__оператор_виклику__"
	LengthOperatorName         = "__довжина__"
	BoolOperatorName           = "__логічний__"
	StringOperatorName         = "__рядок__"
	RepresentationOperatorName = "__представлення__"
)
