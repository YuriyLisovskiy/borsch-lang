package interpreter

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

type Context struct {
	parentObject      types.Type
	rootDir           string
	package_          *types.PackageInstance
	parentPackageName string
}

func NewContext(packageName, parentPackage, rootDir string) *Context {
	return &Context{
		parentObject:      nil,
		rootDir:           rootDir,
		package_:          types.NewPackageInstance(false, packageName, parentPackage, map[string]types.Type{}),
		parentPackageName: parentPackage,
	}
}
