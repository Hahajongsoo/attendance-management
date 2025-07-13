package repositories

import (
	"attendance-management/internal/models"
	"database/sql"
)

type AttendanceRepository interface {
	Create(*models.Attendance) error
	GetByStudentIDAndDate(studentID string, date string) (*models.Attendance, error)
	GetByDate(date string) ([]models.Attendance, error)
	Update(studentID string, date string, attendance *models.Attendance) (int64, error)
	Delete(studentID string, date string) (int64, error)
}

type attendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	query := `
		INSERT INTO attendance (student_id, date, check_in, check_out, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, attendance.StudentID, attendance.Date.Time, attendance.CheckIn.Time, attendance.CheckOut.Time, attendance.Status)
	return err
}

func (r *attendanceRepository) GetByStudentIDAndDate(studentID string, date string) (*models.Attendance, error) {
	query := `
		SELECT * FROM attendance WHERE student_id = $1 AND date = $2
	`
	row := r.db.QueryRow(query, studentID, date)

	var attendance models.Attendance
	err := row.Scan(&attendance.StudentID, &attendance.Date.Time, &attendance.CheckIn.Time, &attendance.CheckOut.Time, &attendance.Status)
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) GetByDate(date string) ([]models.Attendance, error) {
	query := `
		SELECT * FROM attendance WHERE date = $1 ORDER BY check_in ASC
	`
	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []models.Attendance
	for rows.Next() {
		var attendance models.Attendance
		if err := rows.Scan(&attendance.StudentID, &attendance.Date.Time, &attendance.CheckIn.Time, &attendance.CheckOut.Time, &attendance.Status); err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}

func (r *attendanceRepository) Update(studentID string, date string, attendance *models.Attendance) (int64, error) {
	query := `
		UPDATE attendance SET check_in = $1, check_out = $2, status = $3 
		WHERE student_id = $4 AND date = $5
	`
	result, err := r.db.Exec(query, attendance.CheckIn.Time, attendance.CheckOut.Time, attendance.Status, studentID, date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *attendanceRepository) Delete(studentID string, date string) (int64, error) {
	query := `
		DELETE FROM attendance WHERE student_id = $1 AND date = $2
	`
	result, err := r.db.Exec(query, studentID, date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
