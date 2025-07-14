package repositories

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetClassesForStudentByWeekday_Success(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	class_id := 1
	weekday := "월"

	query := `
		SELECT c.* FROM classes c
		JOIN enrollments e ON c.class_id = e.class_id
		WHERE e.student_id = $1 AND c.days LIKE '%' || $2 || '%'
	`
	expectedTime := time.Date(2025, 7, 15, 10, 0, 0, 0, time.Local)

	rows := sqlmock.NewRows([]string{"class_id", "class_name", "days", "start_time", "end_time", "price", "teacher_id"}).
		AddRow(class_id, "월요일", "월/수", expectedTime, expectedTime.Add(1*time.Hour), 10000, 1)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(class_id, weekday).
		WillReturnRows(rows)

	repo := NewClassRepository(db)

	classes, err := repo.GetClassesForStudentByWeekday(class_id, weekday)
	assert.NoError(t, err)
	assert.Len(t, classes, 1)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetClassesForStudentByWeekday_Error(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	class_id := 1
	weekday := "월"

	query := `
		SELECT c.* FROM classes c
		JOIN enrollments e ON c.class_id = e.class_id
		WHERE e.student_id = $1 AND c.days LIKE '%' || $2 || '%'
	`
	rows := sqlmock.NewRows([]string{"class_id", "class_name", "days", "start_time", "end_time", "price", "teacher_id"}).
		AddRow(class_id, "월요일", "월/수", "10:00:00", "11:00:00", 10000, 1)

	repo := NewClassRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(class_id, weekday).
		WillReturnRows(rows)

	_, err := repo.GetClassesForStudentByWeekday(class_id, weekday)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Contains(t, err.Error(), "Scan")
}
