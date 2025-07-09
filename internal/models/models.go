package models

import "time"

type Student struct {
	StudentID   int    `json:"student_id"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	Phone       string `json:"phone"`
	ParentPhone string `json:"parent_phone"`
}

type Attendance struct {
	StudentID int       `json:"student_id"`
	Date      time.Time `json:"date"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	Status    string    `json:"status"`
}
