package models

type Suggestion struct {
	Source  string  `json:"source"`
	Account Account `json:"account"`
}
