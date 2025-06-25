package types

import (
	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	"github.com/aity-cloud/monty/pkg/metrics/collector"
	"github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions"
	"github.com/aity-cloud/monty/pkg/plugins/apis/system"
)

type (
	HTTPAPIExtensionPlugin       = apiextensions.HTTPAPIExtensionClient
	ManagementAPIExtensionPlugin = apiextensions.ManagementAPIExtensionClient
	StreamAPIExtensionPlugin     = apiextensions.StreamAPIExtensionClient
	UnaryAPIExtensionPlugin      = apiextensions.UnaryAPIExtensionClient
	CapabilityBackendPlugin      = capabilityv1.BackendClient
	CapabilityNodePlugin         = capabilityv1.NodeClient
	MetricsPlugin                = collector.RemoteProducer
	SystemPlugin                 = system.SystemPluginServer
)
