package main

import (
	"github.com/aity-cloud/monty/pkg/plugins"
	"github.com/aity-cloud/monty/pkg/plugins/meta"
	_ "github.com/aity-cloud/monty/pkg/storage/etcd"
	_ "github.com/aity-cloud/monty/pkg/storage/jetstream"
	"github.com/aity-cloud/monty/plugins/example/pkg/example"
)

func main() {
	m := plugins.Main{
		Modes: meta.ModeSet{
			meta.ModeGateway: example.Scheme,
			meta.ModeAgent:   example.Scheme,
		},
	}
	m.Exec()
}
