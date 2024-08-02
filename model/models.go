package model

type FilterLines struct {
	Lines []Line `json:"lines"`
}

type Line struct {
	LineId string `json:"line_id"`
	Title  string `json:"title"`
}
