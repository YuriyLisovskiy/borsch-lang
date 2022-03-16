package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type ObjectInstance interface {
	GetClass() *Class
}

var (
	Any *Class = nil
	// TypeClass  *Class = nil
	Nil        *Class = nil
	BoolType   *Class = nil
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

	// def
	// TypeClass = newTypeClass()
	Nil = newNilClass()
	BoolType = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()

	// init
	// InitClass(TypeClass)
	InitClass(Nil)
	InitClass(BoolType)
	InitClass(Dictionary)
	InitClass(Function)
	InitClass(Integer)
	InitClass(List)
	InitClass(Package)
	InitClass(Real)
	InitClass(String)
}

func InitClass(cls *Class) {
	cls.Setup()
	if !cls.IsValid() {
		panic("class is not valid")
	}
}
