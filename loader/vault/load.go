package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type API struct {
	Client *api.Client
}

func (c *API) setConnect() error {
	if c.Client == nil {
		if err := c.Connect(); err != nil {
			return err
		}

		if err := c.Login(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func (c *API) Login(ctx context.Context) error {
	if err := c.setConnect(); err != nil {
		return err
	}

	// A combination of a Role ID and Secret ID is required to log in to Vault with an AppRole.
	// First, let's get the role ID given to us by our Vault administrator.
	roleID := os.Getenv("VAULT_ROLE_ID")
	if roleID == "" {
		return fmt.Errorf("no role ID was provided in APPROLE_ROLE_ID env var")
	}

	secret, err := c.Client.Logical().WriteWithContext(ctx, "auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": os.Getenv("VAULT_ROLE_SECRET"),
	})
	if err != nil {
		return fmt.Errorf("failed to login to vault: %v", err)
	}

	// Set the token
	c.Client.SetToken(secret.Auth.ClientToken)

	return nil
}

func (c *API) Connect() error {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to create vault client: %v", err)
	}

	c.Client = client

	return nil
}

func (c *API) LoadMap(ctx context.Context, mountPath string, key string) (map[string]interface{}, error) {
	if err := c.setConnect(); err != nil {
		return nil, err
	}

	// Get the key
	secret, err := c.Client.KVv2(mountPath).Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %v", err)
	}

	return secret.Data, nil
}

func (c *API) SetMap(ctx context.Context, mountPath string, key string, value map[string]interface{}) error {
	if err := c.setConnect(); err != nil {
		return err
	}

	_, err := c.Client.KVv2(mountPath).Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}

	return nil
}
