package tester

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goda6565/ptf-backends/applications/auth/pkg/logger"
)

func MockDB() (db *gorm.DB, mock sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}

	mockGormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:        "sqlmock_db",
		DriverName: "postgres",
		Conn:       mockDB,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mockGormDB, mock
}

type mockClock struct {
	t time.Time
}

func NewMockClock(t time.Time) *mockClock {
	return &mockClock{t: t}
}

func (m *mockClock) Now() time.Time {
	return m.t
}
