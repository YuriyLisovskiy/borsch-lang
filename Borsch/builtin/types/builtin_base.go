package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
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
		return util.AttributeIsReadOnlyError(i.GetTypeName(), name)
	}

	return util.AttributeNotFoundError(i.GetTypeName(), name)
}

func (i BuiltinInstance) Call(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error) {
	return nil, util.ObjectIsNotCallable("", i.GetTypeName())
}
