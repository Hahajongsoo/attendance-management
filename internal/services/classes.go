package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories"
	"database/sql"
)

type ClassService struct {
	repo repositories.ClassRepository
}

func NewClassService(repo repositories.ClassRepository) *ClassService {
	return &ClassService{repo: repo}
}

func (s *ClassService) GetAllClasses() ([]models.Class, error) {
	return s.repo.GetAll()
}

func (s *ClassService) GetClassByID(classID string) (*models.Class, error) {
	return s.repo.GetByID(classID)
}

func (s *ClassService) GetClassesByTeacherID(teacherID string) ([]models.Class, error) {
	return s.repo.GetByTeacherID(teacherID)
}

func (s *ClassService) CreateClass(class *models.Class) error {
	return s.repo.Create(class)
}

func (s *ClassService) UpdateClass(classID string, class *models.Class) error {
	affected, err := s.repo.Update(classID, class)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *ClassService) DeleteClass(classID string) error {
	affected, err := s.repo.Delete(classID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
