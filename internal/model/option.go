package model

type Option struct {
	ID      int    `json:"id"`
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}
