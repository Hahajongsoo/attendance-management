package services

import (
	"attendance-management/internal/models"
	"attendance-management/internal/repositories/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAttendance_WhenClassExistsAndIsLate(t *testing.T) {
	mockClassRepo := new(mocks.ClassRepository)
	mockAttendanceRepo := new(mocks.AttendanceRepository)

	now := time.Date(2025, 7, 14, 10, 5, 0, 0, time.Local)
	classStart := models.TimeOnly{Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.Local)}

	mockClassRepo.
		On("GetClassesForStudentByWeekday", 1, "월").
		Return([]models.Class{
			{ClassID: 1, StartTime: classStart},
		}, nil)

	mockAttendanceRepo.
		On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
			return att.Status == "지각"
		})).
		Return(nil)

	service := &AttendanceService{
		classRepo:      mockClassRepo,
		attendanceRepo: mockAttendanceRepo,
	}

	att := &models.Attendance{
		StudentID: 1,
		CheckIn:   models.TimeOnly{Time: now},
	}

	err := service.CreateAttendance(att)
	assert.NoError(t, err)

	mockClassRepo.AssertExpectations(t)
	mockAttendanceRepo.AssertExpectations(t)
}

func TestCreateAttendance_WhenClassExistsAndIsNotLate(t *testing.T) {
	mockClassRepo := new(mocks.ClassRepository)
	mockAttendanceRepo := new(mocks.AttendanceRepository)

	now := time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)
	classStart := models.TimeOnly{Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.Local)}

	mockClassRepo.
		On("GetClassesForStudentByWeekday", 1, "월").
		Return([]models.Class{
			{ClassID: 1, StartTime: classStart},
		}, nil)

	mockAttendanceRepo.
		On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
			return att.Status == "출석"
		})).
		Return(nil)

	service := &AttendanceService{
		classRepo:      mockClassRepo,
		attendanceRepo: mockAttendanceRepo,
	}

	att := &models.Attendance{
		StudentID: 1,
		CheckIn:   models.TimeOnly{Time: now},
	}

	err := service.CreateAttendance(att)
	assert.NoError(t, err)

	mockClassRepo.AssertExpectations(t)
	mockAttendanceRepo.AssertExpectations(t)
}

func TestCreateAttendance_WhenClassDoesNotExist(t *testing.T) {
	mockClassRepo := new(mocks.ClassRepository)
	mockAttendanceRepo := new(mocks.AttendanceRepository)

	now := time.Date(2025, 7, 14, 10, 0, 0, 0, time.Local)

	mockClassRepo.
		On("GetClassesForStudentByWeekday", 1, "월").
		Return([]models.Class{}, nil)

	mockAttendanceRepo.
		On("Create", mock.MatchedBy(func(att *models.Attendance) bool {
			return att.Status == "출석"
		})).
		Return(nil)

	service := &AttendanceService{
		classRepo:      mockClassRepo,
		attendanceRepo: mockAttendanceRepo,
	}

	att := &models.Attendance{
		StudentID: 1,
		CheckIn:   models.TimeOnly{Time: now},
	}

	err := service.CreateAttendance(att)
	assert.NoError(t, err)

	mockClassRepo.AssertExpectations(t)
	mockAttendanceRepo.AssertExpectations(t)
}
