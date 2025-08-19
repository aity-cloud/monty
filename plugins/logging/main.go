package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/logging/pkg/agent"
	"github.com/aity-cloud/monty/plugins/logging/pkg/gateway"

	_ "github.com/aity-cloud/monty/plugins/logging/pkg/agent/drivers/kubernetes_manager"
	_ "github.com/aity-cloud/monty/plugins/logging/pkg/gateway/drivers/backend/kubernetes_manager"
	_ "github.com/aity-cloud/monty/plugins/logging/pkg/gateway/drivers/management/kubernetes_manager"
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
