package interpolator

import (
	"reflect"
	"testing"
)

func TestDo2(t *testing.T) {

	varsContent := map[string]interface{}{
		"mix":             "{{ .house }}  {{ .cosa_rara | upper  }} ",
		"house":           "A {{ .the }} casita",
		"the":             "la {{ .cosa_rara | title  }}",
		"animal":          "de {{ .the | title }} mariposa",
		"cosa_rara":       "demo_pato",
		"mapa":            "demo_pato",
		"redirect_pink":   "{{ .colour.pink }}",
		"redirect_orange": "{{ .colour.orange }} {{ .mix }}",
		"cyclic":          "This is a {{ .cyclic }}",
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
		{"Must check the function", args{`{{ eq .the   "DEMO_PATO" }}`, varsContent}, false, false},
		// {"Must fail, cyclic", args{"{{ .cyclic | upper }}", varsContent}, "", true},
		//{"Must fail, key without dot", args{"{{ cosa_rara | upper }}", varsContent}, "", true},
		//{"Must fail, function not exist", args{"{{ .cosa_rara | floupper }}", varsContent}, "", true},
		//{"simple interpolation", args{"{{ .cosa_rara | upper }}", varsContent}, "DEMO_PATO", false},
		//{"medium complex interpolation", args{"{{ .mix }}", varsContent}, "A la Demo_pato casita  DEMO_PATO ", false},
		//{"very complex interpolation", args{"{{ .colour.orange |  upper }}", varsContent}, "DEMO_PATO", false},
		//{"another very complex interpolation", args{"{{ .redirect_pink | upper }}", varsContent}, "ROSA", false},
		//{"another very complex interpolation", args{"{{ .redirect_pink | upper }} -- {{ .mix | title }} -- {{ .mix }}", varsContent}, "ROSA -- A La Demo_pato Casita  DEMO_PATO  -- A la Demo_pato casita  DEMO_PATO ", false},
		//{"another very complex interpolation", args{"{{ .redirect_orange | upper }} {{ .mix | title }} {{ .mix }}", varsContent}, "DEMO_PATO A LA DEMO_PATO CASITA  DEMO_PATO  A La Demo_pato Casita  DEMO_PATO  A la Demo_pato casita  DEMO_PATO ", false},
		//{"connan test", args{"I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero.bad_batman | title }}", values}, "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs With The Using The Conan Sword", false},
		//{"connan test", args{"{{ .colour.orange | upper }} -- {{ .colour.orange | title }} -- {{ .colour.orange }} -- {{ .mix }}", varsContent}, "DEMO_PATO -- Demo_pato -- demo_pato -- A la Demo_pato casita  DEMO_PATO ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do2(tt.args.str, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Do2() = %v, want %v", got, tt.want)
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
		{"Should extract as keys the name of the variables using also functions ", args{`title .insult "hola | | | ma, nsdsds" | tesss`}, []string{".user_name", ".insult"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractKeys(tt.args.str); !reflect.DeepEqual(got, tt.want) {
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
		{"Should return the same string with true flag ", args{"{{ funcion .como  | estoesunafuncion | to_json }}"}, " funcion .como  | estoesunafuncion | to_json ", true},
		{"Should return the same string with true flag ", args{"{{ funcion .como  | estoesunafuncion }}"}, " funcion .como  | estoesunafuncion  | to_json ", false},
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
		{"extract the parameters contained within the string", args{`para empezar esto {{ funcion .como  | estoesunafuncion | to_json }} esto es unba prueba {{ Hola_2 como estas  }}`}, []string{"{{ funcion .como  | estoesunafuncion | to_json }}", "{{ Hola_2 como estas  }}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractParamteres(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamteres() = %v, want %v", got, tt.want)
			}
		})
	}
}
