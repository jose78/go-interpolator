package interpolator

import (
	"reflect"
	"testing"
)

//func TestDo(t *testing.T) {
//
//	a := map[string]interface{}{
//		"name":               "            Jose                 ",
//		"main_topic":         "restore the snyderverse",
//		"favorite_superhero": "batman who laughs",
//		"the":                "the",
//	}
//
//	//b := map[string]interface{}{
//	//	"house":     "A {{ .the }} casita ",
//	//	"colour":    "rosa,",
//	//	"the":       "la {{ .cosa_rara}}",
//	//	"animal":    "de {{ .the }} mariposa",
//	//	"cosa_rara": "demo_pato {{ .the }}",
//	//}
//
//	c := map[string]interface{}{
//		"house":     "A {{ .the | title }} casita ",
//		"colour":    "rosa,",
//		"the":       "la {{ .cosa_rara | title  }}",
//		"animal":    "de {{ .the | title }} mariposa",
//		"cosa_rara": "demo_pato",
//	}
//
//	type args struct {
//		str  string
//		vars map[string]interface{}
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    Content
//		wantErr bool
//	}{
//		{"Should check the content generate by the interpolator is correct",
//			args{"I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", a},
//			content{value: "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs"},
//			false,
//		}, {"Should interpolate all nested templates",
//			args{" {{ .house }} {{ .colour }} {{ .animal }} {{ .cosa_rara }}", c},
//			content{" A La Demo_pato casita  rosa, de La Demo_pato mariposa demo_pato"},
//			false},
//		//{"should fail",
//		//	args{" {{ .house }} {{ .colour }} {{ .the }} {{ .animal }} {{ .cosa_rara}}", b},
//		//	content{" A la demo_pato la demo_pato {{ .the }} casita  rosa, la demo_pato {{ .the }} de la demo_pato {{ .the }} mariposa demo_pato la demo_pato {{ .the }}"},
//		//	true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Do(tt.args.str, tt.args.vars)
//			if tt.wantErr {
//				if (err != nil) != tt.wantErr {
//					t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//			} else {
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("Do() = %v, want %v", got, tt.want)
//				}
//			}
//		})
//	}
//}

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
			got := extractKeys(tt.args.str)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fnExecuteInterpolator(t *testing.T) {

	varsContent := map[string]interface{}{
		"mix":       "{{ .house }}  {{ .cosa_rara | upper  }} ",
		"house":     "A {{ .the }} casita",
		"the":       "la {{ .cosa_rara | title  }}",
		"animal":    "de {{ .the | title }} mariposa",
		"cosa_rara": "demo_pato",
		"mapa":      "demo_pato",
		"cyclic":    "This is a {{ .cyclic }}",
		"colour": map[string]interface{}{
			"red":    "rojo",
			"blue":   "azul",
			"pink":   "rosa",
			"orange": "{{ mapa }}",
		},
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
		{"should interpolate this", args{"{{ .mix }} ", varsContent}, "A la Demo_pato casita  DEMO_PATO  ", false},
		{"should interpolate this", args{"{{ .house }}  {{ .cosa_rara | upper  }} ", varsContent}, "A la Demo_pato casita  DEMO_PATO ", false},
		{"should interpolate this", args{"{{ .house }}", varsContent}, "A la Demo_pato casita", false},
		{"should generate a error of type cyclic", args{"{{ .cyclic }}", varsContent}, "A La Demo_pato casita ", true},
		{"should interpolate a coolor", args{"{{ .colour.pink }}", varsContent}, "rosa", false},
		//{"should interpolate a coolor", args{"{{ .colour.orange }}", varsContent}, "demo_pato", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.str, tt.args.vars)
			if tt.wantErr {
				if err == nil {
					t.Errorf("fnExecuteInterpolator() should return an error and returned: %s", got)
				}
			} else {
				if got != tt.want {
					t.Errorf("fnExecuteInterpolator() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

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
			if got := fnInterpolateString(tt.args.str, tt.args.vars); got != tt.want {
				t.Errorf("fnInterpolateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDo2(t *testing.T) {

	varsContent := map[string]interface{}{
		"mix":       "{{ .house }}  {{ .cosa_rara | upper  }} ",
		"house":     "A {{ .the }} casita",
		"the":       "la {{ .cosa_rara | title  }}",
		"animal":    "de {{ .the | title }} mariposa",
		"cosa_rara": "demo_pato",
		"mapa":      "demo_pato",
		"redirect_pink": "{{ .colour.pink }}",
		"redirect_orange": "{{ .colour.orange }}",
		"cyclic":    "This is a {{ .cyclic }}",
		"colour": map[string]interface{}{
			"red":    "rojo",
			"blue":   "azul",
			"pink":   "rosa",
			"orange": "{{ .mapa }}",
		},
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
		//{"simple interpolation", args{"{{ .cosa_rara | upper }}", varsContent}, "DEMO_PATO", false},
		//{"medium complex interpolation", args{"{{ .mix }}", varsContent}, "A la Demo_pato casita  DEMO_PATO ", false},
		{"very complex interpolation", args{"{{ .colour.orange |  upper }}", varsContent}, "demo_pato", false},
		//{"another very complex interpolation", args{"{{ .redirect_pink | upper }}", varsContent}, "ROSA", false},
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
