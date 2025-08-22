package test

import (
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/plugins/example/pkg/example"
)

func init() {
	test.EnablePlugin(meta.ModeGateway, example.Scheme)
}
