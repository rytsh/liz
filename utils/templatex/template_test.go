package templatex

import (
	"testing"

	"github.com/rytsh/liz/utils/templatex/functions"
)

func TestTemplate_Execute(t *testing.T) {
	type args struct {
		v       any
		content string
	}
	tests := []struct {
		name    string
		args    args
		opts    []functions.Option
		want    string
		wantErr bool
	}{
		{
			name: "test template",
			args: args{
				v:       map[string]interface{}{"name": "test"},
				content: `{{ .name }}`,
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "Os",
			args: args{
				v:       map[string]interface{}{"fileName": "../templatex"},
				content: `{{ fileExists .fileName }}`,
			},
			want:    "true",
			wantErr: false,
		},
		{
			name: "Os workdir",
			args: args{
				v:       map[string]interface{}{"fileName": "custom.go"},
				content: `{{ fileExists .fileName }}`,
			},
			opts:    []functions.Option{functions.WorkDir("./functions/custom")},
			want:    "true",
			wantErr: false,
		},
		{
			name: "readDir",
			args: args{
				v:       map[string]interface{}{"dir": "."},
				content: `{{ range readDir .dir }}{{ .Name }}{{ end}}`,
			},
			opts:    []functions.Option{functions.WorkDir("./functions/custom")},
			want:    "custom.go",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateEngine := New(tt.opts...)

			got, err := templateEngine.Execute(tt.args.v, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Template.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
