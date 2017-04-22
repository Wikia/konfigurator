package model

type InputType string

const (
	SIMPLE  InputType = "simple"
	POD               = "pod"
	VAULT             = "vault"
	FILE              = "file"
	PODYAML           = "pod_yaml"
)
