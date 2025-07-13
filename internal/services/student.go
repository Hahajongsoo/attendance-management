package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
	"errors"
)

type StudentService struct {
	Repo repositories.StudentRepository
}

func NewStudentService(repo repositories.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	return s.Repo.GetAll()
}

func (s *StudentService) GetStudentByID(id string) (*models.Student, error) {
	return s.Repo.GetByID(id)
}

func (s *StudentService) CreateStudent(student *models.Student) error {
	if student.StudentID < 1000 || student.StudentID > 99999 || student.Name == "" {
		return errors.New("invalid input")
	}
	return s.Repo.Create(student)
}

func (s *StudentService) UpdateStudent(id string, student *models.Student) error {
	affected, err := s.Repo.Update(id, student)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *StudentService) DeleteStudent(id string) error {
	affected, err := s.Repo.Delete(id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
