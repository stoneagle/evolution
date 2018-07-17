package config

type ProjectType int8

const (
	ProjectQuant ProjectType = iota + 1
	ProjectTime
	ProjectSystem
)
