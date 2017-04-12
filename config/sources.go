package config

type VariableSource string

const (
	SIMPLE  VariableSource = "simple"
	POD                    = "pod"
	VAULT                  = "valut"
	FILE                   = "file"
	PODYAML                = "pod_yaml"
)
