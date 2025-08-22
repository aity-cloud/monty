package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/aiops/pkg/gateway"
)

func main() {
	m := plugins.Main{
		Modes: meta.ModeSet{
			meta.ModeGateway: gateway.Scheme,
		},
	}
	m.Exec()
}
