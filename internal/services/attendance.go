package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
	"fmt"
	"strconv"
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
	strStudentID := strconv.Itoa(attendance.StudentID)

	existingAttendance, err := s.attendanceRepo.GetByStudentIDAndDate(strStudentID, attendance.Date.Time.Format("2006-01-02"))
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingAttendance == nil {
		today := time.Now().In(time.Local).Weekday()
		Weekday := convertWeekdayToKorean(today)
		classes, err := s.classRepo.GetClassesForStudentByWeekday(strStudentID, Weekday)
		if err != nil {
			return err
		}
		if len(classes) == 0 {
			attendance.Status = "출석"
		} else {
			attendance.Status = determineAttendanceStatus(attendance.CheckIn.Time, classes[0].StartTime.Time)
		}
		return s.attendanceRepo.Create(attendance)
	} else {
		if existingAttendance.CheckOut.Time.Format("15:04") == "00:00" {
			existingAttendance.CheckOut = attendance.CheckIn
			_, err := s.attendanceRepo.Update(strStudentID, attendance.Date.Time.Format("2006-01-02"), existingAttendance)
			return err
		} else {
			return fmt.Errorf("출석 기록이 이미 존재합니다")
		}
	}
}

func determineAttendanceStatus(checkInTime time.Time, classStartTime time.Time) string {
	checkInTime = time.Date(0, 0, 0, checkInTime.Hour(), checkInTime.Minute(), 0, 0, time.Local)
	classStartTime = time.Date(0, 0, 0, classStartTime.Hour(), classStartTime.Minute(), 0, 0, time.Local)
	if checkInTime.After(classStartTime) {
		return "지각"
	}
	return "출석"
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
