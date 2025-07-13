package repositories

import (
	"attendance-management/internal/models"
	"database/sql"
)

type TeacherRepository interface {
	Create(*models.Teacher) error
	GetAll() ([]models.Teacher, error)
	GetByID(teacherID string) (*models.Teacher, error)
	Update(teacherID string, teacher *models.Teacher) (int64, error)
	Delete(teacherID string) (int64, error)
}

type teacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) Create(teacher *models.Teacher) error {
	query := `
		INSERT INTO teachers (teacher_id, password, name, phone_number)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, teacher.TeacherID, teacher.Password, teacher.Name, teacher.Phone)
	return err
}

func (r *teacherRepository) GetAll() ([]models.Teacher, error) {
	query := `SELECT teacher_id, name, phone_number FROM teachers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(&teacher.TeacherID, &teacher.Name, &teacher.Phone); err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (r *teacherRepository) GetByID(teacherID string) (*models.Teacher, error) {
	query := `SELECT teacher_id, name, phone_number FROM teachers WHERE teacher_id = $1`
	row := r.db.QueryRow(query, teacherID)

	var teacher models.Teacher
	err := row.Scan(&teacher.TeacherID, &teacher.Name, &teacher.Phone)
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) Update(teacherID string, teacher *models.Teacher) (int64, error) {
	query := `
		UPDATE teachers SET name = $1, phone_number = $2, password = $3 
		WHERE teacher_id = $4
	`
	result, err := r.db.Exec(query, teacher.Name, teacher.Phone, teacher.Password, teacherID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *teacherRepository) Delete(teacherID string) (int64, error) {
	query := `DELETE FROM teachers WHERE teacher_id = $1`
	result, err := r.db.Exec(query, teacherID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
