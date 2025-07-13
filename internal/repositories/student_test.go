package repositories

import (
	"regexp"
	"testing"

	"attendance-management/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStudentRepository_GetAll(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewStudentRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM students")).WillReturnRows(sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "중1", "010-1234-5678", "010-8765-4321"))

	students, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(students))
	assert.Equal(t, "김철수", students[0].Name)
	assert.Equal(t, "중1", students[0].Grade)
	assert.Equal(t, "010-1234-5678", students[0].Phone)
	assert.Equal(t, "010-8765-4321", students[0].ParentPhone)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentRepository_GetByID(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewStudentRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM students WHERE student_id = $1")).WillReturnRows(sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "중1", "010-1234-5678", "010-8765-4321"))

	student, err := repo.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "김철수", student.Name)
	assert.Equal(t, "중1", student.Grade)
	assert.Equal(t, "010-1234-5678", student.Phone)
	assert.Equal(t, "010-8765-4321", student.ParentPhone)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentRepository_Create(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewStudentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO students (student_id, name, grade, phone, parent_phone) VALUES ($1, $2, $3, $4, $5)")).WillReturnResult(sqlmock.NewResult(1, 1))

	student := models.Student{
		Name:        "김철수",
		Grade:       "중1",
		Phone:       "010-1234-5678",
		ParentPhone: "010-8765-4321",
	}

	err := repo.Create(&student)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentRepository_Update(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewStudentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE students SET name=$1, grade=$2, phone=$3, parent_phone=$4 WHERE student_id=$5")).WillReturnResult(sqlmock.NewResult(1, 1))

	student := models.Student{
		StudentID:   1,
		Name:        "김철수",
		Grade:       "중1",
		Phone:       "010-1234-5678",
		ParentPhone: "010-8765-4321",
	}

	affected, err := repo.Update("1", &student)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentRepository_Delete(t *testing.T) {
	db, mock := SetupMockDB(t)
	defer db.Close()

	repo := NewStudentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM students WHERE student_id = $1")).WillReturnResult(sqlmock.NewResult(1, 1))

	affected, err := repo.Delete("1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	assert.NoError(t, mock.ExpectationsWereMet())
}
