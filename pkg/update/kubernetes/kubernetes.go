package kubernetes

import (
	"github.com/aity-cloud/monty/pkg/oci"
)

const (
	UpdateStrategy = "kubernetes"
)

var (
	ComponentImageMap = map[ComponentType]oci.ImageType{
		AgentComponent:      oci.ImageTypeMinimal,
		ControllerComponent: oci.ImageTypeMonty,
	}
)

type ComponentType string

const (
	AgentComponent      ComponentType = "agent"
	ControllerComponent ComponentType = "controller"
)
