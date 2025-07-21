package models

import (
	"fmt"
	"strings"
	"time"
)

// Student 학생 정보를 나타냅니다
// @Description 학생 정보 모델
type Student struct {
	StudentID   int    `json:"student_id" example:"1001" minimum:"1000" maximum:"99999"`
	Name        string `json:"name" example:"홍길동" maxLength:"20"`
	Grade       string `json:"grade" example:"초등 3학년" maxLength:"10"`
	Phone       string `json:"phone" example:"010-1234-5678" maxLength:"15"`
	ParentPhone string `json:"parent_phone" example:"010-9876-5432" maxLength:"15"`
}

// Attendance 출결 정보를 나타냅니다
// @Description 출결 정보 모델
type Attendance struct {
	StudentID int      `json:"student_id" example:"1001"`
	Date      DateOnly `json:"date" swaggertype:"string" example:"2024-01-15"`
	CheckIn   TimeOnly `json:"check_in" swaggertype:"string" example:"13:55"`
	CheckOut  TimeOnly `json:"check_out" swaggertype:"string" example:"15:35"`
	Status    string   `json:"status" example:"출석" enums:"출석,결석,지각"`
}

// Teacher 교사 정보를 나타냅니다
// @Description 교사 정보 모델
type Teacher struct {
	TeacherID string `json:"teacher_id" example:"T001" maxLength:"30"`
	Password  string `json:"password" example:"hashed_password" maxLength:"100"`
	Name      string `json:"name" example:"김선생님" maxLength:"30"`
	Phone     string `json:"phone" example:"010-1111-2222" maxLength:"20"`
}

// TeacherResponse 교사 응답 정보를 나타냅니다
// @Description 교사 응답 정보 모델
type TeacherResponse struct {
	TeacherID string `json:"teacher_id" example:"T001"`
	Name      string `json:"name" example:"김선생님"`
	Phone     string `json:"phone" example:"010-1111-2222"`
}

// Class 수업 정보를 나타냅니다
// @Description 수업 정보 모델
type Class struct {
	ClassID   int      `json:"class_id" example:"1"`
	ClassName string   `json:"class_name" example:"수학 기초반" maxLength:"50"`
	Days      string   `json:"days" example:"월,수,금" maxLength:"20"`
	StartTime TimeOnly `json:"start_time" swaggertype:"string" example:"14:00"`
	EndTime   TimeOnly `json:"end_time" swaggertype:"string" example:"15:30"`
	Price     int      `json:"price" example:"150000" minimum:"0"`
	TeacherID string   `json:"teacher_id" example:"T001" maxLength:"30"`
}

// Enrollment 수강신청 정보를 나타냅니다
// @Description 수강신청 정보 모델
type Enrollment struct {
	EnrollmentID int      `json:"enrollment_id" example:"1"`
	StudentID    int      `json:"student_id" example:"1001"`
	ClassID      int      `json:"class_id" example:"1"`
	EnrolledDate DateOnly `json:"enrolled_date" swaggertype:"string" example:"2024-01-15"`
}

// Payment 결제 정보를 나타냅니다
// @Description 결제 정보 모델
type Payment struct {
	PaymentID    int      `json:"payment_id" example:"1"`
	StudentID    int      `json:"student_id" example:"1001"`
	ClassID      int      `json:"class_id" example:"1"`
	PaymentDate  DateOnly `json:"payment_date" swaggertype:"string" example:"2024-01-15"`
	Amount       int      `json:"amount" example:"150000" minimum:"0"`
	EnrollmentID int      `json:"enrollment_id" example:"1"`
}

// TimeOnly 시간만을 나타내는 타입입니다
// @Description 시간만을 나타내는 타입
type TimeOnly struct {
	time.Time
}

// UnmarshalJSON JSON에서 시간을 파싱합니다
func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// MarshalJSON 시간을 JSON으로 직렬화합니다
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", t.Format("15:04")), nil
}

// DateOnly 날짜만을 나타내는 타입입니다
// @Description 날짜만을 나타내는 타입
type DateOnly struct {
	time.Time
}

// UnmarshalJSON JSON에서 날짜를 파싱합니다
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

// MarshalJSON 날짜를 JSON으로 직렬화합니다
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", d.Format("2006-01-02")), nil
}

// ToResponse Teacher를 TeacherResponse로 변환합니다
func (t Teacher) ToResponse() TeacherResponse {
	return TeacherResponse{
		TeacherID: t.TeacherID,
		Name:      t.Name,
		Phone:     t.Phone,
	}
}
