package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
)

type EnrollmentService struct {
	repo repositories.EnrollmentRepository
}

func NewEnrollmentService(repo repositories.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{repo: repo}
}

func (s *EnrollmentService) GetAllEnrollments() ([]models.Enrollment, error) {
	return s.repo.GetAll()
}

func (s *EnrollmentService) GetEnrollmentByID(enrollmentID string) (*models.Enrollment, error) {
	return s.repo.GetByID(enrollmentID)
}

func (s *EnrollmentService) GetEnrollmentsByStudentID(studentID string) ([]models.Enrollment, error) {
	return s.repo.GetByStudentID(studentID)
}

func (s *EnrollmentService) GetEnrollmentsByClassID(classID string) ([]models.Enrollment, error) {
	return s.repo.GetByClassID(classID)
}

func (s *EnrollmentService) CreateEnrollment(enrollment *models.Enrollment) error {
	return s.repo.Create(enrollment)
}

func (s *EnrollmentService) UpdateEnrollment(enrollmentID string, enrollment *models.Enrollment) error {
	affected, err := s.repo.Update(enrollmentID, enrollment)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *EnrollmentService) DeleteEnrollment(enrollmentID string) error {
	affected, err := s.repo.Delete(enrollmentID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
