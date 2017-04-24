package model

type InputType string

const (
	SIMPLE  InputType = "simple"
	POD               = "pod"
	VAULT             = "vault"
	CONSUL            = "consul"
	FILE              = "file"
	PODYAML           = "pod_yaml"
)
