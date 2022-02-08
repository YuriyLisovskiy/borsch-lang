package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

// BuiltinInstance is non-callable instance where setting an
// attribute is restricted operation.
//
// Note: these rules are not applicable to FunctionInstance and PackageInstance.
type BuiltinInstance struct {
	ClassInstance
}

func (i BuiltinInstance) SetAttribute(name string, _ common.Value) error {
	if i.HasAttribute(name) {
		return utilities.AttributeIsReadOnlyError(i.GetTypeName(), name)
	}

	return utilities.AttributeNotFoundError(i.GetTypeName(), name)
}

func (i BuiltinInstance) Call(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error) {
	return nil, utilities.ObjectIsNotCallable("", i.GetTypeName())
}
