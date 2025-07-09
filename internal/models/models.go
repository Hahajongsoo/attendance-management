package models

type Student struct {
	StudentID   int    `json:"student_id"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	Phone       string `json:"phone"`
	ParentPhone string `json:"parent_phone"`
}
