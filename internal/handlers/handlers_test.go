package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"attendance-management/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestNewHandler(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)
	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
}

func TestGetStudents(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	expectedStudents := []models.Student{
		{StudentID: 1, Name: "김철수", Grade: "1학년", Phone: "010-1234-5678", ParentPhone: "010-8765-4321"},
		{StudentID: 2, Name: "이영희", Grade: "2학년", Phone: "010-2345-6789", ParentPhone: "010-9876-5432"},
	}

	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321").
		AddRow(2, "이영희", "2학년", "010-2345-6789", "010-9876-5432")

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	handler.GetStudents(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseStudents []models.Student
	err := json.Unmarshal(w.Body.Bytes(), &responseStudents)
	require.NoError(t, err)
	assert.Equal(t, expectedStudents, responseStudents)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStudent(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	student := models.Student{
		StudentID:   3,
		Name:        "박민수",
		Grade:       "3학년",
		Phone:       "010-3456-7890",
		ParentPhone: "010-0987-6543",
	}

	mock.ExpectExec("INSERT INTO students").WithArgs(3, "박민수", "3학년", "010-3456-7890", "010-0987-6543").
		WillReturnResult(sqlmock.NewResult(1, 1))

	jsonData, _ := json.Marshal(student)
	req := httptest.NewRequest("POST", "/students", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateStudent(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStudent(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	student := models.Student{
		StudentID:   1,
		Name:        "김철수",
		Grade:       "2학년",
		Phone:       "010-1234-5678",
		ParentPhone: "010-8765-4321",
	}

	mock.ExpectExec("UPDATE students").WithArgs("김철수", "2학년", "010-1234-5678", "010-8765-4321", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	jsonData, _ := json.Marshal(student)
	req := httptest.NewRequest("PUT", "/students", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateStudent(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteStudent(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	student := models.Student{
		StudentID: 1,
	}

	mock.ExpectExec("DELETE FROM students").WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	jsonData, _ := json.Marshal(student)
	req := httptest.NewRequest("DELETE", "/students", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DeleteStudent(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentHandler(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "GET 요청",
			method:         "GET",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
					AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321")
				mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)
			},
		},
		{
			name:           "POST 요청",
			method:         "POST",
			expectedStatus: http.StatusCreated,
			setupMock: func() {
				mock.ExpectExec("INSERT INTO students").WithArgs(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "PUT 요청",
			method:         "PUT",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("UPDATE students").WithArgs("김철수", "2학년", "010-1234-5678", "010-8765-4321", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "DELETE 요청",
			method:         "DELETE",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("DELETE FROM students").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var req *http.Request
			if tt.method == "GET" {
				req = httptest.NewRequest(tt.method, "/students", nil)
			} else {
				grade := "1학년"
				if tt.method == "PUT" {
					grade = "2학년"
				}
				student := models.Student{StudentID: 1, Name: "김철수", Grade: grade, Phone: "010-1234-5678", ParentPhone: "010-8765-4321"}
				if tt.method == "DELETE" {
					student = models.Student{StudentID: 1}
				}
				jsonData, _ := json.Marshal(student)
				req = httptest.NewRequest(tt.method, "/students", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()

			handler.StudentHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetStudentsError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentsError2(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStudentInvalidJSON(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	req := httptest.NewRequest("POST", "/students", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.CreateStudent(w, req)
	})
}

func TestUpdateStudentInvalidJSON(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	req := httptest.NewRequest("PUT", "/students", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.UpdateStudent(w, req)
	})
}

func TestDeleteStudentInvalidJSON(t *testing.T) {
	db, _ := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	req := httptest.NewRequest("DELETE", "/students", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.DeleteStudent(w, req)
	})
}

// rows.Scan 에러 케이스 테스트들
func TestGetStudentsScanError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// Mock 쿼리는 성공하지만 Scan에서 에러가 발생하는 경우
	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow("invalid_int", "김철수", "1학년", "010-1234-5678", "010-8765-4321") // student_id가 문자열

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	// panic이 발생하는지 확인
	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentsNullValueError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// NULL 값이 포함된 데이터
	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow(nil, "김철수", "1학년", "010-1234-5678", "010-8765-4321") // student_id가 NULL

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	// panic이 발생하는지 확인
	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentsColumnMismatchError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// 컬럼 개수가 다른 경우 (예: 테이블 구조가 변경된 경우)
	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone", "extra_column"}).
		AddRow(1, "김철수", "1학년", "010-1234-5678", "010-8765-4321", "extra_value")

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	// panic이 발생하는지 확인
	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentsInsufficientColumnsError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// 컬럼이 부족한 경우
	rows := sqlmock.NewRows([]string{"student_id", "name", "grade"}).
		AddRow(1, "김철수", "1학년") // phone, parent_phone 컬럼이 없음

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	// panic이 발생하는지 확인
	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentsDataTypeError(t *testing.T) {
	db, mock := createMockDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// 데이터 타입이 맞지 않는 경우
	rows := sqlmock.NewRows([]string{"student_id", "name", "grade", "phone", "parent_phone"}).
		AddRow("not_an_int", "김철수", 123, "010-1234-5678", "010-8765-4321") // grade가 int, phone이 string이어야 하는데 반대

	mock.ExpectQuery("SELECT \\* FROM students").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()

	// panic이 발생하는지 확인
	assert.Panics(t, func() {
		handler.GetStudents(w, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
