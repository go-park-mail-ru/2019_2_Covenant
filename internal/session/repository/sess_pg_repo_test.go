package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/vars"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func configureSessionRepository(db *sql.DB) session.Repository {
	return NewSessionRepository(db)
}

func TestSessionRepository_Get(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	sessRepo := configureSessionRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		data := "data"
		columns := []string{"id", "user_id", "expires", "data"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, time.Now().Add(24*time.Hour), "data")

		mock.ExpectQuery("SELECT").WithArgs(data).WillReturnRows(rows)

		sess, err := sessRepo.Get(data)

		if sess == nil || err != nil {
			fmt.Println("Session -> expected not nil, got: ", sess)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		data := "data"
		columns := []string{"id", "user_id", "expires", "data"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, 1, time.Now().Add(24*time.Hour), "data")

		mock.ExpectQuery("SELECT").WithArgs(data).WillReturnRows(rows)

		sess, err := sessRepo.Get(data)

		if sess != nil || err == nil {
			fmt.Println("Session -> expected nil, got: ", sess)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error of expiring and OK deleting", func(t3 *testing.T) {
		data := "data"
		id := 1
		columns := []string{"id", "user_id", "expires", "data"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, time.Now(), "data")

		rowsDel := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		mock.ExpectQuery("SELECT").WithArgs(data).WillReturnRows(rows)
		mock.ExpectQuery("DELETE").WithArgs(id).WillReturnRows(rowsDel)

		sess, err := sessRepo.Get(data)

		if sess != nil || err == nil {
			fmt.Println("Session -> expected nil, got: ", sess)
			fmt.Println("Error -> expected not nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})

	t.Run("Error of expiring and Error inside deleting", func(t4 *testing.T) {
		data := "data"
		id := 1
		columns := []string{"id", "user_id", "expires", "data"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, time.Now(), "data")

		mock.ExpectQuery("SELECT").WithArgs(data).WillReturnRows(rows)
		mock.ExpectQuery("DELETE").WithArgs(id)
		sess, err := sessRepo.Get(data)

		if sess != nil || err == nil {
			fmt.Println("Session -> expected nil, got: ", sess)
			fmt.Println("Error -> expected not nil, got: ", err)
			t4.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t4.Fail()
		}
	})
}

func TestSessionRepository_Store(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	sessRepo := configureSessionRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		newSession := models.Session {
			UserID: uint64(1),
			Expires: time.Now().Add(24*time.Hour),
			Data: "data",
		}

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(uint64(6))
		mock.ExpectQuery("INSERT").WithArgs(newSession.UserID, newSession.Expires, newSession.Data).WillReturnRows(rows)

		err := sessRepo.Store(&newSession)
		if err != nil {
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		newSession := models.Session {
			UserID: uint64(1),
			Expires: time.Now().Add(24*time.Hour),
			Data: "data",
		}

		mock.ExpectQuery("INSERT").WithArgs(newSession.UserID, newSession.Expires, newSession.Data)

		err := sessRepo.Store(&newSession)
		if err == nil {
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})
}

func TestSessionRepository_DeleteByID(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	sessRepo := configureSessionRepository(dbMock)

	t.Run("Error no rows", func(t1 *testing.T) {
		sessId := uint64(1)
		columns := []string{"id", "user_id", "expires", "data"}

		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("DELETE").WithArgs(sessId).WillReturnRows(rows)

		err := sessRepo.DeleteByID(sessId)

		if err != vars.ErrNotFound {
			fmt.Println("Error -> expected Not Found, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})
}
