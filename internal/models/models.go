package models

import (
	"fmt"
	"strings"
	"time"
)

type Student struct {
	StudentID   int    `json:"student_id"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	Phone       string `json:"phone"`
	ParentPhone string `json:"parent_phone"`
}

type Attendance struct {
	StudentID int      `json:"student_id"`
	Date      DateOnly `json:"date"`
	CheckIn   TimeOnly `json:"check_in"`
	CheckOut  TimeOnly `json:"check_out"`
	Status    string   `json:"status"`
}

type Teacher struct {
	TeacherID string `json:"teacher_id"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type TimeOnly struct {
	time.Time
}

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", t.Format("15:04")), nil
}

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", d.Format("2006-01-02")), nil
}
