package interpolator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

const (
	variable_finder string = `{{[ ]*([a-zA-Z0-9_\.\-\[\]])+[ ]*}}`
)

func Interpolate(str string, vars map[string]interface{}) (interface{}, error) {

	var fnEvaluate func( string,  map[string]interface{}) (interface{}, error)
	fnEvaluate = func(str string, vars map[string]interface{}) (interface{}, error) {
		str = strings.TrimSpace(str[2 : len(str)-2])
		jsonStr, err := json.Marshal(vars)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		} else {
			fmt.Println(string(jsonStr))
		}

		result := gojsonq.New().FromString(string(jsonStr)).Find(str)

		if reflect.TypeOf(result).Name() == "string" {
			var re = regexp.MustCompile(variable_finder)
			if re.Match([]byte(result.(string))) {
				variable := re.FindString(result.(string))
				result, err = fnEvaluate(variable, vars)
			}
		}

		return result, err
	}

	var re = regexp.MustCompile(variable_finder)
	if !re.MatchString(str) {
		return nil, fmt.Errorf("error, parsing the variable %s", str)
	}

	result := str
	for _, value := range re.FindAllString(str, -1) {
		resultLoop, _ := fnEvaluate(value, vars)
		result = strings.Replace(result, value, resultLoop.(string), 1)
	}

	return result, nil
}
