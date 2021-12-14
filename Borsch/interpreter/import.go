package interpreter

import (
	"path/filepath"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeImport(
	node *ast.ImportNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	if node.IsStd {
		node.FilePath = filepath.Join(i.stdRoot, node.FilePath)
	} else if !filepath.IsAbs(node.FilePath) {
		node.FilePath = filepath.Join(rootDir, node.FilePath)
	}

	if node.FilePath == parentPackage {
		return nil, util.RuntimeError("циклічний імпорт заборонений")
	}

	pkg, ok := i.includedPackages[node.FilePath]
	if !ok {
		var err error
		fileContent, err := util.ReadFile(node.FilePath)
		if err != nil {
			return nil, err
		}

		pkg, err = i.ExecuteFile(node.FilePath, thisPackage, fileContent, node.IsStd)
		if err != nil {
			return nil, err
		}

		i.includedPackages[node.FilePath] = pkg
	}

	if node.Name != "" {
		err := i.setVar(thisPackage, node.Name, pkg)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}
