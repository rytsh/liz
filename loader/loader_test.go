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
			name: "test vault",
			c: Configs{
				{
					Export: tempDir + "/test",
					Statics: []ConfigStatic{
						{
							Vault: &ConfigVault{
								Path:       "test",
								PathPrefix: "secret",
								Template:   true,
							},
						},
					},
				},
			},
			skip: true,
			want: want{
				content: `test: "1234"` + "\n",
				path:    tempDir + "/test",
			},
		},
		{
			name: "test consul",
			c: Configs{
				{
					Export: tempDir + "/test",
					Statics: []ConfigStatic{
						{
							Consul: &ConfigConsul{
								Path:       "test",
								PathPrefix: "secret",
								Template:   true,
								Raw:        true,
							},
						},
					},
				},
			},
			skip: true,
			want: want{
				content: `test: 1234`,
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
		{
			name: "test inner path",
			c: Configs{
				{
					Export: tempDir + "/test.json",
					Statics: []ConfigStatic{
						{
							Content: &ConfigContent{
								Content:   `{"test": {"test-2": "mycontent"}}`,
								InnerPath: "test",
							},
						},
					},
				},
			},
			want: want{
				content: `{
  "test-2": "mycontent"
}
`,
				path: tempDir + "/test.json",
			},
		},
		{
			name: "test inner path raw",
			c: Configs{
				{
					Export: tempDir + "/test",
					Statics: []ConfigStatic{
						{
							Content: &ConfigContent{
								Content:   `{"test": {"test-2": "mycontent"}}`,
								InnerPath: "test/test-2",
							},
						},
					},
				},
			},
			want: want{
				content: "mycontent",
				path:    tempDir + "/test",
			},
		},
		{
			name: "test inner path raw with base64",
			c: Configs{
				{
					Export: tempDir + "/test",
					Statics: []ConfigStatic{
						{
							Content: &ConfigContent{
								Content:   `{"test": {"test-2": "bWVyaGFiYQ=="}}`,
								InnerPath: "test/test-2",
								Base64:    true,
							},
						},
					},
				},
			},
			want: want{
				content: "merhaba",
				path:    tempDir + "/test",
			},
		},
	}

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
