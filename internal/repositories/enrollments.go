package repositories

import (
	"attendance-management/internal/models"
	"database/sql"
)

type EnrollmentRepository interface {
	Create(*models.Enrollment) error
	GetAll() ([]models.Enrollment, error)
	GetByID(enrollmentID string) (*models.Enrollment, error)
	GetByStudentID(studentID string) ([]models.Enrollment, error)
	GetByClassID(classID string) ([]models.Enrollment, error)
	Update(enrollmentID string, enrollment *models.Enrollment) (int64, error)
	Delete(enrollmentID string) (int64, error)
}

type enrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	query := `
		INSERT INTO enrollments (student_id, class_id, enrolled_date)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, enrollment.StudentID, enrollment.ClassID, enrollment.EnrolledDate.Time)
	return err
}

func (r *enrollmentRepository) GetAll() ([]models.Enrollment, error) {
	query := `SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetByID(enrollmentID string) (*models.Enrollment, error) {
	query := `SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE enrollment_id = $1`
	row := r.db.QueryRow(query, enrollmentID)

	var enrollment models.Enrollment
	err := row.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time)
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetByStudentID(studentID string) ([]models.Enrollment, error) {
	query := `SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE student_id = $1`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetByClassID(classID string) ([]models.Enrollment, error) {
	query := `SELECT enrollment_id, student_id, class_id, enrolled_date FROM enrollments WHERE class_id = $1`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.EnrollmentID, &enrollment.StudentID, &enrollment.ClassID, &enrollment.EnrolledDate.Time); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}

func (r *enrollmentRepository) Update(enrollmentID string, enrollment *models.Enrollment) (int64, error) {
	query := `
		UPDATE enrollments SET student_id = $1, class_id = $2, enrolled_date = $3 
		WHERE enrollment_id = $4
	`
	result, err := r.db.Exec(query, enrollment.StudentID, enrollment.ClassID, enrollment.EnrolledDate.Time, enrollmentID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *enrollmentRepository) Delete(enrollmentID string) (int64, error) {
	query := `DELETE FROM enrollments WHERE enrollment_id = $1`
	result, err := r.db.Exec(query, enrollmentID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
