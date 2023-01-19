package models

type Transaction struct {
	Id        string `json:"id"`
	Candidate string `json:"candidate"`
	Voter     string `json:"voter"`
	Timestamp int64  `json:"timestamp"`
}
