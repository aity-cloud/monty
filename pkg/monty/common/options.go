package common

// These constants are available to all montyctl sub-commands and are filled
// in by the root command using persistent flags.

var (
	NamespaceFlagValue string
	DisableUsage       bool
)

const (
	DefaultMontyNamespace = "monty-system"
)
