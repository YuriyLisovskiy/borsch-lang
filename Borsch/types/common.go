package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

type ObjectInstance interface {
	GetPrototype() *Class
}

var (
	Bool       *Class = nil
	Dictionary *Class = nil
	Function   *Class = nil
	Integer    *Class = nil
	List       *Class = nil
	Nil        *Class = nil
	Package    *Class = nil
	Real       *Class = nil
	String     *Class = nil
	TypeClass  *Class = nil
	Any        *Class = nil
)

var BuiltinPackage *PackageInstance

func init() {
	Bool = newBoolClass()
	Dictionary = newDictionaryClass()
	Function = newFunctionClass()
	Integer = newIntegerClass()
	List = newListClass()
	Nil = newNilClass()
	Package = NewPackageClass()
	Real = newRealClass()
	String = newStringClass()
	TypeClass = newTypeClass()

	BuiltinPackage = NewPackageInstance(true, "вбудований", nil, map[string]common.Type{})
}
