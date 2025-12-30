package model

type Task struct {
	ID          int64
	Header      string
	Description string
	Completed   bool
}
