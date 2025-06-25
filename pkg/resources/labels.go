package resources

const (
	PretrainedModelLabel = "monty.io/pretrained-model"
	ServiceLabel         = "monty.io/service"
	OpniClusterName      = "monty.io/cluster-name"
	AppNameLabel         = "app.kubernetes.io/name"
	PartOfLabel          = "app.kubernetes.io/part-of"
	InstanceLabel        = "app.kubernetes.io/instance"
	HostTopologyKey      = "kubernetes.io/hostname"
	OpniClusterID        = "monty.io/cluster-id"
	OpniBootstrapToken   = "monty.io/bootstrap-token"
	OpniInferenceType    = "monty.io/inference-type"
	OpniConfigHash       = "monty.io/config-hash"
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
