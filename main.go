package judo_interpolator

import (
	"bytes"
	"fmt"
	"text/template"
	"github.com/Masterminds/sprig/v3"
)

type content struct {
	value string
}
// Provide the funcionality to work easily with the string intepolated
type Content interface {
	// Retrieve the string with the vvalues interpolated.
	get() string
	// Write at console
	print()
	// Write at console
	println()
	// Create an error
	error() error
}


func (body content) get() string {
	return body.value
}

func (body content) print() {
	fmt.Print(body.value)
}

func (body content) println() {
	fmt.Println(body.value)
}

func (body content) error() error {
	return fmt.Errorf(body.value)
}

// Given a string with the templates, it is interpolated with the value of the vars.
func Do(str string, vars map[string]interface{}) Content {
	funcMap := sprig.FuncMap()
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(str)
	if err != nil {
		panic(err)
	}
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return content{value: tmplBytes.String()}
}
