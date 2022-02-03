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
	"github.com/alecthomas/participle/v2/lexer"
)

type Interpreter struct {
	packages    map[string]*types.PackageInstance
	rootContext common.Context
	stacktrace  common.StackTrace
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		packages: map[string]*types.PackageInstance{},
	}

	i.rootContext = &ContextImpl{
		scopes:        []map[string]common.Type{builtin.BuiltinScope},
		classContext:  nil,
		parentContext: nil,
		interpreter:   i,
	}
	return i
}

func (i *Interpreter) Import(state common.State, newPackagePath string) (
	common.Type,
	error,
) {
	parentPackageInstance, _ := state.GetCurrentPackageOrNil().(*types.PackageInstance)
	fullPackagePath, err := getFullPath(newPackagePath, parentPackageInstance)
	if err != nil {
		return nil, err
	}

	if p, ok := i.packages[fullPackagePath]; ok {
		return p, nil
	}

	currPackage := parentPackageInstance
	for currPackage != nil {
		if currPackage.Name == fullPackagePath {
			return nil, util.RuntimeError("циклічний імпорт заборонений")
		}

		currPackage = currPackage.Parent
	}

	packageCode, err := util.ReadFile(fullPackagePath)
	if err != nil {
		return nil, err
	}

	ast, err := state.GetParser().Parse(fullPackagePath, string(packageCode))
	if err != nil {
		return nil, err
	}

	pkg := types.NewPackageInstance(
		i.rootContext.GetChild(),
		false,
		fullPackagePath,
		parentPackageInstance,
		nil,
	)
	ctx := pkg.GetContext()
	if _, err = ast.Evaluate(state.WithContext(ctx).WithPackage(pkg)); err != nil {
		stackTrace := state.GetInterpreter().StackTrace()
		return nil, errors.New(fmt.Sprintf("Відстеження (стек викликів):\n%s", stackTrace.String(err)))
	}

	scope := ctx.TopScope()
	if toExport, err := ctx.GetVar(common.ExportedAttributeName); err == nil {
		switch exported := toExport.(type) {
		case types.ListInstance:
			pkg.Attributes = map[string]common.Type{}
			for _, value := range exported.Values {
				if name, ok := value.(types.StringInstance); ok {
					if attr, ok := scope[name.Value]; ok {
						pkg.Attributes[name.Value] = attr
					}
				}
			}
		case types.StringInstance:
			pkg.Attributes = map[string]common.Type{}
			if attr, ok := scope[exported.Value]; ok {
				pkg.Attributes[exported.Value] = attr
			}
		}
	}

	if pkg.Attributes == nil {
		pkg.Attributes = scope
	}

	i.packages[fullPackagePath] = pkg
	return pkg, nil
}

func (i *Interpreter) Trace(pos lexer.Position, place string, statement string) {
	i.stacktrace.Push(common.NewTraceRow(pos, statement, place))
}

func (i *Interpreter) StackTrace() *common.StackTrace {
	return &i.stacktrace
}

func getFullPath(packagePath string, parentPackage *types.PackageInstance) (string, error) {
	if strings.HasPrefix(packagePath, "!/") {
		packagePath = path.Join(os.Getenv(common.BORSCH_LIB), packagePath[2:])
	} else if !path.IsAbs(packagePath) {
		var err error
		if parentPackage != nil {
			baseDir := path.Dir(parentPackage.Name)
			packagePath = path.Join(baseDir, packagePath)
		} else {
			packagePath, err = filepath.Abs(packagePath)
		}

		if err != nil {
			return "", err
		}
	}

	ext := "." + common.LANGUAGE_FILE_EXT
	if !strings.HasSuffix(packagePath, ext) {
		packagePath += ext
	}

	return packagePath, nil
}
