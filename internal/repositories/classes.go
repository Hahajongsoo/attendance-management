package repositories

import (
	"attendance-management/internal/models"
	"database/sql"
)

type ClassRepository interface {
	Create(*models.Class) error
	GetAll() ([]models.Class, error)
	GetByID(classID string) (*models.Class, error)
	GetByTeacherID(teacherID string) ([]models.Class, error)
	Update(classID string, class *models.Class) (int64, error)
	Delete(classID string) (int64, error)
}

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	query := `
		INSERT INTO classes (class_id, class_name, days, start_time, end_time, price, teacher_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, class.ClassID, class.ClassName, class.Days, class.StartTime.Time, class.EndTime.Time, class.Price, class.TeacherID)
	return err
}

func (r *classRepository) GetAll() ([]models.Class, error) {
	query := `SELECT * FROM classes`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (r *classRepository) GetByID(classID string) (*models.Class, error) {
	query := `SELECT * FROM classes WHERE class_id = $1`
	row := r.db.QueryRow(query, classID)

	var class models.Class
	err := row.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID)
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) GetByTeacherID(teacherID string) ([]models.Class, error) {
	query := `SELECT * FROM classes WHERE teacher_id = $1`
	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ClassID, &class.ClassName, &class.Days, &class.StartTime.Time, &class.EndTime.Time, &class.Price, &class.TeacherID); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (r *classRepository) Update(classID string, class *models.Class) (int64, error) {
	query := `
		UPDATE classes SET class_name = $1, days = $2, start_time = $3, end_time = $4, price = $5, teacher_id = $6 
		WHERE class_id = $7
	`
	result, err := r.db.Exec(query, class.ClassName, class.Days, class.StartTime.Time, class.EndTime.Time, class.Price, class.TeacherID, classID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *classRepository) Delete(classID string) (int64, error) {
	query := `DELETE FROM classes WHERE class_id = $1`
	result, err := r.db.Exec(query, classID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
