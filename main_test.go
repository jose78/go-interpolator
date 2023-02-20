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
	"reflect"
	"testing"
)

func TestDo(t *testing.T) {

	a := map[string]interface{}{
		"name":               "            Jose                 ",
		"main_topic":         "restore the snyderverse",
		"favorite_superhero": "batman who laughs",
		"the":                "the",
	}

	//b := map[string]interface{}{
	//	"house":     "A {{ .the }} casita ",
	//	"colour":    "rosa,",
	//	"the":       "la {{ .cosa_rara}}",
	//	"animal":    "de {{ .the }} mariposa",
	//	"cosa_rara": "demo_pato {{ .the }}",
	//}

	c := map[string]interface{}{
		"house":     "A {{ .the | title }} casita ",
		"colour":    "rosa,",
		"the":       "la {{ .cosa_rara | title  }}",
		"animal":    "de {{ .the | title }} mariposa",
		"cosa_rara": "demo_pato",
	}

	type args struct {
		str  string
		vars map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Content
		wantErr bool
	}{
		{"Should check the content generate by the interpolator is correct",
			args{"I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", a},
			content{value: "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs"},
			false,
		}, {"Should interpolate all nested templates",
			args{" {{ .house }} {{ .colour }} {{ .animal }} {{ .cosa_rara }}", c},
			content{" A La Demo_pato casita  rosa, de La Demo_pato mariposa demo_pato"},
			false},
		//{"should fail",
		//	args{" {{ .house }} {{ .colour }} {{ .the }} {{ .animal }} {{ .cosa_rara}}", b},
		//	content{" A la demo_pato la demo_pato {{ .the }} casita  rosa, la demo_pato {{ .the }} de la demo_pato {{ .the }} mariposa demo_pato la demo_pato {{ .the }}"},
		//	true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.str, tt.args.vars)
			if tt.wantErr {
				if (err != nil) != tt.wantErr {
					t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Do() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestExtractKeys(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Should extrac the keys", args{"Lo más importante en la {{ .vida }} es poder decir lo que {{ .quieres | alto }}"}, []string{"vida", "quieres"}},
		{"Should extrac the keys", args{"la {{ .cosa_rara | title  }}"}, []string{"cosa_rara"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractKeys(tt.args.str); 
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evaluateVars(t *testing.T) {

	type args struct {
		mapsContainer map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluateVars(tt.args.mapsContainer)
		})
	}
}

func Test_fnExecuteInterpolator(t *testing.T) {

	c := map[string]interface{}{
		"house":     "A {{ .the | title }} casita ",
		"colour":    "rosa,",
		"the":       "la {{ .cosa_rara | title  }}",
		"animal":    "de {{ .the | title }} mariposa",
		"cosa_rara": "demo_pato",
		"mapa": "demo_pato",
	}

	type args struct {
		str  string
		vars map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should interpolate this", args{"{{ .house }}",c}, "A La Demo_pato casita "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fnExecuteInterpolator(tt.args.str, tt.args.vars, map[string]string{}); 
			if got != tt.want {
				t.Errorf("fnExecuteInterpolator() = %v, want %v", got, tt.want)
			}
		})
	}
}
