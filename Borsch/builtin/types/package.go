// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
//
// Package objects.

package types

import (
	"fmt"
	"sync"
)

type PackageFlags int32

const (
	// ShareModule signals that an embedded module is threadsafe and read-only, meaninging it could be shared across multiple py.Context instances (for efficiency).
	// Otherwise, PackageImpl will create a separate py.Package instance for each py.Context that imports it.
	// This should be used with extreme caution since any module mutation (write) means possible cross-context data corruption.
	ShareModule PackageFlags = 0x01

	MainModuleName = "__main__"
)

// PackageInfo contains info and about a module and can specify flags that affect how it is imported into a py.Context
type PackageInfo struct {
	Name     string // __name__ (if nil, "__main__" is used)
	Doc      string // __doc__
	FileDesc string // __file__
	Flags    PackageFlags
}

// PackageImpl is used for modules that are ready to be imported into a py.Context.
// The model is that a PackageImpl is read-only and instantiates a Package into a py.Context when imported.
//
// By convention, .Code is executed when a module instance is initialized. If nil,
// then .CodeBuf or .CodeSrc will be auto-compiled to set .Code.
type PackageImpl struct {
	Info            PackageInfo
	Methods         []*Method      // Package-bound global method functions
	Globals         Dict           // Package-bound global variables
	OnContextClosed func(*Package) // Callback for when a py.Context is closing to release resources
}

// PackageStore is a container of Package imported into an owning py.Context.
type PackageStore struct {
	// Registry of installed modules
	modules map[string]*Package
	// Builtin module
	Builtins *Package
	// this should be the frozen module importlib/_bootstrap.py generated
	// by Modules/_freeze_importlib.c into Python/importlib.h
	Importlib *Package
}

func RegisterModule(module *PackageImpl) {
	gRuntime.RegisterModule(module)
}

func GetModuleImpl(moduleName string) *PackageImpl {
	gRuntime.mu.RLock()
	defer gRuntime.mu.RUnlock()
	impl := gRuntime.PackageImpls[moduleName]
	return impl
}

type Runtime struct {
	mu           sync.RWMutex
	PackageImpls map[string]*PackageImpl
}

var gRuntime = Runtime{
	PackageImpls: make(map[string]*PackageImpl),
}

func (rt *Runtime) RegisterModule(impl *PackageImpl) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.PackageImpls[impl.Info.Name] = impl
}

func NewModuleStore() *PackageStore {
	return &PackageStore{
		modules: make(map[string]*Package),
	}
}

// Package is a runtime instance of a PackageImpl bound to the py.Context that imported it.
type Package struct {
	PackageImpl *PackageImpl // Parent implementation of this Package instance
	Globals     Dict         // Initialized from PackageImpl.Globals
	Context     Context      // Parent context that "owns" this Package instance
}

var PackageClass = NewClass("пакет", "об_єкт пакету")

func (value *Package) Class() *Class {
	return PackageClass
}

func (value *Package) __represent__() (Object, error) {
	name, ok := value.Globals["__name__"].(String)
	if !ok {
		name = "???"
	}

	return String(fmt.Sprintf("<module %s>", string(name))), nil
}

func (value *Package) GetDict() Dict {
	return value.Globals
}

// Call calls a named method of a module.
func (value *Package) Call(state State, name string, args Tuple) (Object, error) {
	attr, err := GetAttrString(value, name)
	if err != nil {
		return nil, err
	}

	return Call(state, attr, args)
}

// Interfaces
var _ IGetDict = (*Package)(nil)

// NewPackage adds a new Package instance to this PackageStore.
// Each given Method prototype is used to create a new "live" Method bound this the newly created Package.
// This func also sets appropriate module global attribs based on the given PackageInfo (e.g. __name__).
func (store *PackageStore) NewPackage(ctx Context, impl *PackageImpl) (*Package, error) {
	name := impl.Info.Name
	if name == "" {
		name = MainModuleName
	}

	m := &Package{
		PackageImpl: impl,
		Globals:     impl.Globals.Copy(),
		Context:     ctx,
	}

	// Insert the methods into the module dictionary
	// Copy each method an insert each "live" with a ptr back to the module (which can also lead us to the host Context)
	for _, method := range impl.Methods {
		methodInst := new(Method)
		*methodInst = *method
		methodInst.Package = m
		m.Globals[method.Name] = methodInst
	}

	// Set some module globals
	m.Globals["__name__"] = String(name)
	m.Globals["__doc__"] = String(impl.Info.Doc)
	m.Globals["__package__"] = Nil
	if len(impl.Info.FileDesc) > 0 {
		m.Globals["__file__"] = String(impl.Info.FileDesc)
	}

	// Register the module
	store.modules[name] = m
	// Make a note of some modules
	switch name {
	case "builtins":
		store.Builtins = m
	case "importlib":
		store.Importlib = m
	}

	// fmt.Printf("Registered module %q\n", moduleName)
	return m, nil
}

func (store *PackageStore) GetModule(name string) (*Package, error) {
	m, ok := store.modules[name]
	if !ok {
		return nil, ErrorNewf(ImportError, "Пакет '%s' не знайдено", name)
	}

	return m, nil
}

func (store *PackageStore) MustGetModule(name string) *Package {
	m, err := store.GetModule(name)
	if err != nil {
		panic(err)
	}
	return m
}

// OnContextClosed signals all module instances that the parent py.Context has closed
func (store *PackageStore) OnContextClosed() {
	for _, m := range store.modules {
		if m.PackageImpl.OnContextClosed != nil {
			m.PackageImpl.OnContextClosed(m)
		}
	}
}
