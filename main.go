/*
The main goal of interpolator is to help you to interpolate your vars inside the string and evaluate functions related with this vars. It is dased on the 'go templates' and using the large list of functions provided by the Sprig Functions Project.
it is an example ahout how to use it:

	package main

	import "github.com/judoctl/interpolator"

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
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	EXTRACT_PARAMTERES = `(?m){{\s* ([a-zA-Z0-9."'|_-]* )+\s*}}`
	EXTRACT_KEYS       = `(?m)\s*(\.[a-zA-Z0-9\._-]+)\s*`
	TO_JSON            = "toJson"
)

var (
	funcMap = getFunctions
)
type parameter struct {
	originalStr        string
	paramter           string
	FlagContainsToJson bool
	Keys               []string
}

// Given a prase, it extracts all parameters to be interpolated\. for instance, given "Hello I'm a {{ function_name .param_1 | function_name_2 }}  and test {{ function_name_3 .param_2 | function_name_3 }} " it should return an array with the items {{ function_name .param_1 | function_name_2 }} and {{ function_name_3 .param_2 | function_name_3 }}
func extractParamteres(str string) []string {
	result := []string{}
	re := regexp.MustCompile(EXTRACT_PARAMTERES)
	result = append(result, re.FindAllString(str, -1)...)
	return result
}

// Given a parameter, it retrieves the keys of the parameter. For instance, given {{ function_name .param_1 | function_name_2 }} it should return an array with the value of .param_1
func extractKeys(str string) parameter {

	replacePipesInQuotes := func(input string) string {
		re := regexp.MustCompile(`"([^"\\]*(?:\\.[^"\\]*)*)"`)
		return re.ReplaceAllStringFunc(input, func(s string) string {
			result := strings.ReplaceAll(s, ",", "_______________COMMA_______________")
			return strings.ReplaceAll(result, "|", "_______________PIPE_______________")
		})
	}
	var flagContainsJson bool
	originalStr := str
	str, flagContainsJson = appensJsonContent(str)
	str = replacePipesInQuotes(str)
	strSplited := strings.Split(str, "|")

	var re = regexp.MustCompile(EXTRACT_KEYS)
	lstKeys := []string{}
	for _, match := range re.FindAllString(strSplited[0], -1) {
		lstKeys = append(lstKeys, strings.TrimSpace(match))
	}

	return parameter{originalStr, str, flagContainsJson, lstKeys}
}

// Given a parameter it check if the function to_json is contained at last position, if not then it will append.
func appensJsonContent(str string) (result string, flagContainsJson bool) {

	result = str
	if strings.HasPrefix(str, "{{") {
		result = str[2 : len(str)-2]
	} else {
		if !strings.HasPrefix(strings.TrimSpace(str), ".") {
			result = fmt.Sprintf("\"%s\"", str)
		}
	}
	resultSplited := strings.Split(result, "|")

	lastItem := resultSplited[len(resultSplited)-1]

	flagContainsJson = strings.TrimSpace(lastItem) == TO_JSON

	if !flagContainsJson {
		result = fmt.Sprintf("{{%s | %s }}", result, TO_JSON)
	} else {
		result = fmt.Sprintf("{{%s}}", result)
	}

	return
}

func Do(str string, vars map[string]interface{}) (result interface{}, err error) {

	parameters := extractParamteres(str)

	if len(parameters) == 1 {
		parameter := extractKeys(parameters[0])
		result, err = execute(parameter, vars)
		if reflect.TypeOf(result).Name() == "string" ||
			len(strings.TrimSpace(strings.Replace(str, parameter.originalStr, "", -1))) != 0 {
			result = strings.Replace(str, parameter.originalStr, result.(string), 1)
		}

		return result, err
		// verificar si aparte del parámetro hay máß cosas... si las hay habría que meterlas
	} else {

		for _, param := range parameters {
			parameter := extractKeys(param)
			result, err = execute(parameter, vars)
			if err != nil {
				return nil, fmt.Errorf("error, executing the interpolation of %s: %v", parameter.paramter, err)
			}
			str = strings.Replace(str, parameter.originalStr, result.(string), 1)
		}
		result = str
	}
	return
}

func generateName() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 15)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	str := "template_" + string(b)
	return str
}

func execute(param parameter, vars map[string]interface{}) (interface{}, error) {

	mainStr := param.paramter
	for _, item := range param.Keys {
		mainStr = strings.Replace(mainStr, item, fmt.Sprintf(" ( %s | eval) ", item), 1)
	}

	eval := func(strToInterpolate interface{}) (result interface{}, err error) {
		if reflect.TypeOf(strToInterpolate).Name() == "string" {
			lstKeys := extractKeys(strToInterpolate.(string))
			if len(lstKeys.Keys) > 0 {
				return Do(strToInterpolate.(string), vars)
			}
		}
		return strToInterpolate, err
	}

	lstFunctionsMapp := funcMap()
	lstFunctionsMapp["eval"] = eval

	tmpl, err := template.New(generateName()).Funcs(lstFunctionsMapp).Parse(mainStr)
	if err != nil {
		return "", fmt.Errorf("error, parsing the next string %s:%v", mainStr, err)
	}

	tmpl.Option()
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, vars)
	if err != nil {
		return "", fmt.Errorf("error, applying the values over the string %s:%v", mainStr, err)
	}

	var result interface{} = tmplBytes.String()

	if !param.FlagContainsToJson {

		json.Unmarshal(tmplBytes.Bytes(), &result)
	}

	return result, nil

}


type typeValidateFunc func(str string, vars map[string]interface{}) (result interface{}, err error)

// Type of function getFunctions, to use your custom functions
type TypeProviderFunctions func() template.FuncMap   

// Data to overwrithe the default behavior, it must be set through the configuration function
type Configuration struct {
	FnProviderFunction    TypeProviderFunctions // Update the list of functions to be used within the templates
}


// Configure optional values of struts_validation:
// funcMap: Function to set de defailt list of funcMap to be used during the template process.
func Configure(conf Configuration) typeValidateFunc {
	if conf.FnProviderFunction != nil {
		funcMap = conf.FnProviderFunction
	}
	return Do
}

// Generate the default empty funcMaps to be used
func getFunctions() template.FuncMap {
	fnMap := template.FuncMap{}
	return fnMap
}
