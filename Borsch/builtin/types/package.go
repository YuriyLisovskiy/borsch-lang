package types

import "fmt"

var PackageClass = ObjectClass.ClassNew("пакет", map[string]Object{}, true, nil, nil)

type Package struct {
	Filename string
	Dict     map[string]Object
	Parent   *Package
	Context  Context
}

func PackageNew(filename string, parent *Package, ctx Context) *Package {
	return &Package{
		Filename: filename,
		Dict: map[string]Object{
			"__назва__": String(filename),
		},
		Parent:  parent,
		Context: ctx,
	}
}

func (value *Package) Class() *Class {
	return PackageClass
}

func (value *Package) represent(Context) (Object, error) {
	return String(fmt.Sprintf("<пакет %s>", value.Name())), nil
}

func (value *Package) string(ctx Context) (Object, error) {
	return value.represent(ctx)
}

func (value *Package) getAttribute(_ Context, name string) (Object, error) {
	if attr, ok := value.Dict[name]; ok {
		return attr, nil
	}

	if attr := value.Class().GetAttributeOrNil(name); attr != nil {
		return attr, nil
	}

	return nil, ErrorNewf("пакет '%s' не містить атрибута '%s'", value.Name(), name)
}

func (value *Package) setAttribute(_ Context, name string, newValue Object) error {
	attr, ok := value.Dict[name]
	if !ok {
		attr = value.Class().GetAttributeOrNil(name)
	}

	if attr != nil && attr.Class() != newValue.Class() {
		return ErrorNewf(
			"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
			newValue.Class().Name,
			name,
			attr.Class().Name,
		)
	}

	value.Dict[name] = value
	return nil
}

func (value *Package) deleteAttribute(_ Context, name string) (Object, error) {
	if attr, ok := value.Dict[name]; ok {
		delete(value.Dict, name)
		return attr, nil
	}

	if attr := value.Class().DeleteAttributeOrNil(name); attr != nil {
		return attr, nil
	}

	return nil, ErrorNewf("пакет '%s' не містить атрибута '%s'", value.Name(), name)
}

func (value *Package) Name() string {
	name, ok := value.Dict["__назва__"].(String)
	if !ok {
		name = "???"
	}

	return string(name)
}
