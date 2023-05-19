package mapx

import (
	"testing"
)

func TestIsMapSubset(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		sub map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple one test",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
				},
				sub: map[string]interface{}{
					"abc": 1,
				},
			},
			want: true,
		},
		{
			name: "mix type",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				sub: map[string]interface{}{
					"abc": 1,
					"def": map[string]interface{}{
						"abc": 1,
					},
				},
			},
			want: true,
		},
		{
			name: "mix type",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": 1,
						"xyz": 2,
					},
				},
				sub: map[string]interface{}{
					"abc": 1,
					"def": map[string]interface{}{
						"abc": "sfdfds",
					},
				},
			},
			want: false,
		},
		{
			name: "mix type in sub different type",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": "sfdfds",
						"xyz": 2,
					},
				},
				sub: map[string]interface{}{
					"abc": 1,
					"def": map[string]interface{}{
						"abc": map[string]interface{}{
							"abc": 1,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "mix type in sub",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": []interface{}{
							"abc",
							"xyz",
							"xyz2",
						},
					},
				},
				sub: map[string]interface{}{
					"abc": 1,
					"def": map[string]interface{}{
						"abc": []interface{}{
							"abc",
							"xyz2",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "mix type in sub false",
			args: args{
				m: map[string]interface{}{
					"abc": 1,
					"xyz": 2,
					"def": map[string]interface{}{
						"abc": []interface{}{
							"abc",
							"xyz",
							"xyz2",
						},
					},
				},
				sub: map[string]interface{}{
					"abc": 1,
					"def": map[string]interface{}{
						"abc": []interface{}{
							"abc",
							"xyz3",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMapSubset(tt.args.m, tt.args.sub); (got == nil) != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
