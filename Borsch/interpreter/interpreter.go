package interpreter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
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
		scopes:        []map[string]common.Value{builtin.GlobalScope},
		classContext:  nil,
		parentContext: nil,
		interpreter:   i,
	}
	return i
}

func (i *Interpreter) Import(state common.State, newPackagePath string) (
	common.Value,
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
			return nil, utilities.RuntimeError("циклічний імпорт заборонений")
		}

		currPackage = currPackage.Parent
	}

	packageCode, err := readFile(fullPackagePath)
	if err != nil {
		return nil, err
	}

	ast, err := state.GetParser().Parse(fullPackagePath, string(packageCode))
	if err != nil {
		return nil, err
	}

	pkg := types.NewPackageInstance(
		i.rootContext.GetChild(),
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
	attrs := map[string]common.Value{}
	if toExport, err := ctx.GetVar(common.ExportedAttributeName); err == nil {
		switch exported := toExport.(type) {
		case types.ListInstance:
			for _, value := range exported.Values {
				if name, ok := value.(types.StringInstance); ok {
					if attr, ok := scope[name.Value]; ok {
						attrs[name.Value] = attr
					}
				}
			}
		case types.StringInstance:
			if attr, ok := scope[exported.Value]; ok {
				attrs[exported.Value] = attr
			}
		default:
			attrs = scope
		}
	} else {
		attrs = scope
	}

	pkg.SetAttributes(attrs)
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

func readFile(filePath string) (content []byte, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = utilities.RuntimeError(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
		return
	}

	content, err = ioutil.ReadFile(filePath)
	return
}
