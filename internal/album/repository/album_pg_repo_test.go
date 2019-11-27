package repository

import (
	"2019_2_Covenant/internal/album"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func configureAlbumReposirory(db *sql.DB) album.Repository {
	return NewAlbumRepository(db)
}


func TestAlbumRepository_FindLike(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	albumRepo := configureAlbumReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		name := "a"
		count := uint64(3)

		columns := []string{"id", "artist_id", "name", "photo", "year"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, "News of the World", "path", time.Now()).
			AddRow(2, 2, "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?", "path", time.Now()).
			AddRow(3, 3, "Love at First Sting", "path", time.Now())

		mock.ExpectQuery("select").WithArgs(name, count).WillReturnRows(rows)

		tracks, err := albumRepo.FindLike(name, count)

		if tracks == nil || err != nil {
			fmt.Println("Tracks -> expected not nil, got: ", tracks)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Test with errors", func(t2 *testing.T) {
		name := "a"
		count := uint64(3)

		mock.ExpectQuery("select").WithArgs(name, count).WillReturnError(fmt.Errorf("some error"))

		tracks, err := albumRepo.FindLike(name, count)

		if tracks != nil || err == nil {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Test with error of rows", func(t4 *testing.T) {
		name := "a"
		count := uint64(3)

		columns := []string{"id", "artist_id", "name", "photo", "year"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, "News of the World", "path", time.Now()).
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("select").WillReturnRows(rows)

		tracks, err := albumRepo.FindLike(name, count)

		if tracks != nil || err == nil {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Error -> expected not nil, got: ", err)
			t4.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t4.Fail()
		}
	})

	t.Run("Test with no result", func(t5 *testing.T) {
		name := "a"
		count := uint64(3)

		columns := []string{"id", "artist_id", "name", "photo", "year"}
		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("select").WillReturnRows(rows)

		tracks, err := albumRepo.FindLike(name, count)

		if tracks != nil || err != nil {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Error -> expected nil, got: ", err)
			t5.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t5.Fail()
		}
	})
}
