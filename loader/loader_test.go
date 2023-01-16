package loader

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestConfigs_Load(t *testing.T) {
	tempDir := t.TempDir()

	type want struct {
		content string
		path    string
	}

	tests := []struct {
		name     string
		c        Configs
		skip     bool
		wantErr  bool
		want     want
		wantWait time.Duration
	}{
		{
			name: "test",
			c: Configs{
				{
					Export: tempDir + "/test",
					Statics: []ConfigStatic{
						{
							Content: &ConfigContent{
								Content: "test",
								Raw:     true,
							},
						},
					},
				},
			},
			want: want{
				content: "test",
				path:    tempDir + "/test",
			},
		},
		{
			name: "test with dynamic",
			c: Configs{
				{
					Export: tempDir + "/testd.yml",
					Statics: []ConfigStatic{
						{
							Content: &ConfigContent{
								Content: `test: "mycontent"`,
							},
							Vault: &ConfigVault{
								Path:       "test",
								PathPrefix: "secret",
								AdditionalPaths: []ConfigVaultAdditional{
									{
										Map:  "myValue2/in",
										Path: "inner/x",
									},
								},
							},
						},
					},
					Dynamics: []ConfigDynamic{
						{
							Consul: &ConfigConsul{
								Path: "test",
							},
						},
					},
				},
			},
			skip:     true,
			wantWait: 1 * time.Second,
			want: want{
				content: `myValue: "1234"` + "\n" + `test: mycontent` + "\n",
				path:    tempDir + "/testd.yml",
			},
		},
	}

	// set consul address
	os.Setenv("CONSUL_HTTP_ADDR", "http://localhost:8500")
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Setenv("VAULT_ROLE_ID", "88eda05e-b98e-dda4-7251-e97a0638adc9")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip && testing.Short() {
				t.Skip()
			}

			if err := tt.c.Load(context.Background(), nil, nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("Configs.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(tt.wantWait)

			// check export file content
			v, err := os.ReadFile(tt.want.path)
			if err != nil {
				t.Errorf("Configs.Load() error = %v", err)
			}

			if string(v) != tt.want.content {
				t.Errorf("Configs.Load() content = \n%q\n, want \n%q\n", string(v), tt.want.content)
			}
		})
	}
}
