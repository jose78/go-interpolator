package interpolator

import (
	"reflect"
	"testing"
)

func TestInterpolate(t *testing.T) {

	varsContent := map[string]interface{}{
		"mix":             "{{ house }}  {{ cosa_rara | upper  }} ",
		"house":           "A {{ the }} casita",
		"the":             "la {{ cosa_rara | title  }}",
		"animal":          "de {{ the | title }} mariposa",
		"cosa_rara":       "demo_pato",
		"mapa":            "demo_pato",
		"redirect_pink":   "{{ colour.pink }}",
		"redirect_orange": "{{ colour.orange }} {{ mix }}",
		"cyclic":          "This is a {{ cyclic }}",
		"lst":             []string{"uno", "dos", "tres"},
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
		name       string
		args       args
		wantResult interface{}
		wantErr    bool
	}{
		//{"test", args{"colour.red", varsContent}, "rojo", false},
		//{"test", args{"lst", varsContent}, []string{"uno", "dos", "tres"}, false},
		{"test", args{"{{colour.orange}} {{colour.orange}} ", varsContent}, "demo_pato demo_pato ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Interpolate(tt.args.str, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interpolate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Interpolate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
