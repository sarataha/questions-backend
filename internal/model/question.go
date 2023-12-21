package model

type Question struct {
	ID      int      `json:"id"`
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}
