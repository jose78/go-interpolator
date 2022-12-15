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
	Get() string
	// Write at console
	Print()
	// Write at console
	Println()
	// Create an error
	Error() error
}


func (body content) Get() string {
	return body.value
}

func (body content) Print() {
	fmt.Print(body.value)
}

func (body content) Println() {
	fmt.Println(body.value)
}

func (body content) Error() error {
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
