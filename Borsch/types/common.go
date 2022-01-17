package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

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
	TypeClassTypeHash
)

type Instance interface {
	GetClass() *Class
}

var (
	Bool       = newBoolClass()
	Dictionary = newDictionaryClass()
	Function   = newFunctionClass()
	Integer    = newIntegerClass()
	List       = newListClass()
	Nil        = newNilClass()
	Package    = NewPackageClass()
	Real       = newRealClass()
	String     = newStringClass()
	TypeClass  = newTypeClass()
)

var BuiltinPackage = NewPackageInstance(true, "вбудований", "", map[string]common.Type{})
