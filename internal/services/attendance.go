package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
	"time"
)

type AttendanceService struct {
	attendanceRepo repositories.AttendanceRepository
	classRepo      repositories.ClassRepository
}

func NewAttendanceService(attendanceRepo repositories.AttendanceRepository, classRepo repositories.ClassRepository) *AttendanceService {
	return &AttendanceService{attendanceRepo: attendanceRepo, classRepo: classRepo}
}

func (s *AttendanceService) GetAttendanceByStudentIDAndDate(studentID string, date string) (*models.Attendance, error) {
	return s.attendanceRepo.GetByStudentIDAndDate(studentID, date)
}

func (s *AttendanceService) GetAttendanceByDate(date string) ([]models.Attendance, error) {
	return s.attendanceRepo.GetByDate(date)
}

func (s *AttendanceService) CreateAttendance(attendance *models.Attendance) error {
	today := time.Now().In(time.Local).Weekday()
	Weekday := convertWeekdayToKorean(today)
	classes, err := s.classRepo.GetClassesForStudentByWeekday(attendance.StudentID, Weekday)
	if err != nil {
		return err
	}
	// 해당 요일에 수업이 없으면 일단 출석으로 처리
	if len(classes) == 0 {
		attendance.Status = "출석"
		return s.attendanceRepo.Create(attendance)
	}

	if attendance.CheckIn.Time.Format("15:04") > classes[0].StartTime.Time.Format("15:04") {
		attendance.Status = "지각"
	} else {
		attendance.Status = "출석"
	}
	return s.attendanceRepo.Create(attendance)
}

func (s *AttendanceService) UpdateAttendance(studentID string, date string, attendance *models.Attendance) error {
	affected, err := s.attendanceRepo.Update(studentID, date, attendance)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *AttendanceService) DeleteAttendance(studentID string, date string) error {
	affected, err := s.attendanceRepo.Delete(studentID, date)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func convertWeekdayToKorean(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "월"
	case time.Tuesday:
		return "화"
	case time.Wednesday:
		return "수"
	case time.Thursday:
		return "목"
	case time.Friday:
		return "금"
	case time.Saturday:
		return "토"
	case time.Sunday:
		return "일"
	}
	return ""
}
