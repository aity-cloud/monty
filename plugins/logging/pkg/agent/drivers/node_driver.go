package drivers

import (
	"context"

	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/plugins/logging/apis/node"
)

type LoggingNodeDriver interface {
	ConfigureNode(*node.LoggingCapabilityConfig)
}

var NodeDrivers = driverutil.NewCache[LoggingNodeDriver]()

func NewListenerFunc(ctx context.Context, fn func(cfg *node.LoggingCapabilityConfig)) chan<- *node.LoggingCapabilityConfig {
	listenerC := make(chan *node.LoggingCapabilityConfig, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case cfg := <-listenerC:
				fn(cfg)
			}
		}
	}()
	return listenerC
}
