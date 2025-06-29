package main

import (
	//mage:import
	"github.com/aity-cloud/monty/magefiles/targets"
)

var Default = targets.Default

var Aliases = map[string]any{
	"test":     targets.Test.All,
	"build":    targets.Build.All,
	"generate": targets.Generate.All,
	"crd":      targets.CRD.All,
}
