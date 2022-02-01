package interpreter

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Interpreter struct {
	parser   common.Parser
	packages map[string]*types.PackageInstance
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		parser:   ParserInstance,
		packages: map[string]*types.PackageInstance{},
	}
}

func (i *Interpreter) Import(newPackagePath string, parentPackage common.Type) (common.Type, error) {
	var parentPackageInstance *types.PackageInstance
	if parentPackage != nil {
		var ok bool
		if parentPackageInstance, ok = parentPackage.(*types.PackageInstance); !ok {
			return nil, errors.New("non-package instance received")
		}
	}

	if strings.HasPrefix(newPackagePath, "!/") {
		newPackagePath = path.Join(os.Getenv(common.BORSCH_LIB), newPackagePath[2:])
	} else if !path.IsAbs(newPackagePath) {
		var err error
		if parentPackageInstance != nil {
			baseDir := path.Dir(parentPackageInstance.Name)
			newPackagePath = path.Join(baseDir, newPackagePath)
		} else {
			newPackagePath, err = filepath.Abs(newPackagePath)
		}

		if err != nil {
			return nil, err
		}
	}

	if p, ok := i.packages[newPackagePath]; ok {
		return p, nil
	}

	currPackage := parentPackageInstance
	for currPackage != nil {
		if currPackage.Name == newPackagePath {
			return nil, util.RuntimeError("циклічний імпорт заборонений")
		}

		currPackage = currPackage.Parent
	}

	packageCode, err := util.ReadFile(newPackagePath)
	if err != nil {
		return nil, err
	}

	ast, err := i.parser.Parse(newPackagePath, string(packageCode))
	if err != nil {
		return nil, err
	}

	context := &ContextImpl{
		scopes:      []map[string]common.Type{builtin.BuiltinScope},
		interpreter: i,
	}
	pkg := types.NewPackageInstance(
		context,
		false,
		newPackagePath,
		parentPackageInstance,
		nil,
	)
	context.package_ = pkg
	if _, err = ast.Evaluate(context); err != nil {
		return nil, errors.New(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
	}

	pkg.Attributes = context.scopes[len(context.scopes)-1]
	i.packages[newPackagePath] = pkg
	return context.package_, nil
}
