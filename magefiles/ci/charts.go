//go:build mage

package main

import (
	// mage:import charts
	charts "github.com/rancher/charts-build-scripts/pkg/actions"

	// mage:import targets
	"github.com/aity-cloud/monty/magefiles/targets"

	"github.com/magefile/mage/mg"
)

func Charts() {
	mg.SerialDeps(
		targets.CRD.All,
		func() {
			charts.Charts("monty")
		}, func() {
			charts.Charts("monty-agent")
		})
}
