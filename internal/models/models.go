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

type TeacherResponse struct {
	TeacherID string `json:"teacher_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type Class struct {
	ClassID   int      `json:"class_id"`
	ClassName string   `json:"class_name"`
	Days      string   `json:"days"`
	StartTime TimeOnly `json:"start_time"`
	EndTime   TimeOnly `json:"end_time"`
	Price     int      `json:"price"`
	TeacherID string   `json:"teacher_id"`
}

type Enrollment struct {
	EnrollmentID int      `json:"enrollment_id"`
	StudentID    int      `json:"student_id"`
	ClassID      int      `json:"class_id"`
	EnrolledDate DateOnly `json:"enrolled_date"`
}

type Payment struct {
	PaymentID    int      `json:"payment_id"`
	StudentID    int      `json:"student_id"`
	ClassID      int      `json:"class_id"`
	PaymentDate  DateOnly `json:"payment_date"`
	Amount       int      `json:"amount"`
	EnrollmentID int      `json:"enrollment_id"`
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

func (t Teacher) ToResponse() TeacherResponse {
	return TeacherResponse{
		TeacherID: t.TeacherID,
		Name:      t.Name,
		Phone:     t.Phone,
	}
}
