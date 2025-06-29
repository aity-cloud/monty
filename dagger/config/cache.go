package config

import (
	"os"

	"dagger.io/dagger"
)

const (
	CacheModeVolumes = "volumes"
	CacheModeNone    = "none"
)

type CacheVolume struct {
	*dagger.CacheVolume
	Path string
}

type Caches struct {
	GoMod        func(*dagger.Container) *dagger.Container
	GoBuild      func(*dagger.Container) *dagger.Container
	GoBin        func(*dagger.Container) *dagger.Container
	Mage         func(*dagger.Container) *dagger.Container
	Yarn         func(*dagger.Container) *dagger.Container
	NodeModules  func(*dagger.Container) *dagger.Container
	GolangciLint func(*dagger.Container) *dagger.Container
}

func SetupCaches(client *dagger.Client, cacheMode string) Caches {
	if _, ok := os.LookupEnv("CI"); ok {
		cacheMode = CacheModeVolumes
	}
	identity := func(ctr *dagger.Container) *dagger.Container { return ctr }
	switch cacheMode {
	case CacheModeVolumes:
		return Caches{
			GoMod: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/go/pkg/mod", client.CacheVolume("monty_gomod"))
			},
			GoBuild: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/root/.cache/go-build", client.CacheVolume("monty_gobuild"))
			},
			GoBin: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/go/bin", client.CacheVolume("monty_gobin"))
			},
			Mage: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/root/.magefile", client.CacheVolume("monty_mage"))
			},
			Yarn: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/cache/yarn", client.CacheVolume("monty_yarn"))
			},
			NodeModules: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/cache/node_modules", client.CacheVolume("monty_node_modules"))
			},
			GolangciLint: func(ctr *dagger.Container) *dagger.Container {
				return ctr.WithMountedCache("/root/.cache/golangci-lint", client.CacheVolume("monty_golangci_lint"))
			},
		}
	case CacheModeNone:
		fallthrough
	default:
		return Caches{
			GoMod:        identity,
			GoBuild:      identity,
			GoBin:        identity,
			Mage:         identity,
			Yarn:         identity,
			NodeModules:  identity,
			GolangciLint: identity,
		}
	}
}
