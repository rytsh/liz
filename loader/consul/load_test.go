package consul

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/consul/api"
)

func TestConsul_LoadRaw(t *testing.T) {
	type fields struct {
		Client *api.Client
	}
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test LoadRaw",
			args: args{
				key:   "test12345-67890",
				value: []byte("test"),
			},
			want:    []byte("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consul{
				Client: tt.fields.Client,
			}

			ctx := context.Background()

			if err := c.SetRaw(ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Consul.SetRaw() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := c.LoadRaw(ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Consul.LoadRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if deep.Equal(got, tt.want) != nil {
				t.Errorf("Consul.LoadRaw() = %v, want %v", got, tt.want)
			}

			// cleanup
			if err := c.Delete(ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Consul.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
