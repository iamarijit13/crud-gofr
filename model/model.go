package model

type Employee struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	CTC   float32 `json:"ctc"`
}
