package test

import (
	"time"

	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/test"
	sloapi "github.com/aity-cloud/monty/plugins/slo/apis/slo"
	"github.com/aity-cloud/monty/plugins/slo/pkg/slo"
)

func init() {
	sloapi.MinEvaluateInterval = time.Second
	test.EnablePlugin(meta.ModeGateway, slo.Scheme)
}
