package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func configureUserRepository(db *sql.DB) user.Repository {
	return NewUserRepository(db)
}

func TestUserRepository_Fetch(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, "marshal", "m1@ya.ru", "path", "123456").
			AddRow(2, "plaksenka", "p2@ya.ru", "path", "123456").
			AddRow(3, "svya", "6:28", "path", "s3@ya.ru")

		count := uint64(3)

		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnRows(rows)

		users, err := userRepo.Fetch(count)

		if users == nil || err != nil {
			fmt.Println("Users -> expected not nil, got: ", users)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Errors", func(t2 *testing.T) {
		count := uint64(1)
		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnError(fmt.Errorf("some error"))

		users, err := userRepo.Fetch(count)

		if users != nil || err == nil {
			fmt.Println("Users -> expected nil, got: ", users)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Test with scanning error", func(t3 *testing.T) {
		count := uint64(1)
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, "marshal", "m1@ya.ru", "path", "123456")
		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnRows(rows)

		users, err := userRepo.Fetch(count)

		if users != nil || err == nil {
			fmt.Println("Users -> expected nil, got: ", users)
			fmt.Println("Error -> expected not nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})

	t.Run("Test with error of rows", func(t4 *testing.T) {
		count := uint64(1)
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, "marshal", "m1@ya.ru", "path", "123456").
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		users, err := userRepo.Fetch(count)

		if users != nil || err == nil {
			fmt.Println("Users -> expected nil, got: ", users)
			fmt.Println("Error -> expected not nil, got: ", err)
			t4.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t4.Fail()
		}
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		email := "m1@ya.ru"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(email).WillReturnRows(rows)

		getUser, err := userRepo.GetByEmail(email)

		if getUser == nil || err != nil {
			fmt.Println("User -> expected not nil, got: ", getUser)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		email := "m1@ya.ru"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(email).WillReturnRows(rows)

		getUser, err := userRepo.GetByEmail(email)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error not exist", func(t2 *testing.T) {
		email := "m1@ya.ru"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("SELECT").WithArgs(email).WillReturnRows(rows)

		getUser, err := userRepo.GetByEmail(email)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		ID := uint64(1)
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnRows(rows)

		getUser, err := userRepo.GetByID(ID)

		if getUser == nil || err != nil {
			fmt.Println("User -> expected not nil, got: ", getUser)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		ID := uint64(1)
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnRows(rows)

		getUser, err := userRepo.GetByID(ID)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error not exist", func(t3 *testing.T) {
		ID := uint64(1)
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnRows(rows)

		getUser, err := userRepo.GetByID(ID)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})
}

func TestUserRepository_GetByNickname(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		nickname := "nickname"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(nickname).WillReturnRows(rows)

		getUser, err := userRepo.GetByNickname(nickname)

		if getUser == nil || err != nil {
			fmt.Println("User -> expected not nil, got: ", getUser)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		nickname := "nickname"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, "marshal", "m1@ya.ru", "path", "123456")

		mock.ExpectQuery("SELECT").WithArgs(nickname).WillReturnRows(rows)

		getUser, err := userRepo.GetByNickname(nickname)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error not exist", func(t2 *testing.T) {
		nickname := "nickname"
		columns := []string{"id", "nickname", "email", "avatar", "password"}

		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("SELECT").WithArgs(nickname).WillReturnRows(rows)

		getUser, err := userRepo.GetByNickname(nickname)

		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})
}

func TestUserRepository_Store(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		newUser := models.User{
			Nickname: "krulex",
			Email:    "p4@mail.ru",
			Password: "123456",
		}

		rows := sqlmock.NewRows([]string{"id", "avatar"}).
			AddRow(uint64(6), "path")
		mock.ExpectQuery("INSERT").WithArgs(newUser.Nickname, newUser.Email, newUser.Password).WillReturnRows(rows)

		err := userRepo.Store(&newUser)
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
		newUser := models.User{
			Nickname: "krulex",
			Email:    "p4@mail.ru",
			Password: "123456",
		}

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(uint64(6))
		mock.ExpectQuery("INSERT").WithArgs(newUser.Nickname, newUser.Email, newUser.Password).WillReturnRows(rows)

		err := userRepo.Store(&newUser)
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

func TestUserRepository_UpdateAvatar(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		id := uint64(1)
		avatar := "new path"

		columns := []string{"nickname", "email", "avatar"}
		rows := sqlmock.NewRows(columns).
			AddRow("marshal", "m1@ya.ru", "path")

		mock.ExpectQuery("UPDATE").WithArgs(avatar, id).WillReturnRows(rows)

		getUser, err := userRepo.UpdateAvatar(id, avatar)
		if getUser == nil || err != nil {
			fmt.Println("User -> expected not nil, got: ", getUser)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		id := uint64(1)
		avatar := "new path"

		columns := []string{"nickname", "email"}
		rows := sqlmock.NewRows(columns).
			AddRow("marshal", "m1@ya.ru")

		mock.ExpectQuery("UPDATE").WithArgs(avatar, id).WillReturnRows(rows)

		getUser, err := userRepo.UpdateAvatar(id, avatar)
		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})
}

func TestUserRepository_UpdateNickname(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		id := uint64(1)
		nickname := "new nickname"
		email := "e@mail"

		columns := []string{"nickname", "email", "avatar"}
		rows := sqlmock.NewRows(columns).
			AddRow("new nickname", "e@mail", "path")

		mock.ExpectQuery("UPDATE").WithArgs(nickname, email, id).WillReturnRows(rows)

		getUser, err := userRepo.Update(id, nickname, email)
		if getUser == nil || err != nil {
			fmt.Println("User -> expected not nil, got: ", getUser)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error of scanning", func(t2 *testing.T) {
		id := uint64(1)
		nickname := "new nickname"
		email := "e@mail"

		columns := []string{"nickname", "email"}
		rows := sqlmock.NewRows(columns).
			AddRow("new nickname", "e@mail")

		mock.ExpectQuery("UPDATE").WithArgs(nickname, email, id).WillReturnRows(rows)

		getUser, err := userRepo.Update(id, nickname, email)
		if getUser != nil || err == nil {
			fmt.Println("User -> expected nil, got: ", getUser)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	userRepo := configureUserRepository(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		id := uint64(1)
		password := "new password"

		mock.ExpectExec("UPDATE").WithArgs(password, id).WillReturnResult(sqlmock.NewResult(1, 1))

		err := userRepo.UpdatePassword(id, password)
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
		id := uint64(1)
		password := "new pass"

		mock.ExpectExec("UPDATE").WithArgs(password, id).WillReturnError(fmt.Errorf("some err"))

		err := userRepo.UpdatePassword(id, password)
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
