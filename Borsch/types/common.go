package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

type ObjectInstance interface {
	GetPrototype() *Class
}

type AttributesInitializer func() map[string]common.Type

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

func Init() {
	BuiltinPackage = NewPackageInstance(true, "вбудований", nil, map[string]common.Type{})

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

	Bool.InitAttributes()
	Dictionary.InitAttributes()
	Function.InitAttributes()
	Integer.InitAttributes()
	List.InitAttributes()
	Nil.InitAttributes()
	Package.InitAttributes()
	Real.InitAttributes()
	String.InitAttributes()
	TypeClass.InitAttributes()
}
