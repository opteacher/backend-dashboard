package model

type Prop struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type Model struct {
	Name    string   `json:"name"`
	Props   []Prop   `json:"props"`
	Methods []string `json:"methods"`
}

const MODELS_NAME = "models"
