package model

type Prop struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Model struct {
	Name    string   `json:"name" orm:",NOT_NULL|PRIMARY_KEY"`
	Props   []Prop   `json:"props"`
	Methods []string `json:"methods"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
}

const MODELS_NAME = "models"
