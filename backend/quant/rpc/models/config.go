package models

type Classify struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

type Item struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
