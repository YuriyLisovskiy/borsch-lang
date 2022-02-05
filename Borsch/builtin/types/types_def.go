package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type ObjectInstance interface {
	GetClass() *Class
	GetAddress() string
}

var (
	Any        *Class = nil
	TypeClass  *Class = nil
	Nil        *Class = nil
	Bool       *Class = nil
	Dictionary *Class = nil
	Function   *Class = nil
	Integer    *Class = nil
	List       *Class = nil
	Package    *Class = nil
	Real       *Class = nil
	String     *Class = nil
)

var BuiltinPackage *PackageInstance

func Init() {
	BuiltinPackage = NewPackageInstance(nil, "вбудований", nil, map[string]common.Value{})

	TypeClass = newTypeClass()
	Nil = newNilClass()
	Bool = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()

	TypeClass.InitAttributes()
	Nil.InitAttributes()
	Bool.InitAttributes()
	Dictionary.InitAttributes()
	Function.InitAttributes()
	Integer.InitAttributes()
	List.InitAttributes()
	Package.InitAttributes()
	Real.InitAttributes()
	String.InitAttributes()
}
