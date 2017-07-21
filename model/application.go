package model

type Application struct {
	Name        string
	Namespace   string
	Definitions map[string]string
}
