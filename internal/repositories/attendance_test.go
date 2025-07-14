package repositories

import (
	"attendance-management/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAttendanceRepository_Create(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewAttendanceRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO attendance (student_id, date, check_in, check_out, status) VALUES ($1, $2, $3, $4, $5)")).
		WithArgs(1, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Now().Truncate(time.Minute), time.Now().Truncate(time.Minute), "present").
		WillReturnResult(sqlmock.NewResult(1, 1))

	attendance := models.Attendance{
		StudentID: 1,
		Date:      models.DateOnly{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		CheckIn:   models.TimeOnly{Time: time.Now().Truncate(time.Minute)},
		CheckOut:  models.TimeOnly{Time: time.Now().Truncate(time.Minute)},
		Status:    "present",
	}

	err := repo.Create(&attendance)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAttendanceRepository_GetByStudentIDAndDate(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewAttendanceRepository(db)

	rows := sqlmock.NewRows([]string{"student_id", "date", "check_in", "check_out", "status"}).
		AddRow(1, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Now().Truncate(time.Minute), time.Now().Truncate(time.Minute), "present")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM attendance WHERE student_id = $1 AND date = $2")).
		WithArgs("1", "2025-01-01").
		WillReturnRows(rows)

	_, err := repo.GetByStudentIDAndDate("1", "2025-01-01")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
