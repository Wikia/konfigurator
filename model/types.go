package model

type InputType string

const (
	SIMPLE  InputType = "simple"
	POD               = "pod"
	VAULT             = "valut"
	FILE              = "file"
	PODYAML           = "pod_yaml"
)
