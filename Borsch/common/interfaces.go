package common

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
	NewContext(packageName string, parentPackage Type) Context
}

type Context interface {
	GetParser() Parser
	PushScope(scope map[string]Type)
	PopScope() map[string]Type
	GetVar(name string) (Type, error)
	SetVar(name string, value Type) error
	GetClass(name string) (Type, error)
	BuildPackage() error
	GetPackage() Type
	GetChild() Context
}

type Evaluatable interface {
	Evaluate(Context) (Type, error)
}

type OperatorEvaluatable interface {
	Evaluate(Context, Type) (Type, error)
}

type Type interface {
	String(Context) string
	Representation(Context) string
	AsBool(Context) bool
	GetTypeName() string
	GetAttribute(string) (Type, error)
	SetAttribute(string, Type) (Type, error)
}

type SequentialType interface {
	Length(Context) int64
	GetElement(Context, int64) (Type, error)
	SetElement(Context, int64, Type) (Type, error)
	Slice(Context, int64, int64) (Type, error)
}

type CallableType interface {
	Call(Context, *[]Type, *map[string]Type) (Type, error)
}
