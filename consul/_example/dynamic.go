package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rytsh/liz/consul"
)

func main() {
	c := consul.API{}

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	defer wg.Wait()
	defer ctxCancel()

	var stop func()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// listen Ctrl+C
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		select {
		case <-ctx.Done():
		case <-signalCh:
			if stop != nil {
				stop()
			}
		}
	}()

	ch, stop, err := c.DynamicValue(ctx, wg, "test")
	if err != nil {
		panic(err)
	}

	time.AfterFunc(15*time.Second, func() {
		stop()
	})

	for v := range ch {
		log.Printf("value: %s", v)
	}
}
