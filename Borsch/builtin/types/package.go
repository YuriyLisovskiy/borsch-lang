package types

import "fmt"

var PackageClass = ObjectClass.ClassNew("пакет", map[string]Object{}, true, nil, nil)

var (
	Initialized = false
)

type Package struct {
	Filename string
	Dict     StringDict
	Parent   *Package
	Context  Context
}

func PackageNew(filename string, parent *Package, ctx Context) *Package {
	pkg := &Package{
		Filename: filename,
		Dict: map[string]Object{
			"__назва__": String(filename),
		},
		Parent:  parent,
		Context: ctx,
	}
	pkg.init()
	return pkg
}

func (value *Package) Class() *Class {
	return PackageClass
}

func (value *Package) init() {
	initInstance(value, &value.Dict, value.Class())
}

func (value *Package) represent(Context) (Object, error) {
	return String(fmt.Sprintf("<пакет %s>", value.Name())), nil
}

func (value *Package) string(ctx Context) (Object, error) {
	return value.represent(ctx)
}

func (value *Package) getAttribute(_ Context, name string) (Object, error) {
	return getAttributeFrom(&value.Dict, name, value.Class())
}

func (value *Package) setAttribute(_ Context, name string, newValue Object) error {
	attr, ok := value.Dict[name]
	if !ok {
		attr = value.Class().GetAttributeOrNil(name)
	}

	return setAttributeTo(value, &value.Dict, attr, name, newValue)
}

func (value *Package) deleteAttribute(_ Context, name string) (Object, error) {
	if attr, ok := value.Dict[name]; ok {
		delete(value.Dict, name)
		return attr, nil
	}

	if attr := value.Class().DeleteAttributeOrNil(name); attr != nil {
		return attr, nil
	}

	return nil, NewErrorf("пакет '%s' не містить атрибута '%s'", value.Name(), name)
}

func (value *Package) Name() string {
	name, ok := value.Dict["__назва__"].(String)
	if !ok {
		name = "???"
	}

	return fmt.Sprintf("<пакет \"%s\">", name)
}
