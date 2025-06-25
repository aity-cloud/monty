package test

import (
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/plugins/logging/pkg/agent"
	"github.com/aity-cloud/monty/plugins/logging/pkg/gateway"
)

func init() {
	test.EnablePlugin(meta.ModeGateway, gateway.Scheme)
	test.EnablePlugin(meta.ModeAgent, agent.Scheme)
}
