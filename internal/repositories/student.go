package repositories

import (
	"attendance-management/internal/models"
	"database/sql"
)

type StudentRepository interface {
	GetAll() ([]models.Student, error)
	GetByID(id string) (*models.Student, error)
	Create(student *models.Student) error
	Update(id string, student *models.Student) (int64, error)
	Delete(id string) (int64, error)
}

type studentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{DB: db}
}

func (r *studentRepository) GetAll() ([]models.Student, error) {
	rows, err := r.DB.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.StudentID, &s.Name, &s.Grade, &s.Phone, &s.ParentPhone); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *studentRepository) GetByID(id string) (*models.Student, error) {
	row := r.DB.QueryRow("SELECT * FROM students WHERE student_id = $1", id)

	var s models.Student
	if err := row.Scan(&s.StudentID, &s.Name, &s.Grade, &s.Phone, &s.ParentPhone); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepository) Create(student *models.Student) error {
	query := `
		INSERT INTO students (student_id, name, grade, phone, parent_phone)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(query,
		student.StudentID, student.Name, student.Grade, student.Phone, student.ParentPhone,
	)
	return err
}

func (r *studentRepository) Update(id string, student *models.Student) (int64, error) {
	result, err := r.DB.Exec(`
		UPDATE students SET name=$1, grade=$2, phone=$3, parent_phone=$4 WHERE student_id=$5`,
		student.Name, student.Grade, student.Phone, student.ParentPhone, id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *studentRepository) Delete(id string) (int64, error) {
	result, err := r.DB.Exec("DELETE FROM students WHERE student_id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
