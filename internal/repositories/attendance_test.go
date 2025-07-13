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

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO attendance (student_id, date, check_in, check_out, status) VALUES ($1, $2, $3, $4, $5)")).WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestAttendanceRepository_GetByStudentID(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewAttendanceRepository(db)

	attendanceTime := time.Now()
	expectedDate := attendanceTime.Truncate(time.Hour * 24)
	expectedTime := attendanceTime.Truncate(time.Minute)

	rows := sqlmock.NewRows([]string{"student_id", "date", "check_in", "check_out", "status"}).
		AddRow(1, expectedDate, expectedTime, expectedTime, "present")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM attendance WHERE student_id = $1")).
		WithArgs("1").
		WillReturnRows(rows)

	attendances, err := repo.GetByStudentIDAndDate("1", "2025-01-01")
	assert.NoError(t, err)
	assert.Equal(t, 1, attendances.StudentID)
	assert.Equal(t, expectedDate, attendances.Date.Time)
	assert.Equal(t, expectedTime, attendances.CheckIn.Time)
	assert.Equal(t, expectedTime, attendances.CheckOut.Time)
	assert.Equal(t, "present", attendances.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}
