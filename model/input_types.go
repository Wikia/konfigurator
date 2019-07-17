package model

type InputType string

const (
	SIMPLE         InputType = "simple"
	REFERENCE      InputType = "reference"
	VAULT          InputType = "vault"
	CONSUL         InputType = "consul"
	LAYERED_CONSUL InputType = "layered_consul"
)
