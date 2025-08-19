package shared

// Jetstream streams
const (
	// global streams
	AgentClusterHealthStatusStream        = "agent-cluster-health-status"
	AgentClusterHealthStatusSubjects      = "agent-cluster-health-status.*"
	AgentClusterHealthStatusDurableReplay = "agent-cluster-health-status-consumer"

	//streams
	AgentDisconnectStream         = "monty_alerting_agent"
	AgentDisconnectStreamSubjects = "monty_alerting_agent.*"
	AgentHealthStream             = "monty_alerting_health"
	AgentHealthStreamSubjects     = "monty_alerting_health.*"
	CortexStatusStream            = "monty_alerting_cortex_status"
	CortexStatusStreamSubjects    = "monty_alerting_cortex_status.*"
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
