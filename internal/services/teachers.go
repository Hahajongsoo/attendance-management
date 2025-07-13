package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type TeacherService struct {
	repo repositories.TeacherRepository
}

func NewTeacherService(repo repositories.TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

func (s *TeacherService) GetAllTeachers() ([]models.Teacher, error) {
	return s.repo.GetAll()
}

func (s *TeacherService) GetTeacherByID(teacherID string) (*models.Teacher, error) {
	return s.repo.GetByID(teacherID)
}

func (s *TeacherService) CreateTeacher(teacher *models.Teacher) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	teacher.Password = string(hashedPassword)
	return s.repo.Create(teacher)
}

func (s *TeacherService) UpdateTeacher(teacherID string, teacher *models.Teacher) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	teacher.Password = string(hashedPassword)

	affected, err := s.repo.Update(teacherID, teacher)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *TeacherService) DeleteTeacher(teacherID string) error {
	affected, err := s.repo.Delete(teacherID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
