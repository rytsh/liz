package vault

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

func TestAPI_LoadMap(t *testing.T) {
	type fields struct {
		Client  *api.Client
		Approle *approle.SecretID
	}
	type args struct {
		ctx       context.Context
		mountPath string
		key       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "TestAPI_LoadMap",
			args: args{
				ctx:       context.Background(),
				mountPath: "secret",
				key:       "test",
			},
			want: map[string]interface{}{
				"test": "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &API{
				Client: tt.fields.Client,
			}
			got, err := c.LoadMap(tt.args.ctx, tt.args.mountPath, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.LoadMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if deep.Equal(got, tt.want) != nil {
				t.Errorf("API.LoadMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
