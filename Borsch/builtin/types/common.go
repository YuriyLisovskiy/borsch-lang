package types

const (
	AnyTypeHash uint64 = iota + 1
	NilTypeHash
	RealTypeHash
	IntegerTypeHash
	StringTypeHash
	BoolTypeHash
	ListTypeHash
	DictionaryTypeHash
	PackageTypeHash
	FunctionTypeHash
)

type Type interface {
	String() string
	Representation() string
	GetTypeHash() uint64
	GetTypeName() string
	AsBool() bool
	GetAttribute(string) (Type, error)
	SetAttribute(string, Type) (Type, error)
}

type SequentialType interface {
	Length() int64
	GetElement(int64) (Type, error)
	SetElement(int64, Type) (Type, error)
	Slice(int64, int64) (Type, error)
}

type CallableType interface {
	Call(*[]Type, *map[string]Type) (Type, error)
}

type Instance interface {
	GetClass() *Class
}

var (
	Bool = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Nil = newNilClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()
)

var BuiltinPackage *PackageInstance = nil

func init() {
	BuiltinPackage = NewPackageInstance(true, "", "", map[string]Type{})
}
