package mapx

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		key []string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name: "simple one test",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				key: []string{"def", "abc"},
			},
			want:  1,
			want1: true,
		},
		{
			name: "over map",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				key: []string{"def", "abc", "xyz"},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "get map",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				key: []string{"def"},
			},
			want:  map[string]interface{}{"abc": 1, "xyz": 2},
			want1: true,
		},
		{
			name: "get map",
			args: args{
				m: map[string]interface{}{
					"abc": []int{1, 2, 3},
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				key: []string{"abc", "1"},
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Get(tt.args.m, tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
