package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/agent"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/alerting"

	_ "github.com/aity-cloud/monty/plugins/alerting/pkg/agent/drivers/default_driver"
	_ "github.com/aity-cloud/monty/plugins/alerting/pkg/alerting/drivers/alerting_manager"
)

func main() {
	m := plugins.Main{
		Modes: meta.ModeSet{
			meta.ModeGateway: alerting.Scheme,
			meta.ModeAgent:   agent.Scheme,
		},
	}
	m.Exec()
}
