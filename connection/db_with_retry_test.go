package connection

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type DbWithRetryTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

const (
	defaultRetry = 3
	defaultDelay = time.Duration(200) * time.Millisecond
)

func TestDbWithRetryTestSuite(s *testing.T) {
	suite.Run(s, new(DbWithRetryTestSuite))
}

func (s *DbWithRetryTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	config := DBWithRetryConf{
		DefaultRetry: defaultRetry,
		DefaultDelay: defaultDelay,
	}

	retryDb := NewDBWithRetry(db, config)

	s.DB, err = gorm.Open("postgres", retryDb)
	if err != nil {
		s.T().Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	s.DB.LogMode(true)
}

func (s *DbWithRetryTestSuite) TestDbRetryWithoutTransaction() {
	// Expect to retry 3 times: 2 failed and last 1 success
	for i := 0; i <= defaultRetry; i++ {
		if i != defaultRetry {
			s.mock.ExpectExec("UPDATE products").
				WillReturnError(fmt.Errorf("connection reset by peer [%d]", i))
		} else {
			s.mock.ExpectExec("UPDATE products").
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}

	if err := s.DB.Exec("UPDATE products SET views = views + 1").Error; err != nil {
		s.T().Errorf("Expect no error return but got error: %s", err)
	}

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *DbWithRetryTestSuite) TestDbRetryWithTransaction() {
	userID, productID := 2, 3

	s.mock.ExpectBegin()
	// Expect to retry 3 times: 2 failed and last 1 success
	for i := 0; i <= defaultRetry; i++ {
		if i != defaultRetry {
			s.mock.ExpectExec("UPDATE products").
				WillReturnError(fmt.Errorf("connection reset by peer [%d]", i))
		} else {
			s.mock.ExpectExec("UPDATE products").
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}

	s.mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	if err := recordStats(s.DB, userID, productID); err != nil {
		s.T().Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func recordStats(db *gorm.DB, userID, productID int) (err error) {
	tx := db.Begin()
	err = tx.Error
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit().Error
		default:
			tx.Rollback()
		}
	}()

	if err = tx.Exec("UPDATE products SET views = views + 1").Error; err != nil {
		return
	}
	if err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", userID, productID).Error; err != nil {
		return
	}
	return
}
