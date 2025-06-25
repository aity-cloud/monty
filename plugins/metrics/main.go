package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/agent"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway"

	_ "github.com/aity-cloud/monty/plugins/metrics/pkg/agent/drivers/opni_manager_otel"
	_ "github.com/aity-cloud/monty/plugins/metrics/pkg/agent/drivers/prometheus_operator"
	_ "github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/drivers/opni_manager"
)

func main() {
	m := plugins.Main{
		Modes: meta.ModeSet{
			meta.ModeGateway: gateway.Scheme,
			meta.ModeAgent:   agent.Scheme,
		},
	}
	m.Exec()
}
