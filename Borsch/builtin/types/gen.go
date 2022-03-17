// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"text/template"
)

const filename = "utility_arithmetic.go"

type Ops []struct {
	Name                string
	Title               string
	Operator            string
	TwoReturnParameters bool
	Unary               bool
	Binary              bool
	Ternary             bool
	NoInplace           bool
	Reversed            string
	Conversion          string
	FailReturn          string
}

type Data struct {
	UnaryOps      Ops
	BinaryOps     Ops
	ComparisonOps Ops
}

var data = Data{
	UnaryOps: Ops{
		// {Name: "neg", Title: "Negate", Operator: "-", Unary: true},
		// {Name: "pos", Title: "MakePositive", Operator: "+", Unary: true},
		// / {Name: "abs", Title: "Abs", Operator: "abs", Unary: true},
		// {Name: "invert", Title: "Invert", Operator: "~", Unary: true},
		// / {Name: "complex", Title: "MakeComplex", Operator: "complex", Unary: true, Conversion: "Complex"},
		// {Name: "int", Title: "MakeInt", Operator: "int", Unary: true, Conversion: "Int"},
		// {Name: "real", Title: "MakeReal", Operator: "real", Unary: true, Conversion: "Real"},
		// / {Name: "iter", Title: "Iter", Operator: "iter", Unary: true},
	},
	BinaryOps: Ops{
		{Name: "add", Title: "Add", Operator: "+", Binary: true},
		// {Name: "sub", Title: "Subtract", Operator: "-", Binary: true},
		// {Name: "mul", Title: "Multiply", Operator: "*", Binary: true},
		// {Name: "div", Title: "Divide", Operator: "/", Binary: true},
		// {Name: "mod", Title: "Mod", Operator: "%%", Binary: true},
		// {Name: "left_shift", Title: "LeftShift", Operator: "<<", Binary: true},
		// {Name: "right_shift", Title: "RightShift", Operator: ">>", Binary: true},
		// {Name: "and", Title: "And", Operator: "&", Binary: true},
		// {Name: "xor", Title: "Xor", Operator: "^", Binary: true},
		// {Name: "or", Title: "Or", Operator: "|", Binary: true},
		// {Name: "pow", Title: "Pow", Operator: "**", Ternary: true},
	},
	ComparisonOps: Ops{
		// {Name: "greater_than", Title: "Greater", Operator: ">", Reversed: "less_than"},
		// {Name: "greater_or_equal", Title: "GreaterOrEqual", Operator: ">=", Reversed: "less_or_equal"},
		// {Name: "less_than", Title: "LessThan", Operator: "<", Reversed: "greater_than"},
		// {Name: "less_or_equal", Title: "LessOrEqual", Operator: "<=", Reversed: "greater_or_equal"},
		// {Name: "equal", Title: "Equal", Operator: "==", Reversed: "equal", FailReturn: "False"},
		// {Name: "not_equal", Title: "NotEqual", Operator: "!=", Reversed: "not_equal", FailReturn: "True"},
	},
}

func main() {
	t := template.Must(template.New("main").Parse(program))
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to open %q: %v", filename, err)
	}
	if err := t.Execute(out, data); err != nil {
		log.Fatal(err)
	}
	err = out.Close()
	if err != nil {
		log.Fatalf("Failed to close %q: %v", filename, err)
	}
	err = exec.Command("go", "fmt", filename).Run()
	if err != nil {
		log.Fatalf("Failed to gofmt %q: %v", filename, err)
	}
}

var program = `// Automatically generated - DO NOT EDIT
// Regenerate with: go generate

// Arithmetic operations

package types

{{ range .UnaryOps }}
// {{.Title}} the Object returning an Object.
//
// Will raise TypeError if {{.Title}} can't be run on this object.
func {{.Title}}(a Object) (Object, error) {
{{ if .Conversion }}
	if _, ok := a.({{.Conversion}}); ok {
		return a, nil
	}
{{end}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res, err := A.__{{.Name}}__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для {{.Operator}}: '%s'", a.Class().Name)
}
{{ end }}

{{ range .BinaryOps }}
// {{.Title}} {{ if .Binary }}two{{ end }}{{ if .Ternary }}three{{ end }} objects together returning an Object.
{{ if .Ternary}}//
// If c != NilTypeClass then it won't attempt to call __reversed_{{.Name}}__
{{ end }}//
// Will raise TypeError if {{.Name}} can't be run on these objects.
func {{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object {{ if .TwoReturnParameters}}, Object{{ end }}, error) {
	// Try using a to {{.Name}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res {{ if .TwoReturnParameters}}, res2{{ end }}, err := A.__{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if err != nil {
			return nil {{ if .TwoReturnParameters }}, nil{{ end }}, err
		}

		if res != NotImplemented {
			return res {{ if .TwoReturnParameters }}, res2{{ end }}, nil
		}
	}

	// Now using b to reversed_{{.Name}} if different in type to a
	if {{ if .Ternary }} c == NilTypeClass && {{ end }} a.Class() != b.Class() {
		if B, ok := b.(I__reversed_{{.Name}}__); ok {
			res {{ if .TwoReturnParameters}}, res2 {{ end }}, err := B.__reversed_{{.Name}}__(a)
			if err != nil {
				return nil {{ if .TwoReturnParameters }}, nil{{ end }}, err
			}

			if res != NotImplemented {
				return res{{ if .TwoReturnParameters}}, res2{{ end }}, nil
			}
		}
	}

	return nil{{ if .TwoReturnParameters}}, nil{{ end }}, ErrorNewf(TypeError, "непідтримувані типи операндів для {{.Operator}}: '%s' та '%s'", a.Class().Name, b.Class().Name)
}

{{ if not .NoInplace }}
func InPlace{{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object, error) {
	if A, ok := a.(I__in_place_{{.Name}}__); ok {
		res, err := A.__in_place_{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return {{.Title}}(a, b {{ if .Ternary }}, c{{ end }})
}
{{end}}
{{end}}

{{ range .ComparisonOps }}
// {{.Title}} two objects returning a boolean result.
//
// Will raise TypeError if {{.Title}} can't be run on this object.
func {{.Title}}(a Object, b Object) (Object, error) {
	// Try using a to {{.Name}}.
	if A, ok := a.(I__{{.Name}}__); ok {
		res, err := A.__{{.Name}}__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to {{.Reversed}} with reversed parameters.
	if B, ok := b.(I__{{.Reversed}}__); ok {
		res, err := B.__{{.Reversed}}__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

{{ if .FailReturn}}
if a.Class() != b.Class() {
	return {{ .FailReturn }}, nil
}
{{ end }}
	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для {{.Operator}}: '%s' та '%s'", a.Class().Name, b.Class().Name)
}
{{ end }}
`
