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
	Nil         *Class = nil
	Dictionary  *Class = nil
	Function    *Class = nil
	Integer     *Class = nil
	List        *Class = nil
	Package     *Class = nil
	Real        *Class = nil
	StringClass *Class = nil
)

var BuiltinPackage *PackageInstance

func Init() {
	BuiltinPackage = NewPackageInstance(nil, "вбудований", nil, map[string]common.Value{})

	// def
	// TypeClass = newTypeClass()
	Nil = newNilClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Package = NewPackageClass()
	Real = newRealClass()
	StringClass = newStringClass()
}
