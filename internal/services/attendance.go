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
	attendanceRepo      repositories.AttendanceRepository
	classRepo           repositories.ClassRepository
	notificationService NotificationService
}

func NewAttendanceService(attendanceRepo repositories.AttendanceRepository, classRepo repositories.ClassRepository, notificationService NotificationService) *AttendanceService {
	return &AttendanceService{
		attendanceRepo:      attendanceRepo,
		classRepo:           classRepo,
		notificationService: notificationService,
	}
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
	// 등원 기록이 없는 경우 등원 기록 생성(등원 시간 추가하여 등원 처리)
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
		s.notificationService.SendAttendanceMessage(strStudentID, "등원", attendance.CheckIn.Time)
		return s.attendanceRepo.Create(attendance)
	} else {
		// 등원 기록이 있는 경우 등원 기록 업데이트(하원 시간 추가하여 하원 처리)
		if existingAttendance.CheckOut.Time.Format("15:04") == "00:00" {
			if existingAttendance.CheckIn.Time.After(attendance.CheckIn.Time) {
				return fmt.Errorf("하원 시간은 등원 시간보다 늦어야 합니다")
			} else {
				// 현재 입력의 등원 시간을 하원 시간으로 변경
				existingAttendance.CheckOut = attendance.CheckIn
			}
			_, err := s.attendanceRepo.Update(strStudentID, attendance.Date.Time.Format("2006-01-02"), existingAttendance)
			s.notificationService.SendAttendanceMessage(strStudentID, "하원", existingAttendance.CheckOut.Time)
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
