package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	EXTRACT_PARAMTERES = `(?m){{\s* ([a-zA-Z0-9."'|_-]* )+\s*}}`
	TO_JSON            = "to_json"
)

func Do2(str string, vars map[string]interface{}) (result interface{}, err error) {

	parameters := extractParamteres(str)

	for _, param := range parameters {
		paramUpdated, flagParamContainToJson := appensJsonContent(param)
		fmt.Println(paramUpdated)
		fmt.Println(flagParamContainToJson)

	}

	fmt.Println(parameters)
	return

}

func appensJsonContent(str string) (result string, flagContainsJson bool) {

	result = str[2 : len(str)-2]
	resultSplited := strings.Split(result, "|")

	lastItem := resultSplited[len(resultSplited)-1]

	flagContainsJson = strings.TrimSpace(lastItem) == TO_JSON

	if !flagContainsJson {
		result = fmt.Sprintf("%s | %s ", result, TO_JSON)
	}

	return
}


