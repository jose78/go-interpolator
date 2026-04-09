package interpolator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

func TestDo2(t *testing.T) {

	customFuncMap := func() template.FuncMap {
		return sprig.FuncMap()
	}

	runner := Configure(Configuration{FnProviderFunction: customFuncMap})

	colorsContent := map[string]interface{}{
		"red":    "rojo",
		"blue":   "azul",
		"pink":   "rosa",
		"orange": "{{ .mapa }}",
	}
	lst := []string{"uno", "dos", "tres"}
	varsContent := map[string]interface{}{
		"the_222":         "la {{ .cosa_rara | title  }}",
		"mix":             "{{ .house }}  {{ .cosa_rara | upper  }} ",
		"house":           "A {{ .the }} casita",
		"the":             "la {{ .cosa_rara | title  }}",
		"animal":          "de {{ .the | title }} mariposa",
		"cosa_rara":       "uno",
		"mapa":            "demo_pato",
		"redirect_pink":   "{{ .colour.pink }}",
		"redirect_orange": "{{ .colour.orange }} {{ .mix }}",
		"cyclic":          "This is a {{ .cyclic }}",
		"lst":             lst,
		"colour": map[string]interface{}{
			"red":    "rojo",
			"blue":   "azul",
			"pink":   "rosa",
			"orange": "{{ .mapa }}",
		},
	}

	values := make(map[string]interface{})
	values["name"] = "            Jose                 "
	values["main_topic"] = "restore the snyderverse"
	values["arms"] = map[string]interface{}{
		"sword": "using the conan sword",
	}
	values["favorite_superhero"] = map[string]interface{}{
		"bad_batman": "batman who laughs with the {{ .arms.sword }}",
	}

	type args struct {
		str  string
		vars map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Must evaluate correctly the key", args{`{{ .the_222  }}`, varsContent}, "la Uno", false},
		{"Must check the function", args{`{{ eq .the   "DEMO_PATO" }}`, varsContent}, false, false},
		////{"Must fail, cyclic", args{"{{ .cyclic | upper }}", varsContent}, "", true},
		{"Must fail, key without dot", args{"{{ cosa_rara | upper }}", varsContent}, "", true},
		{"Must return a simple value", args{"{{ .cosa_rara }}", varsContent}, "uno", false},
		{"Must fail, function not exist", args{"{{ .cosa_rara | floupper }}", varsContent}, "", true},
		{"simple interpolation", args{"{{ .cosa_rara | upper }}", varsContent}, "UNO", false},
		{"medium complex interpolation", args{"{{ .mix }}", varsContent}, "A la Uno casita  UNO ", false},
		{"very complex interpolation", args{"{{ .colour }}", varsContent}, colorsContent, false},
		{"another very complex interpolation", args{"{{ .redirect_pink | upper }}", varsContent}, "ROSA", false},
		{"another very complex interpolation", args{"{{ .redirect_pink | upper }} -- {{ .mix | title }} -- {{ .mix }}", varsContent}, "ROSA -- A La Uno Casita  UNO  -- A la Uno casita  UNO ", false},
		{"another very complex interpolation", args{"{{ .redirect_orange | upper }} {{ .mix | title }} {{ .mix }}", varsContent}, "DEMO_PATO A LA UNO CASITA  UNO  A La Uno Casita  UNO  A la Uno casita  UNO ", false},
		{"connan test", args{"I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero.bad_batman | title }}", values}, "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs With The Using The Conan Sword", false},
		{"connan test", args{"{{ .colour.orange | upper }} -- {{ .colour.orange | title }} -- {{ .colour.orange }} -- {{ .mix }}", varsContent}, "DEMO_PATO -- Demo_pato -- demo_pato -- A la Uno casita  UNO ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runner(tt.args.str, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var flagEq bool
			rt := reflect.TypeOf(got)
			switch rt.Kind() {
			case reflect.Slice:
				flagEq = reflect.DeepEqual(got, tt.want)
			case reflect.Array:
				flagEq = reflect.DeepEqual(got, tt.want)
			case reflect.Map:
				flagEq = reflect.DeepEqual(got, tt.want)
			default:
				flagEq = got == tt.want
			}

			if !flagEq {
				t.Errorf("Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractKeys(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		//{"Should extract as keys the name of the variables", args{"Hola {{ .user_name }} como estás, lo cierto es que esto es {{ .insult }}"}, []string{".user_name", ".insult"}},
		//{"Should extract as keys the name of the variables using also functions ", args{"Hola {{ .user_name | upper }} como estás, lo cierto es que esto es {{ .insult | title}}"}, []string{".user_name", ".insult"}},
		{"Should extract as keys the name of the variables using also functions ", args{`title .insult "hola | | | ma, nsdsds" | tesss`}, []string{".insult"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractKeys(tt.args.str); !reflect.DeepEqual(got.Keys, tt.want) {
				t.Errorf("extractKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appensJsonContent(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name                 string
		args                 args
		wantResult           string
		wantFlagContainsJson bool
	}{
		{"Should return the same string with true flag ", args{"{{ funcion .como  | estoesunafuncion | toJson }}"}, "{{ funcion .como  | estoesunafuncion | toJson }}", true},
		{"Should return same string with the sufix of '| toJson ' with false flag ", args{"{{ funcion .como  | estoesunafuncion }}"}, "{{ funcion .como  | estoesunafuncion  | toJson }}", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotFlagContainsJson := appensJsonContent(tt.args.str)
			if gotResult != tt.wantResult {
				t.Errorf("appensJsonContent() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotFlagContainsJson != tt.wantFlagContainsJson {
				t.Errorf("appensJsonContent() gotFlagContainsJson = %v, want %v", gotFlagContainsJson, tt.wantFlagContainsJson)
			}
		})
	}
}

func Test_extractParamteres(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"extract the parameters contained within the string", args{`para empezar esto {{ funcion .como  | estoesunafuncion | toJson }} esto es unba prueba {{ Hola_2 como estas  }}`}, []string{"{{ funcion .como  | estoesunafuncion | toJson }}", "{{ Hola_2 como estas  }}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractParamteres(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamteres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DoRangeTemplate(t *testing.T) {
	customFuncMap := func() template.FuncMap {
		return sprig.FuncMap()
	}

	runner := Configure(Configuration{FnProviderFunction: customFuncMap})

	vars := map[string]interface{}{
		"app": map[string]interface{}{
			"name":     "demo",
			"replicas": 2,
		},
		"containers": []interface{}{
			map[string]interface{}{
				"name":  "c1",
				"image": "img1",
				"ports": []interface{}{
					map[string]interface{}{
						"containerPort": 8080,
						"protocol":      "TCP",
						"name":          "http",
					},
				},
			},
		},
	}

	templateStr := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .app.name }}
spec:
  replicas: {{ .app.replicas }}
  template:
    spec:
      containers:
      {{- range .containers }}
        - name: {{ .name }}
          image: {{ .image }}
          ports:
          {{- range .ports }}
            - containerPort: {{ .containerPort }}
              protocol: {{ .protocol }}
              name: {{ .name }}
          {{- end }}
      {{- end }}`

	got, err := runner(templateStr, vars)
	if err != nil {
		t.Fatalf("runner error: %v", err)
	}

	gotStr, ok := got.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", got)
	}

	if !strings.Contains(gotStr, "apiVersion: apps/v1") ||
		!strings.Contains(gotStr, "name: demo") ||
		!strings.Contains(gotStr, "replicas: 2") ||
		!strings.Contains(gotStr, "- name: c1") ||
		!strings.Contains(gotStr, "image: img1") ||
		!strings.Contains(gotStr, "containerPort: 8080") {
		t.Fatalf("unexpected rendered result:\n%s", gotStr)
	}
}

func Test_DoControlStructures(t *testing.T) {
	customFuncMap := func() template.FuncMap {
		return sprig.FuncMap()
	}

	runner := Configure(Configuration{FnProviderFunction: customFuncMap})
	vars := map[string]interface{}{
		"enabled": true,
		"count":   2,
		"parent": map[string]interface{}{
			"child": "x",
		},
		"name": "bob",
		"items": []interface{}{
			map[string]interface{}{"name": "one"},
			map[string]interface{}{"name": "two"},
		},
	}

	tests := []struct {
		name string
		tmpl string
		want string
	}{
		{"if true", `{{ if .enabled }}ok{{ end }}`, "ok"},
		{"if false else", `{{ if .missing }}yes{{ else }}no{{ end }}`, "no"},
		{"else if", `{{ if eq .count 1 }}one{{ else if eq .count 2 }}two{{ else }}many{{ end }}`, "two"},
		{"with block", `{{ with .parent }}child={{ .child }}{{ end }}`, "child=x"},
		{"define and template", `{{ define "inner" }}Hello {{ .name }}{{ end }}{{ template "inner" . }}`, "Hello bob"},
		{"block default", `{{ define "outer" }}prefix {{ block "inner" . }}inner-default{{ end }} suffix{{ end }}{{ template "outer" . }}`, "prefix inner-default suffix"},
		{"variable assignment", `{{ $greet := "Hi" }}{{ $greet }}`, "Hi"},
		{"nested range", `{{ range .items }}{{ .name }}-{{ end }}`, "one-two-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runner(tt.tmpl, vars)
			if err != nil {
				t.Fatalf("runner error: %v", err)
			}
			gotStr, ok := got.(string)
			if !ok {
				t.Fatalf("expected string result, got %T", got)
			}
			if gotStr != tt.want {
				t.Fatalf("template %q = %q, want %q", tt.name, gotStr, tt.want)
			}
		})
	}
}

func Test_SuperHero(t *testing.T) {
	customFuncMap := func() template.FuncMap {
		return sprig.FuncMap()
	}

	runner := Configure(Configuration{FnProviderFunction: customFuncMap})

	values := make(map[string]interface{})
	values["name"] = "            Jose                 "
	values["main_topic"] = "restore the snyderverse"
	values["hero"] = "{{ .hero_redirect }}"
	values["favorite_superhero"] = "{{ .hero | upper }} who laughs"
	values["hero_redirect"] = "batman"
	str, _ := runner("I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", values)

	resultExpected := "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with BATMAN Who Laughs"
	if str != resultExpected {
		t.Errorf("extractParamteres() = %v, want %v", str, resultExpected)
	}

	fmt.Println(str)
}
