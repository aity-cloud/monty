package test

import (
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/agent"
	"github.com/aity-cloud/monty/plugins/metrics/pkg/gateway"

	_ "github.com/aity-cloud/monty/plugins/metrics/pkg/gateway/services"
)

func init() {
	test.EnablePlugin(meta.ModeGateway, gateway.Scheme)
	test.EnablePlugin(meta.ModeAgent, agent.Scheme)
}
