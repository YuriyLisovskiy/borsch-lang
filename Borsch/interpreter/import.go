package interpreter

import (
	"path/filepath"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeImport(ctx *Context, node *ast.ImportNode) error {
	if node.IsStd {
		node.FilePath = filepath.Join(i.stdRoot, node.FilePath)
	} else if !filepath.IsAbs(node.FilePath) {
		node.FilePath = filepath.Join(ctx.rootDir, node.FilePath)
	}

	if node.FilePath == ctx.parentPackageName {
		return util.RuntimeError("циклічний імпорт заборонений")
	}

	pkg, ok := i.includedPackages[node.FilePath]
	if !ok {
		var err error
		fileContent, err := util.ReadFile(node.FilePath)
		if err != nil {
			return err
		}

		newContext := &Context{
			rootDir: filepath.Dir(node.FilePath),
			package_: types.NewPackageInstance(
				node.IsStd,
				node.FilePath,
				ctx.package_.Name,
				map[string]types.Type{},
			),
			parentPackageName: ctx.package_.Name,
		}
		newContext.parentObject = newContext.package_
		err = i.ExecuteFile(newContext, fileContent)
		if err != nil {
			return err
		}

		pkg = newContext.package_
		i.includedPackages[node.FilePath] = pkg
	}

	if node.Name != "" {
		err := i.setVar(ctx.GetPackageFromParent(), node.Name, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
