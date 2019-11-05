package repository

import (
	"2019_2_Covenant/internal/track"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func configureTrackReposirory(db *sql.DB) track.Repository {
	return NewTrackRepository(db)
}

func TestTrackRepository_Fetch(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	trackRepo := configureTrackReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World").
			AddRow(2, 2, 2, "bad guy", "3:14", "path", "Billie Eilish", "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?").
			AddRow(3, 3, 3, "Still Loving You", "6:28", "path", "Scorpions", "Love at First Sting")

		count := uint64(3)

		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count)

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
		count := uint64(1)
		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnError(fmt.Errorf("some error"))

		tracks, err := trackRepo.Fetch(count)

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

	t.Run("Test with scanning error", func(t3 *testing.T) {
		count := uint64(1)
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World")

		mock.ExpectQuery("SELECT").WithArgs(count).WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count)

		if tracks != nil || err == nil {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
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
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World").
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count)

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
}
