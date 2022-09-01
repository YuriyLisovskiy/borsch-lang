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
)

type InterpreterImpl struct {
	packages    map[string]*types.Package
	rootContext types.Context
	parser      Parser
	state       State
}

func NewInterpreter(parser Parser, initialState State) Interpreter {
	i := &InterpreterImpl{
		packages: map[string]*types.Package{},
		parser:   parser,
		state:    initialState,
	}

	GlobalScope["імпорт"] = types.MethodNew(
		"імпорт", BuiltinPackage, []types.MethodParameter{
			{
				Class:      types.StringClass,
				Name:       "шлях",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]types.MethodReturnType{
			{
				Class:      types.PackageClass,
				IsNullable: false,
			},
		},
		func(ctx types.Context, args types.Tuple, kwargs types.StringDict) (types.Object, error) {
			return i.Import(string(args[0].(types.String)))
		},
	)

	i.rootContext = &ContextImpl{
		scopes:        []map[string]types.Object{GlobalScope},
		parentContext: nil,
	}
	return i
}

func (i *InterpreterImpl) Import(newPackagePath string) (
	types.Object,
	error,
) {
	parentPkg, _ := i.state.PackageOrNil().(*types.Package)
	fullPackagePath, err := getFullPath(newPackagePath, parentPkg)
	if err != nil {
		return nil, err
	}

	if p, ok := i.packages[fullPackagePath]; ok {
		return p, nil
	}

	currPackage := parentPkg
	for currPackage != nil {
		if currPackage.Filename == fullPackagePath {
			return nil, errors.New("циклічний імпорт заборонений")
		}

		currPackage = currPackage.Parent
	}

	packageCode, err := readFile(fullPackagePath)
	if err != nil {
		return nil, err
	}

	return i.Evaluate(fullPackagePath, string(packageCode), parentPkg)
}

func (i *InterpreterImpl) Evaluate(packageName, code string, parentPkg *types.Package) (types.Object, error) {
	ast, err := i.parser.Parse(packageName, code)
	if err != nil {
		return nil, err
	}

	pkg := types.PackageNew(packageName, parentPkg, i.rootContext.Derive())
	ctx := pkg.Context
	if _, err = ast.Evaluate(i.state.NewChild().WithContext(ctx).WithPackage(pkg)); err != nil {
		return nil, err
	}

	scope := ctx.TopScope()
	attrs := map[string]types.Object{}
	if toExport, err := ctx.GetVar(builtin.ExportedAttributeName); err == nil {
		switch exported := toExport.(type) {
		// case types.ListInstance:
		// 	for _, value := range exported.Values {
		// 		if name, ok := value.(types.StringInstance); ok {
		// 			if attr, ok := scope[name.Value]; ok {
		// 				attrs[name.Value] = attr
		// 			}
		// 		}
		// 	}
		case types.String:
			if attr, ok := scope[string(exported)]; ok {
				attrs[string(exported)] = attr
			}
		default:
			attrs = scope
		}
	} else {
		attrs = scope
	}

	pkg.Dict = attrs
	i.packages[packageName] = pkg
	return pkg, nil
}

func (i *InterpreterImpl) StackTrace() *common.StackTrace {
	return i.state.StackTrace()
}

func (i *InterpreterImpl) Parser() Parser {
	return i.parser
}

func getFullPath(packagePath string, parentPackage *types.Package) (string, error) {
	if strings.HasPrefix(packagePath, "!/") {
		packagePath = path.Join(os.Getenv(builtin.BORSCH_LIB), packagePath[2:])
	} else if !path.IsAbs(packagePath) {
		var err error
		if parentPackage != nil {
			baseDir := path.Dir(parentPackage.Filename)
			packagePath = path.Join(baseDir, packagePath)
		} else {
			packagePath, err = filepath.Abs(packagePath)
		}

		if err != nil {
			return "", err
		}
	}

	ext := "." + builtin.LANGUAGE_FILE_EXT
	if !strings.HasSuffix(packagePath, ext) {
		packagePath += ext
	}

	return packagePath, nil
}

func readFile(filePath string) (content []byte, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("файл з ім'ям '%s' не існує", filePath))
		return
	}

	content, err = ioutil.ReadFile(filePath)
	return
}
