package shared

// Jetstream streams
const (
	// global streams
	AgentClusterHealthStatusStream        = "agent-cluster-health-status"
	AgentClusterHealthStatusSubjects      = "agent-cluster-health-status.*"
	AgentClusterHealthStatusDurableReplay = "agent-cluster-health-status-consumer"

	//streams
	AgentDisconnectStream         = "opni_alerting_agent"
	AgentDisconnectStreamSubjects = "opni_alerting_agent.*"
	AgentHealthStream             = "opni_alerting_health"
	AgentHealthStreamSubjects     = "opni_alerting_health.*"
	CortexStatusStream            = "opni_alerting_cortex_status"
	CortexStatusStreamSubjects    = "opni_alerting_cortex_status.*"
	// buckets
	AlertingConditionBucket            = "monty-alerting-condition-bucket"
	AlertingEndpointBucket             = "monty-alerting-endpoint-bucket"
	AgentDisconnectBucket              = "monty-alerting-agent-bucket"
	AgentStatusBucket                  = "monty-alerting-agent-status-bucket"
	StatusBucketPerCondition           = "monty-alerting-condition-status-bucket"
	StatusBucketPerClusterInternalType = "monty-alerting-cluster-condition-type-status-bucket"
	GeneralIncidentStorage             = "monty-alerting-general-incident-bucket"
	RouterStorage                      = "monty-alerting-router-bucket"
)
