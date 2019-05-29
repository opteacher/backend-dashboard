package model

type Relation struct {
	BegModel int64 `json:"selBegMdl"`
	EndModel int64 `json:"selEndMdl"`
	BegModelNum int `json:"begMdlNum"`
	EndModelNum int `json:"endMdlNum"`
}

const RELATIONS_NAME = "relations"