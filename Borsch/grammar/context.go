package grammar

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Context struct {
	scopes   []Scope
	package_ *types.PackageInstance
}

func NewContext(packageFilename string, parentPackage *types.PackageInstance) *Context {
	parentPackageName := ""
	if parentPackage != nil {
		parentPackageName = parentPackage.Name
	}
	return &Context{
		package_: types.NewPackageInstance(false, packageFilename, parentPackageName, map[string]types.Type{}),
	}
}

func (c *Context) PushScope(scope Scope) {
	c.scopes = append(c.scopes, scope)
}

func (c *Context) PopScope() map[string]types.Type {
	if len(c.scopes) == 0 {
		panic("fatal: not enough scopes")
	}

	lastScopeIdx := len(c.scopes) - 1
	scope := c.scopes[lastScopeIdx]
	c.scopes = c.scopes[:lastScopeIdx]
	return scope
}

func (c *Context) getVar(name string) (types.Type, error) {
	lastScopeIdx := len(c.scopes) - 1
	for idx := lastScopeIdx; idx >= 0; idx-- {
		if val, ok := c.scopes[idx][name]; ok {
			return val, nil
		}
	}

	return nil, util.RuntimeError(fmt.Sprintf("ідентифікатор '%s' не визначений", name))
}

func (c *Context) setVar(name string, value types.Type) error {
	scopesLen := len(c.scopes)
	for idx := 0; idx < scopesLen; idx++ {
		if oldValue, ok := c.scopes[idx][name]; ok {
			if oldValue.GetTypeHash() != value.GetTypeHash() {
				if scopesLen == 1 {
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
