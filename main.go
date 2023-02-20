/*
The main goal of interpolator is to help you to interpolate your vars inside the string and evaluate functions related with this vars. It is dased on the 'go templates' and using the large list of functions provided by the Sprig Functions Project.
it is an example ahout how to use it:

	package main

	import "github.com/judoDSL/interpolator"

	func main () {
		values := make(map[string] interface{})
		values["name"] = "            Jose                 "
		values["main_topic"] = "restore the snyderverse"
		values["favorite_superhero"] = "batman who laughs"
		interpolator.Do("I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", values).Println()
	}

It would be the result of the execution:

	[jose78@~/ws/test_interpolator] $  go run main.go
	I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs
*/
package interpolator

import (
	"bytes"
	"fmt"

	//"math/rand"
	"reflect"
	"regexp"

	//"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// Custom type of map[string]interface{}
type VarsContent map[string]interface{}

type content struct {
	value string
}

// Interface that provides the functionality to easily work with the interpolated string
type Content interface {
	// Retrieve the string with the vvalues interpolated.
	Get() string
	// Write to console
	Print()
	// Write to console
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

func ExtractKeys(str string) []string {
	var replaceRegexPattern = regexp.MustCompile(`{{|\|(.*?)}}|\.|\}}`)
	var re = regexp.MustCompile(`{{[ ]*.([a-zA-Z\_\-| ]*) [0-9a-zA-Z \[\],.]*[ ]*}}`)
	lstKeys := []string{}
	for _, match := range re.FindAllString(str, -1) {
		keyStracted := replaceRegexPattern.ReplaceAllString(match, "")
		lstKeys = append(lstKeys, strings.TrimSpace(keyStracted))
	}
	return lstKeys
}

func evaluateVars(mapsContainer map[string]interface{}) {

}

func fnExecuteInterpolator(str string, vars map[string]interface{}, keysEvaluated map[string]string) (string, error) {
	var dataError *string = nil
	funcMap := sprig.FuncMap()
	tmpl, err := template.New("template").Funcs(funcMap).Parse(str)
	if err != nil {
		panic(err)
	}
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	result := tmplBytes.String()
	lstKeys := ExtractKeys(result)
	if len(lstKeys) == 0 {
		lstKeys = ExtractKeys(str)
		for _, item := range lstKeys {
			keysEvaluated[item] = ""
		}
	} else {
		for _, item := range lstKeys {
			value, ok := keysEvaluated[item]
			if ok {
				return *dataError, fmt.Errorf("error, cyclic interpolation detected over the key %s", value)
			}
			valueEvaluated, err := fnExecuteInterpolator(vars[item].(string), vars, keysEvaluated)
			if err != nil {
				return *dataError,  fmt.Errorf("error, generated to execute a recursive interpolation using the key %s with the content %s: %v", item, vars[item].(string), err)
			}
			vars[item] = valueEvaluated
		}
		result, err = fnExecuteInterpolator(result, vars, keysEvaluated)
		if err != nil {
			return *dataError,  fmt.Errorf("error, generated to resolve the text %s: %v", result, err)
			
		}
	}
	return result, nil
}

// Given a string with the templates, it is interpolated with the value of the vars.
func Do(str string, vars map[string]interface{}) (Content, error) {

	var fnEvaluateVars func(internalVars map[string]interface{}, keysEvaluated map[string]string) error
	fnEvaluateVars = func(internalVars map[string]interface{}, keysEvaluated map[string]string) error {
		var re = regexp.MustCompile(`{{[ ]*.[0-9a-zA-Z \[\],.|]+[ ]*}}`)

		var str = "Lo cierto es que en estos momentos estamos en las antípodas {{ .Hola }}  sdasd {{ .Adios | title }} vaya ikaos"

		for i, match := range re.FindAllString(str, -1) {
			fmt.Println(match, "found at index", i)
		}

		for key, value := range internalVars {
			if _, ok := keysEvaluated[key]; ok {
				return fmt.Errorf("error, ciclic evaluation key:[%s] - value:[%s]", key, value)
			}

			if reflect.TypeOf(value).Name() == "string" {
				lstMatches := re.FindAllString(str, -1)
				if len(lstMatches) == 0 {
					keysEvaluated[key] = ""
					internalVars[key] = ""
					//fnExecuteInterpolator(value.(string), internalVars)
				}
			}
		}

		flagNotContentToBeParsed := true
		var trace string = ""

		for flagNotContentToBeParsed {
			flagNotContentToBeParsed = false
			for key, value := range internalVars {
				trace = trace + fmt.Sprintf("Key:[%s] - value:[%s]\n", key, value.(string))
				if reflect.TypeOf(value).Name() == "string" && strings.Contains(value.(string), "{{") {
					if strings.Contains(value.(string), fmt.Sprintf(".%s", key)) {
						return fmt.Errorf("error, cyclic interpolation with key '%s': %s", key, trace)
					}
					result := ""
					//fnExecuteInterpolator(value.(string), internalVars)
					vars[key] = result
					flagNotContentToBeParsed = true
				}
			}
		}
		return nil
	}

	keysEvaluated := make(map[string]string)
	err := fnEvaluateVars(vars, keysEvaluated)
	return content{value: fnExecuteInterpolator(str, vars, nil)}, err
}
