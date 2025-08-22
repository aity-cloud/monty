package resources

const (
	PretrainedModelLabel = "monty.io/pretrained-model"
	ServiceLabel         = "monty.io/service"
	MontyClusterName     = "monty.io/cluster-name"
	AppNameLabel         = "app.kubernetes.io/name"
	PartOfLabel          = "app.kubernetes.io/part-of"
	InstanceLabel        = "app.kubernetes.io/instance"
	HostTopologyKey      = "kubernetes.io/hostname"
	MontyClusterID       = "monty.io/cluster-id"
	MontyBootstrapToken  = "monty.io/bootstrap-token"
	MontyInferenceType   = "monty.io/inference-type"
	MontyConfigHash      = "monty.io/config-hash"
)

func CombineLabels(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

type OpensearchLabels map[string]string

func NewOpensearchLabels() OpensearchLabels {
	return map[string]string{
		"app": "opensearch",
	}
}

func NewGatewayLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name": "monty-gateway",
	}
}
