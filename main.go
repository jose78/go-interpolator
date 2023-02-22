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

var (
	index int = 0
)

func extractKeys(str string) []string {
	var replaceRegexPattern = regexp.MustCompile(`{{|\|(.*?)}}|\.|\}}`)
	var re = regexp.MustCompile(`{{[ ]*.([a-zA-Z\_\-| ]*) [0-9a-zA-Z \[\],.]*[ ]*}}`)
	lstKeys := []string{}
	for _, match := range re.FindAllString(str, -1) {
		keyStracted := replaceRegexPattern.ReplaceAllString(match, "")
		lstKeys = append(lstKeys, strings.TrimSpace(keyStracted))
	}
	return lstKeys
}


 func fnInterpolateString(str string, vars map[string]interface{}) string {
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
	return tmplBytes.String()
}


func executeInterpolator(str string, vars map[string]interface{}, currentKey string, keysEvaluated map[string]string) (string, error) {
	fmt.Printf("\n\n\n*************************************************************\nindex:%v - currentKey:%s, current str:[%s]\n", index, currentKey, str)
	index = index + 1

	result := fnInterpolateString(str, vars)
	lstKeys := extractKeys(result)
	fmt.Printf("generated result:[%s] and list of keys detected:%v", result, lstKeys)
	if len(lstKeys) == 0 {
		vars[currentKey] = result
		lstKeys = extractKeys(str)
		if len(lstKeys) > 0 {
			for _, item := range lstKeys {
				keysEvaluated[item] = ""
				fmt.Printf("current key: %s, key evaluated: %s\n", currentKey, item)
			}
		} else {
			fmt.Printf("str %s evaluated with 0 items, current list: %s\n", str, currentKey)
			result = vars[currentKey].(string)
		}
	} else {
		for _, item := range lstKeys {
			value, ok := keysEvaluated[item]
			if ok || item == currentKey {
				return "", fmt.Errorf("error, cyclic interpolation detected over the key %s", value)
			}
			valueEvaluated, err := executeInterpolator(vars[item].(string), vars, item, keysEvaluated)
			if err != nil {
				return "", fmt.Errorf("error, generated to execute a recursive interpolation using the key '%s' with the content [%s], err: [%v]", item, vars[item].(string), err)
			}
			vars[item] = valueEvaluated
		}
		fmt.Printf("Maps used to be iterpoltaed:%v", vars)
		result = fnInterpolateString(result, vars)
	}
	fmt.Printf("returned result:%s\n", result)
	return result, nil
}

//// Given a string with the templates, it is interpolated with the value of the vars.
func Do(str string, vars map[string]interface{}) (string, error) {
	result, err:= executeInterpolator(str, vars, "", map[string]string{})
	return result, err
}
