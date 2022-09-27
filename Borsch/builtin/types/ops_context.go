package types

type Context interface {
	PushScope(scope map[string]Object)
	PopScope() map[string]Object
	TopScope() map[string]Object
	GetVar(name string) (Object, error)
	SetVar(name string, value Object) error
	GetClass(name string) (Object, error)
	Derive() Context
}
