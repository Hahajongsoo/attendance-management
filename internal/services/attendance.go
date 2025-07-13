package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
)

type AttendanceService struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceService(repo repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) GetAttendanceByStudentIDAndDate(studentID string, date string) (*models.Attendance, error) {
	return s.repo.GetByStudentIDAndDate(studentID, date)
}

func (s *AttendanceService) GetAttendanceByDate(date string) ([]models.Attendance, error) {
	return s.repo.GetByDate(date)
}

func (s *AttendanceService) CreateAttendance(attendance *models.Attendance) error {
	return s.repo.Create(attendance)
}

func (s *AttendanceService) UpdateAttendance(studentID string, date string, attendance *models.Attendance) error {
	affected, err := s.repo.Update(studentID, date, attendance)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *AttendanceService) DeleteAttendance(studentID string, date string) error {
	affected, err := s.repo.Delete(studentID, date)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
