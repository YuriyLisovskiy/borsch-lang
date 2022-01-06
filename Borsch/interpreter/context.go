package interpreter

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

type Context struct {
	parentObject      types.Type
	rootDir           string
	package_          *types.PackageInstance
	parentPackageName string
}

func NewContext(packageName, parentPackage, rootDir string) *Context {
	context := &Context{
		parentObject:      nil,
		rootDir:           rootDir,
		package_:          types.NewPackageInstance(false, packageName, parentPackage, map[string]types.Type{}),
		parentPackageName: parentPackage,
	}
	context.parentObject = context.package_
	return context
}

func (c *Context) WithParent(parent types.Type) *Context {
	return &Context{
		parentObject:      parent,
		rootDir:           c.rootDir,
		package_:          c.package_,
		parentPackageName: c.parentPackageName,
	}
}

func (c *Context) GetPackageFromParent() *types.PackageInstance {
	if p, ok := c.parentObject.(*types.PackageInstance); ok {
		return p
	}

	return nil
}
