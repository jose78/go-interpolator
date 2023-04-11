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
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// Given a prase, it extracts all parameters to be interpolated\. for instance, given "Hello I'm a {{ function_name .param_1 | function_name_2 }}  and test {{ function_name_3 .param_2 | function_name_3 }} " it should return an array with the items {{ function_name .param_1 | function_name_2 }} and {{ function_name_3 .param_2 | function_name_3 }}
func extractParamteres(str string) []string {
	//var str = `{{ funcion .como "estas" | esto es una funcion | to_json }} esto es unba prueba {{ Hola_2 como estas  }}`
	result := []string{}
	re := regexp.MustCompile(EXTRACT_PARAMTERES)

	for _, match := range re.FindAllString(str, -1) {
		result = append(result, match)
	}
	return result
}

type parameter struct {
	OriginalStr        string
	Paramter           string
	FirstItem          string
	Function           string
	FlagContainsToJson bool
	Keys               []string
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

	functionName := ""
	if strings.Contains(strSplited[0], " ") {
		functionName = strings.Split(strSplited[0][2:], " ")[0]
	}

	return parameter{originalStr, str, strSplited[0], functionName, flagContainsJson, lstKeys}
}

// Given a parameter it check if the function to_json is contained at last position, if not then it will append.
func appensJsonContent(str string) (result string, flagContainsJson bool) {

	result = str[2 : len(str)-2]
	resultSplited := strings.Split(result, "|")

	lastItem := resultSplited[len(resultSplited)-1]

	flagContainsJson = strings.TrimSpace(lastItem) == TO_JSON

	if !flagContainsJson {
		result = fmt.Sprintf("{{%s | %s }}", result, TO_JSON)
	}

	return
}

const (
	EXTRACT_PARAMTERES = `(?m){{\s* ([a-zA-Z0-9."'|_-]* )+\s*}}`
	EXTRACT_KEYS       = `(?m)\s*(\.[a-zA-Z0-1\._-]+)\s*`
	TO_JSON            = "toJson"
)

func Do2(str string, vars map[string]interface{}) (result interface{}, err error) {

	parameters := extractParamteres(str)

	if len(parameters) == 1 {
		parameter := extractKeys(parameters[0])
		result, err = execute(parameter, vars)
		if reflect.TypeOf(result).Name() == "string" ||
			len(strings.TrimSpace(strings.Replace(str, parameter.OriginalStr, "", -1))) != 0 {
			result = strings.Replace(str, parameter.OriginalStr, result.(string), 1)
		}

		return result, err
		// verificar si aparte del parámetro hay máß cosas... si las hay habría que meterlas
	} else {

		for _, param := range parameters {
			parameter := extractKeys(param)
			result, err = execute(parameter, vars)
			if err != nil {
				return nil, fmt.Errorf("error, executing the interpolation of %s: %v", parameter.Paramter, err)
			}
			str = strings.Replace(str, parameter.OriginalStr, result.(string), 1)
		}
		result = str
	}
	fmt.Println(parameters)
	return
}

func execute(param parameter, vars map[string]interface{}) (interface{}, error) {

	mainStr := param.Paramter
	for _, item := range param.Keys {
		mainStr = strings.Replace(mainStr, item, fmt.Sprintf(" ( %s | eval) ", item), 1)
	}

	eval := func(strToInterpolate interface{}) (result interface{}, err error) {
		if reflect.TypeOf(strToInterpolate).Name() == "string" {
			lstKeys := extractKeys(strToInterpolate.(string))
			if len(lstKeys.Keys) > 0 {
				return Do2(strToInterpolate.(string), vars)
			}
		}
		return strToInterpolate, err
	}

	funcMap := sprig.FuncMap()
	funcMap["eval"] = eval

	tmpl, err := template.New("template").Funcs(funcMap).Parse(mainStr)
	if err != nil {
		return "", fmt.Errorf("error, parsing the next string %s:%v", mainStr, err)
	}
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

//func interpolateString(str string, vars map[string]interface{}, mapResults map[string]interface{}) (string, error) {
//	eval := func(strToInterpolate string) (result string, err error) {
//		lstKeys := extractKeys(strToInterpolate)
//		//if _, ok := mapResults[strToInterpolate]; ok{
//		//	return "", fmt.Errorf("error, cyclic interpolation. The string [%s] has been generated previopusly, keys:%v", strToInterpolate,lstKeys)
//		//}
//		mapResults[strToInterpolate] = ""
//
//		result = appendEval(strToInterpolate, lstKeys)
//		result, err = interpolateString(result, vars, mapResults)
//		return result, err
//	}
//	funcMap := sprig.FuncMap()
//	funcMap["eval"] = eval
//	tmpl, err := template.New("template").Funcs(funcMap).Parse(str)
//	if err != nil {
//		return "", fmt.Errorf("error, parsing the next string %s:%v", str, err)
//	}
//	var tmplBytes bytes.Buffer
//	err = tmpl.Execute(&tmplBytes, vars)
//	if err != nil {
//		return "", fmt.Errorf("error, applying the values over the string %s:%v", str, err)
//	}
//	return tmplBytes.String(), nil
//}

func appendEval(str string, lstKeys []string) string {
	freq := make(map[string]int)
	for _, key := range lstKeys {
		freq[key] = freq[key] + 1
	}
	for key, value := range freq {
		pattern := fmt.Sprintf("{{[ ]+%s", key)
		var re = regexp.MustCompile(pattern)
		lstMatched := re.FindAllString(str, -1)
		if len(lstMatched) > 0 {
			match := lstMatched[0]
			str = strings.Replace(str, match, fmt.Sprintf("%s | eval", match), value)
		} else {
			panic("there are a gost key")
		}
	}
	return str
}

// // Given a string with the templates, it is interpolated with the value of the vars.
//func Do(str string, vars map[string]interface{}) (result string, err error) {
//	result = str
//	flagAskToResolveInterpolation := false
//	mapResults := map[string]interface{}{}
//
//	for !flagAskToResolveInterpolation {
//		lstKeys := extractKeys(result)
//		flagAskToResolveInterpolation = len(lstKeys) == 0
//		if !flagAskToResolveInterpolation {
//			result = appendEval(result, lstKeys)
//		}
//		result, err = interpolateString(result, vars, mapResults)
//		flagAskToResolveInterpolation = flagAskToResolveInterpolation || err != nil
//
//	}
//	return result, err
//}
//
