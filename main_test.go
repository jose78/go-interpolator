package judo_interpolator

import (
	"reflect"
	"testing"
)

func TestDo(t *testing.T) {
	type args struct {
		str  string
		vars map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want Content
	}{
		{name: "Should check the content generate by the judo_interpolator is correct",
			args: args{str: "I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}",
				vars: map[string]interface{}{
					"name":               "            Jose                 ",
					"main_topic":         "restore the snyderverse",
					"favorite_superhero": "batman who laughs",
				}},
			want: content{value: "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs"},
		},
		{name: "Should check the content generate by the judo_interpolator is correct using varsContent",
		args: args{str: "I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}",
			vars: VarsContent{
				"name":               "            Jose                 ",
				"main_topic":         "restore the snyderverse",
				"favorite_superhero": "batman who laughs",
			}},
		want: content{value: "I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs"},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Do(tt.args.str, tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
