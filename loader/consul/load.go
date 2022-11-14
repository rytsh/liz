package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type API struct {
	Client       *api.Client
	KV           *api.KV
	QueryOptions api.QueryOptions
	WriteOptions api.WriteOptions
}

func (c *API) setConnect() error {
	if c.KV == nil && c.Client == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	if c.KV == nil {
		c.KV = c.Client.KV()
	}

	return nil
}

func (c *API) Connect() error {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to create consul client: %v", err)
	}

	c.Client = client
	c.KV = client.KV()

	return nil
}

func (c *API) LoadRaw(ctx context.Context, key string) ([]byte, error) {
	if err := c.setConnect(); err != nil {
		return nil, err
	}

	// Get the key
	pair, _, err := c.KV.Get(key, c.QueryOptions.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %v", err)
	}

	return pair.Value, nil
}

func (c *API) SetRaw(ctx context.Context, key string, value []byte) error {
	if err := c.setConnect(); err != nil {
		return err
	}

	// Set the key
	pair := &api.KVPair{Key: key, Value: value}

	_, err := c.KV.Put(pair, c.WriteOptions.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}

	return nil
}

func (c *API) Delete(ctx context.Context, key string) error {
	if err := c.setConnect(); err != nil {
		return err
	}

	// Delete the key
	_, err := c.KV.Delete(key, c.WriteOptions.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete key: %v", err)
	}

	return nil
}
