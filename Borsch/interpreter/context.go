package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type ContextImpl struct {
	scopes        []map[string]common.Value
	parentContext common.Context
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
	lastScopeIdx := len(c.scopes) - 1
	for i := lastScopeIdx; i >= 0; i-- {
		if val, ok := c.scopes[i][name]; ok {
			return val, nil
		}
	}

	if c.parentContext != nil {
		return c.parentContext.GetVar(name)
	}

	return nil, errors.New(fmt.Sprintf("ідентифікатор '%s' не визначений", name))
}

func (c *ContextImpl) SetVar(name string, value common.Value) error {
	if isKeyword(name) {
		return errors.New(
			fmt.Sprintf(
				"неможливо записати значення у '%s', оскільки це ключове слово",
				name,
			),
		)
	}

	size := len(c.scopes)
	for i := 0; i < size; i++ {
		if old, found := c.scopes[i][name]; found {
			oldClass := old.(types.ObjectInstance).GetClass()
			if oldClass != value.(types.ObjectInstance).GetClass() && oldClass != types.Nil {
				if i == size-1 {
					return errors.New(
						fmt.Sprintf(
							"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
							value.GetTypeName(), name, old.GetTypeName(),
						),
					)
				}

				break
			}

			c.scopes[i][name] = value
			return nil
		}
	}

	c.scopes[size-1][name] = value
	return nil
}

func (c *ContextImpl) GetClass(name string) (common.Value, error) {
	if variable, err := c.GetVar(name); err == nil {
		if _, ok := variable.(*types.Class); ok {
			return variable, nil
		}

		return nil, errors.New(fmt.Sprintf("'%s' не є ідентифікатором типу", name))
	}

	return nil, errors.New(fmt.Sprintf("невідомий тип '%s'", name))
}

func (c *ContextImpl) Derive() common.Context {
	return &ContextImpl{
		scopes:        []map[string]common.Value{},
		parentContext: c,
	}
}
