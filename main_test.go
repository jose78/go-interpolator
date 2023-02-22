package interpolator

import (
	"testing"
)

func Test_fnInterpolateString(t *testing.T) {

	varsContent := map[string]interface{}{
		"mix":       "{{ .house }}",
		"house":     "A {{ .the }} casita",
		"the":       "la {{ .cosa_rara | title  }}",
		"animal":    "de {{ .the | title }} mariposa",
		"cosa_rara": "demo_pato",
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
		{"Prueba", args{"{{ .the }}", varsContent}, "la {{ .cosa_rara | title  }}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interpolateString(tt.args.str, tt.args.vars); got != tt.want {
				t.Errorf("fnInterpolateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		want    string
		wantErr bool
	}{
		{"simple interpolation", args{"{{ .cosa_rara | upper }}", varsContent}, "DEMO_PATO", false},
		//{"medium complex interpolation", args{"{{ .mix }}", varsContent}, "A la Demo_pato casita  DEMO_PATO ", false},
		//{"very complex interpolation", args{"{{ .colour.orange |  upper }}", varsContent}, "DEMO_PATO", false},
		//{"another very complex interpolation", args{"{{ .redirect_pink | upper }}", varsContent}, "ROSA", false},
		//{"another very complex interpolation", args{"{{ .redirect_pink | upper }} {{ .mix | title }} {{ .mix }}", varsContent}, "ROSA A La Demo_pato Casita  DEMO_PATO  A la Demo_pato casita  DEMO_PATO ", false},
		//{"another very complex interpolation", args{"{{ .redirect_orange | upper }} {{ .mix | title }} {{ .mix }}", varsContent}, "DEMO_PATO A LA DEMO_PATO CASITA  DEMO_PATO  A La Demo_pato Casita  DEMO_PATO  A la Demo_pato casita  DEMO_PATO ", false},
		{"connan", args{"I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero.bad_batman | title }}", values}, "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs With The Using The Conan Sword", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.str, tt.args.vars)
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
