package model

type Relation struct {
	Model1  int64 `json:"model1"`
	Model2  int64 `json:"model2"`
	Model1n int   `json:"model1n"`
	Model2n int   `json:"model2n"`
}

const RELATIONS_NAME = "relations"
