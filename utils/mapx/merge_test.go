package mapx

import (
	"testing"

	"github.com/go-test/deep"
)

func TestMerge(t *testing.T) {
	type args struct {
		value map[string]interface{}
		to    map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "merge",
			args: args{
				value: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"x": "bar",
					},
				},
				to: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": "bar",
					"x":   "bar",
				},
			},
		},
		{
			name: "merge mix",
			args: args{
				value: map[string]interface{}{
					"foo": []interface{}{"bar"},
					"bar": map[string]interface{}{
						"x": map[string]interface{}{
							"foo": "bar",
						},
					},
				},
				to: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"x": map[string]interface{}{
							"foo": []string{"bar"},
						},
					},
				},
			},
			want: map[string]interface{}{
				"foo": []interface{}{"bar"},
				"bar": map[string]interface{}{
					"x": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args.value, tt.args.to)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Merge() = %v", diff)
			}
		})
	}
}
