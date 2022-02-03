package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ContextImpl struct {
	scopes []map[string]common.Value
	// package_      *types.PackageInstance
	classContext  common.Context
	parentContext common.Context
	interpreter   common.Interpreter
}

func (c *ContextImpl) PushScope(scope map[string]common.Value) {
	c.scopes = append(c.scopes, scope)
}

func (c *ContextImpl) PopScope() map[string]common.Value {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	lastScopeIdx := len(c.scopes) - 1
	scope := c.scopes[lastScopeIdx]
	c.scopes = c.scopes[:lastScopeIdx]
	return scope
}

func (c *ContextImpl) TopScope() map[string]common.Value {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	return c.scopes[len(c.scopes)-1]
}

func (c *ContextImpl) GetVar(name string) (common.Value, error) {
	switch name {
	case "нуль":
		return types.NewNilInstance(), nil
	case "нульовий":
		return types.Nil, nil
	}

	lastScopeIdx := len(c.scopes) - 1
	for idx := lastScopeIdx; idx >= 0; idx-- {
		if val, ok := c.scopes[idx][name]; ok {
			return val, nil
		}
	}

	if c.parentContext != nil {
		return c.parentContext.GetVar(name)
	}

	return nil, util.RuntimeError(fmt.Sprintf("ідентифікатор '%s' не визначений", name))
}

func (c *ContextImpl) SetVar(name string, value common.Value) error {
	switch name {
	case "нуль", "нульовий":
		return util.RuntimeError(fmt.Sprintf("неможливо записати значення у '%s'", name))
	}

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

func (c *ContextImpl) GetClass(name string) (common.Value, error) {
	var variable common.Value
	var err error
	if c.classContext != nil {
		variable, err = c.classContext.GetVar(name)
	} else {
		variable, err = c.GetVar(name)
	}

	if err == nil {
		if _, ok := variable.(*types.Class); ok {
			return variable, nil
		}
	}

	if c.parentContext != nil {
		return c.parentContext.GetClass(name)
	}

	return nil, util.RuntimeError(fmt.Sprintf("невідомий тип '%s'", name))
}

func (c *ContextImpl) GetChild() common.Context {
	return c.getChildContext()
}

func (c *ContextImpl) getChildContext() *ContextImpl {
	return &ContextImpl{
		scopes: nil,
		// package_:      c.package_,
		classContext:  nil,
		parentContext: c,
		interpreter:   c.interpreter,
	}
}
