package services

import (
	"database/sql"
	"regexp"
	"testing"

	"attendance-management/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return db, mock
}

func TestStudentService_GetAll(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := repositories.NewStudentRepository(db)
	service := NewStudentService(repo)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM students")).WillReturnRows(sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "중1", "010-1234-5678", "010-8765-4321"))

	students, err := service.GetAllStudents()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(students))
	assert.Equal(t, "김철수", students[0].Name)
	assert.Equal(t, "중1", students[0].Grade)
	assert.Equal(t, "010-1234-5678", students[0].Phone)
	assert.Equal(t, "010-8765-4321", students[0].ParentPhone)

	assert.NoError(t, mock.ExpectationsWereMet())
}
