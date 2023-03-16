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
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func extractKeys(str string) []string {
	var replaceRegexPattern = regexp.MustCompile(`{{|\|(.*?)}}|\}}`)
	var re = regexp.MustCompile(`{{[ ]*.([a-zA-Z\_\-|. ]*) [0-9a-zA-Z \[\],.]*[ ]*}}`)
	lstKeys := []string{}
	for _, match := range re.FindAllString(str, -1) {
		keyStracted := replaceRegexPattern.ReplaceAllString(match, "")
		lstKeys = append(lstKeys, strings.TrimSpace(keyStracted))
	}
	return lstKeys
}

func interpolateString(str string, vars map[string]interface{}, mapResults map[string]interface{}) (string, error) {
	eval := func(strToInterpolate string) (result string, err error) {
		lstKeys := extractKeys(strToInterpolate)
		//if _, ok := mapResults[strToInterpolate]; ok{
		//	return "", fmt.Errorf("error, cyclic interpolation. The string [%s] has been generated previopusly, keys:%v", strToInterpolate,lstKeys)
		//}
		mapResults[strToInterpolate] = ""
		
		result = appendEval(strToInterpolate, lstKeys)
		result, err = interpolateString(result, vars, mapResults)
		return result, err
	}
	funcMap := sprig.FuncMap()
	funcMap["eval"] = eval
	tmpl, err := template.New("template").Funcs(funcMap).Parse(str)
	if err != nil {
		return "", fmt.Errorf ("error, parsing the next string %s:%v",str, err)
	}


	tmpl.Option()
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, vars)
	if err != nil {
		return "", fmt.Errorf ("error, applying the values over the string %s:%v",str, err)
	}
	return tmplBytes.String(), nil
}

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
			fmt.Printf("TODO: verificar si el valor asociado a este element es una list o un dict %s \n", match)
			str = strings.Replace(str, match, fmt.Sprintf("%s | eval", match), value)
		} else {
			panic("there are a gost key")
		}
	}
	return str
}

// // Given a string with the templates, it is interpolated with the value of the vars.
func Do(str string, vars map[string]interface{}) (result string, err error) {
	result = str
	flagAskToResolveInterpolation := false
	mapResults := map[string]interface{}{}

	for !flagAskToResolveInterpolation {
		lstKeys := extractKeys(result)
		flagAskToResolveInterpolation = len(lstKeys) == 0
		if !flagAskToResolveInterpolation {
			result = appendEval(result, lstKeys)
		}
		result, err = interpolateString(result, vars, mapResults)
		flagAskToResolveInterpolation = flagAskToResolveInterpolation || err != nil

	}
	return result, err
}
