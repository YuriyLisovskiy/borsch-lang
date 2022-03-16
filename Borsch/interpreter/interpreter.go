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
)

type Interpreter struct {
	packageStore *types.PackageStore
	rootContext  types.Context
	stacktrace   types.StackTrace
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		packageStore: types.NewModuleStore(),
	}

	i.rootContext = &ContextImpl{
		scopes: []map[string]types.Object{
			{
				"друкр": types.MustNewMethod(
					"друкр",
					func(state types.State, self types.Object, args types.Tuple) (types.Object, error) {
						stringResult := ""
						for _, arg := range args {
							argStr, err := types.StrAsString(arg)
							if err != nil {
								return nil, err
							}

							stringResult += argStr
						}

						fmt.Print(stringResult)
						return types.Nil, nil
					},
					0,
					"", // TODO: add doc
				),
				"імпорт": types.MustNewMethod(
					"імпорт",
					func(state types.State, self, packagePath types.Object) (types.Object, error) {
						strPath, err := types.StrAsString(packagePath)
						if err != nil {
							return nil, err
						}
						
						return i.Import(state, strPath)
					},
					0,
					"", // TODO: add doc
				),
			},
		},
		parentContext: nil,
	}
	return i
}

func (i *Interpreter) Import(state types.State, newPackagePath string) (
	types.Object,
	error,
) {
	parentPackageInstance, _ := state.GetCurrentPackageOrNil().(*types.Package)
	fullPackagePath, err := getFullPath(newPackagePath, parentPackageInstance)
	if err != nil {
		return nil, err
	}

	package_, err := i.packageStore.GetModule(fullPackagePath)
	if err == nil {
		return package_, nil
	}

	// currPackage := parentPackageInstance
	// for currPackage != nil {
	// 	if currPackage.PackageImpl.Info.Name == fullPackagePath {
	// 		return nil, errors.New("циклічний імпорт заборонений")
	// 	}
	//
	// 	currPackage = currPackage.Parent
	// }

	packageCode, err := readFile(fullPackagePath)
	if err != nil {
		return nil, err
	}

	ast, err := state.GetParser().Parse(fullPackagePath, string(packageCode))
	if err != nil {
		return nil, err
	}

	ctx := i.rootContext.Derive()

	// pkg := types.NewPackageInstance(
	// 	i.rootContext.Derive(),
	// 	fullPackagePath,
	// 	parentPackageInstance,
	// 	nil,
	// )
	// ctx := pkg.GetContext()
	if _, err = ast.Evaluate(state.WithContext(ctx).WithPackage(nil)); err != nil {
		return nil, err
	}

	scope := ctx.TopScope()
	attrs := map[string]types.Object{}
	if toExport, err := ctx.GetVar(builtin.ExportedAttributeName); err == nil {
		switch exported := toExport.(type) {
		case *types.List:
			for _, value := range exported.Items {
				if name, ok := value.(types.String); ok {
					if attr, ok := scope[string(name)]; ok {
						attrs[string(name)] = attr
					}
				}
			}
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

	packageImpl := &types.PackageImpl{
		Info: types.PackageInfo{
			Name:     fullPackagePath,
			Doc:      "", // TODO: add doc
			FileDesc: "",
			Flags:    0,
		},
		Methods:         nil,
		Globals:         attrs,
		OnContextClosed: nil,
	}

	package_, err = i.packageStore.NewPackage(ctx, packageImpl)
	if err != nil {
		return nil, err
	}

	return package_, nil
}

func (i *Interpreter) StackTrace() *types.StackTrace {
	return &i.stacktrace
}

func getFullPath(packagePath string, parentPackage *types.Package) (string, error) {
	if strings.HasPrefix(packagePath, "!/") {
		packagePath = path.Join(os.Getenv(builtin.BORSCH_LIB), packagePath[2:])
	} else if !path.IsAbs(packagePath) {
		var err error
		if parentPackage != nil {
			baseDir := path.Dir(parentPackage.PackageImpl.Info.Name)
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
