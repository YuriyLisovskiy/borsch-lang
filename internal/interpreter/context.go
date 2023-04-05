package interpreter

import (
	"fmt"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

type ContextImpl struct {
	scopes        []map[string]types2.Object
	parentContext types2.Context
}

func (c *ContextImpl) PushScope(scope map[string]types2.Object) {
	c.scopes = append(c.scopes, scope)
}

func (c *ContextImpl) PopScope() map[string]types2.Object {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	lastScopeIdx := len(c.scopes) - 1
	scope := c.scopes[lastScopeIdx]
	c.scopes = c.scopes[:lastScopeIdx]
	return scope
}

func (c *ContextImpl) TopScope() map[string]types2.Object {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	return c.scopes[len(c.scopes)-1]
}

func (c *ContextImpl) GetVar(name string) (types2.Object, error) {
	lastScopeIdx := len(c.scopes) - 1
	for i := lastScopeIdx; i >= 0; i-- {
		if val, ok := c.scopes[i][name]; ok {
			return val, nil
		}
	}

	if c.parentContext != nil {
		return c.parentContext.GetVar(name)
	}

	return nil, types2.NewIdentifierErrorf("ідентифікатор '%s' не визначений", name)
}

func (c *ContextImpl) SetVar(name string, value types2.Object) error {
	if isKeyword(name) {
		return types2.NewIdentifierErrorf(
			"неможливо записати значення у '%s', оскільки це ключове слово",
			name,
		)
	}

	size := len(c.scopes)
	for i := 0; i < size; i++ {
		if old, found := c.scopes[i][name]; found {
			oldClass := old.Class()
			if oldClass != value.Class() && oldClass != types2.NilClass {
				if i == size-1 {
					return types2.NewTypeErrorf(
						"неможливо записати значення типу '%s' у змінну '%s' з типом '%s'",
						value.Class().Name, name, old.Class().Name,
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

func (c *ContextImpl) GetClass(name string) (types2.Object, error) {
	if variable, err := c.GetVar(name); err == nil {
		if _, ok := variable.(*types2.Class); ok {
			return variable, nil
		}

		return nil, types2.NewIdentifierErrorf("'%s' не є ідентифікатором типу", name)
	}

	return nil, types2.NewIdentifierErrorf(fmt.Sprintf("невідомий тип '%s'", name))
}

func (c *ContextImpl) Derive() types2.Context {
	return &ContextImpl{
		scopes:        []map[string]types2.Object{},
		parentContext: c,
	}
}
