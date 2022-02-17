package types

import (
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type ObjectInstance interface {
	GetClass() *Class
	GetAddress() string
}

var (
	AnyClass      *Class = nil
	NilClass      *Class = nil
	BoolClass     *Class = nil
	DictClass     *Class = nil
	FunctionClass *Class = nil
	IntClass      *Class = nil
	ListClass     *Class = nil
	PackageClass  *Class = nil
	RealClass     *Class = nil
	StringClass   *Class = nil
)

var BuiltinPackage *PackageInstance

// delayedReady holds types waiting to be intialised
var delayedReady []*Class

// TypeDelayReady stores the list of types to initialise
//
// Call MakeReady when all initialised
func TypeDelayReady(t *Class) {
	delayedReady = append(delayedReady, t)
}

// TypeMakeReady readies all the types
func TypeMakeReady() (err error) {
	for _, t := range delayedReady {
		err = t.Ready()
		if err != nil {
			return fmt.Errorf("error initialising go type %s: %v", t.Name, err)
		}
	}

	delayedReady = nil
	return nil
}

func init() {
	err := TypeMakeReady()
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	BuiltinPackage = NewPackageInstance(nil, "вбудований", nil, map[string]common.Object{})

	// def
	TypeClass = newTypeClass()
	NilClass = newNilClass()
	BoolClass = newBoolClass()
	DictClass = newDictionaryClass()
	FunctionClass = newFunctionClass()
	IntClass = newIntegerClass()
	ListClass = newListClass()
	PackageClass = NewPackageClass()
	RealClass = newRealClass()
	StringClass = newStringClass()

	// init
	InitClass(TypeClass)
	InitClass(NilClass)
	InitClass(BoolClass)
	InitClass(DictClass)
	InitClass(FunctionClass)
	InitClass(IntClass)
	InitClass(ListClass)
	InitClass(PackageClass)
	InitClass(RealClass)
	InitClass(StringClass)
}

func InitClass(cls *Class) {
	cls.Setup()
	if !cls.IsValid() {
		panic("class is not valid")
	}
}
