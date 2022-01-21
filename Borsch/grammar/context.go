package grammar

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ContextImpl struct {
	scopes       []map[string]common.Type
	package_     *types.PackageInstance
	classContext common.Context
}

func (c *ContextImpl) GetParser() common.Parser {
	return ParserInstance
}

func (c *ContextImpl) PushScope(scope map[string]common.Type) {
	c.scopes = append(c.scopes, scope)
}

func (c *ContextImpl) PopScope() map[string]common.Type {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	lastScopeIdx := len(c.scopes) - 1
	scope := c.scopes[lastScopeIdx]
	c.scopes = c.scopes[:lastScopeIdx]
	return scope
}

func (c *ContextImpl) GetVar(name string) (common.Type, error) {
	lastScopeIdx := len(c.scopes) - 1
	for idx := lastScopeIdx; idx >= 0; idx-- {
		if val, ok := c.scopes[idx][name]; ok {
			return val, nil
		}
	}

	return nil, util.RuntimeError(fmt.Sprintf("ідентифікатор '%s' не визначений", name))
}

func (c *ContextImpl) SetVar(name string, value common.Type) error {
	scopesLen := len(c.scopes)
	for idx := 0; idx < scopesLen; idx++ {
		if oldValue, ok := c.scopes[idx][name]; ok {
			oldValuePrototype := oldValue.(types.ObjectInstance).GetPrototype()
			if oldValuePrototype != value.(types.ObjectInstance).GetPrototype() && oldValuePrototype != types.Nil {
				if idx == scopesLen-1 {
					return util.RuntimeError(
						fmt.Sprintf(
							"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
							value.GetTypeName(), name, oldValue.GetTypeName(),
						),
					)
				}

				// TODO: надрукувати нормальне попередження!
				fmt.Println(
					fmt.Sprintf(
						"Увага: несумісні типи даних '%s' та '%s', змінна '%s' стає недоступною в поточному полі видимості",
						value.GetTypeName(), oldValue.GetTypeName(), name,
					),
				)
				break
			}

			c.scopes[idx][name] = value
			return nil
		}
	}

	c.scopes[scopesLen-1][name] = value
	return nil
}

func (c *ContextImpl) GetClass(name string) (common.Type, error) {
	variable, err := c.classContext.GetVar(name)
	if err == nil {
		if _, ok := variable.(*types.Class); ok {
			return variable, nil
		}
	}

	return nil, util.RuntimeError(fmt.Sprintf("невідомий тип '%s'", name))
}

// GetPackage returns pointer to current evaluating package
// instance without its scope. This method can be called during
// the evaluation process.
func (c *ContextImpl) GetPackage() common.Type {
	return c.package_
}

// BuildPackage sets scope to package instance.
// This method should be called after the package is evaluated.
// After building the package, call GetPackage to retrieve it.
func (c *ContextImpl) BuildPackage() error {
	if len(c.scopes) == 0 {
		return errors.New("not enough scopes")
	}

	c.package_.Attributes = c.scopes[len(c.scopes)-1]
	return nil
}
