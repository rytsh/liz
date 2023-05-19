package consul

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
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

// DynamicValue return a channel for getting latest value of key.
// This function will start a goroutine for watching key.
// The caller should call stop function when it is no longer needed.
func (c *API) DynamicValue(ctx context.Context, wg *sync.WaitGroup, key string) (<-chan []byte, func(), error) {
	if err := c.setConnect(); err != nil {
		return nil, nil, err
	}

	plan, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  key,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("wath.Parse %w", err)
	}

	// not add any buffer, this is useful for getting latest change only
	vChannel := make(chan []byte)

	plan.HybridHandler = func(_ watch.BlockingParamVal, raw interface{}) {
		if raw == nil {
			return
		}

		v, ok := raw.(*api.KVPair)
		if ok {
			vChannel <- v.Value
			return
		}
	}

	runCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// this select-case for listen ctx done and plan run result same time
		select {
		case <-ctx.Done():
			plan.Stop()
		case <-runCh:
		}

		close(vChannel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runCh <- plan.RunWithClientAndHclog(c.Client, hclog.NewNullLogger())
	}()

	return vChannel, plan.Stop, nil
}
