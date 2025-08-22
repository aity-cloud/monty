package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/topology/pkg/topology/agent"
	"github.com/aity-cloud/monty/plugins/topology/pkg/topology/gateway"
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
