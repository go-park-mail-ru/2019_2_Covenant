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
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name", "T_path"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World", "path").
			AddRow(2, 2, 2, "bad guy", "3:14", "path", "Billie Eilish", "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?", "path").
			AddRow(3, 3, 3, "Still Loving You", "6:28", "path", "Scorpions", "Love at First Sting", "path")

		count := uint64(3)
		offset := uint64(0)

		mock.ExpectQuery("SELECT").WithArgs(count, offset).WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count, offset)

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
		offset := uint64(0)
		mock.ExpectQuery("SELECT").WithArgs(count, offset).WillReturnError(fmt.Errorf("some error"))

		tracks, err := trackRepo.Fetch(count, offset)

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
		offset := uint64(0)
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rows := sqlmock.NewRows(columns).
			AddRow(-1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World")

		mock.ExpectQuery("SELECT").WithArgs(count, offset).WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count, offset)

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
		offset := uint64(0)
		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World").
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		tracks, err := trackRepo.Fetch(count, offset)

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

func TestTrackRepository_StoreFavourite(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	trackRepo := configureTrackReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)
		mock.ExpectExec("INSERT").WithArgs(userID, trackID).WillReturnResult(result).WillReturnError(nil)

		err := trackRepo.StoreFavourite(userID, trackID)

		if  err != nil {
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error already exist", func(t2 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(uint64(1))
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		err := trackRepo.StoreFavourite(userID, trackID)

		if  fmt.Sprintln(err) == "already exist" {
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error inserting", func(t3 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)
		mock.ExpectExec("INSERT").WithArgs(userID, trackID).WillReturnResult(result).WillReturnError(fmt.Errorf("some error"))

		err := trackRepo.StoreFavourite(userID, trackID)

		if  err == nil {
			fmt.Println("Error -> expected nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})
}

func TestTrackRepository_RemoveFavourite(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	trackRepo := configureTrackReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		favID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(favID)
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)
		mock.ExpectExec("DELETE").WithArgs(favID).WillReturnResult(result).WillReturnError(nil)

		err := trackRepo.RemoveFavourite(userID, trackID)

		if  err != nil {
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Error not found", func(t2 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		err := trackRepo.RemoveFavourite(userID, trackID)

		if  fmt.Sprintln(err) == "not found" {
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Error deleting", func(t3 *testing.T) {
		userID := uint64(1)
		trackID := uint64(1)
		favID := uint64(1)
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(favID)
		mock.ExpectQuery("SELECT").WithArgs(userID, trackID).WillReturnRows(rows)

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)
		mock.ExpectExec("DELETE").WithArgs(favID).WillReturnResult(result).WillReturnError(fmt.Errorf("some error"))

		err := trackRepo.RemoveFavourite(userID, trackID)

		if  err == nil {
			fmt.Println("Error -> expected nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})
}

func TestTrackRepository_FetchFavourites(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	trackRepo := configureTrackReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		userID := uint64(1)
		count := uint64(3)
		offset := uint64(0)

		rowCount := sqlmock.NewRows([]string{"count"}).AddRow(uint64(3))
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rowCount)

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name", "T_path"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World", "path").
			AddRow(2, 2, 2, "bad guy", "3:14", "path", "Billie Eilish", "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?", "path").
			AddRow(3, 3, 3, "Still Loving You", "6:28", "path", "Scorpions", "Love at First Sting", "path")

		mock.ExpectQuery("SELECT").WithArgs(userID, count, offset).WillReturnRows(rows)

		tracks, total, err := trackRepo.FetchFavourites(userID, count, offset)

		if tracks == nil || err != nil || total != uint64(3) {
			fmt.Println("Tracks -> expected not nil, got: ", tracks)
			fmt.Println("Total -> expected 3, got: ", total)
			fmt.Println("Error -> expected nil, got: ", err)
			t1.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t1.Fail()
		}
	})

	t.Run("Test with errors", func(t2 *testing.T) {
		userID := uint64(1)
		count := uint64(1)
		offset := uint64(0)

		rowCount := sqlmock.NewRows([]string{"count"}).AddRow(uint64(3))
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rowCount)

		mock.ExpectQuery("SELECT").WithArgs(userID, count, offset).WillReturnError(fmt.Errorf("some error"))

		tracks, total, err := trackRepo.FetchFavourites(userID, count, offset)

		if tracks != nil || err == nil || total != uint64(3) {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Total -> expected 3, got: ", total)
			fmt.Println("Error -> expected not nil, got: ", err)
			t2.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t2.Fail()
		}
	})

	t.Run("Test with scanning error", func(t3 *testing.T) {
		userID := uint64(1)
		count := uint64(1)
		offset := uint64(0)

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}

		rowCount := sqlmock.NewRows([]string{"count"}).AddRow(uint64(3))
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rowCount)

		rows := sqlmock.NewRows(columns).
			AddRow(-1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World")

		mock.ExpectQuery("SELECT").WithArgs(userID, count, offset).WillReturnRows(rows)

		tracks, total, err := trackRepo.FetchFavourites(userID, count, offset)

		if tracks != nil || err == nil || total != uint64(3) {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Total -> expected 3, got: ", total)
			fmt.Println("Error -> expected not nil, got: ", err)
			t3.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t3.Fail()
		}
	})

	t.Run("Test with error of rows", func(t4 *testing.T) {
		userID := uint64(1)
		count := uint64(1)
		offset := uint64(0)

		rowCount := sqlmock.NewRows([]string{"count"}).AddRow(uint64(3))
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rowCount)

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World").
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		tracks, total, err := trackRepo.FetchFavourites(userID, count, offset)

		if tracks != nil || err == nil || total != uint64(3) {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Total -> expected 3, got: ", total)
			fmt.Println("Error -> expected not nil, got: ", err)
			t4.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t4.Fail()
		}
	})

	t.Run("Test with error of rows", func(t5 *testing.T) {
		userID := uint64(1)
		count := uint64(1)
		offset := uint64(0)

		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnError(fmt.Errorf("some error"))

		tracks, total, err := trackRepo.FetchFavourites(userID, count, offset)

		if tracks != nil || err == nil || total != 0 {
			fmt.Println("Tracks -> expected nil, got: ", tracks)
			fmt.Println("Total -> expected 0, got: ", total)
			fmt.Println("Error -> expected not nil, got: ", err)
			t5.Fail()
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Println("unmet expectation error: ", err)
			t5.Fail()
		}
	})
}

func TestTrackRepository_FindLike(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer dbMock.Close()

	trackRepo := configureTrackReposirory(dbMock)

	t.Run("Test OK", func(t1 *testing.T) {
		name := "a"
		count := uint64(3)

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name", "T_path"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World", "path").
			AddRow(2, 2, 2, "bad guy", "3:14", "path", "Billie Eilish", "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?", "path").
			AddRow(3, 3, 3, "Still Loving You", "6:28", "path", "Scorpions", "Love at First Sting", "path")

		mock.ExpectQuery("SELECT").WithArgs(name, count).WillReturnRows(rows)

		tracks, err := trackRepo.FindLike(name, count)

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

		mock.ExpectQuery("SELECT").WithArgs(name, count).WillReturnError(fmt.Errorf("some error"))

		tracks, err := trackRepo.FindLike(name, count)

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

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 1, 1, "We Are the Champions", "3:00", "path", "Queen", "News of the World").
			RowError(0, fmt.Errorf("some error"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		tracks, err := trackRepo.FindLike(name, count)

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

		columns := []string{"T_id", "T_album_id", "Ar_id", "T_name", "T_duration", "Al_photo", "Ar_name", "Al_name"}
		rows := sqlmock.NewRows(columns)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		tracks, err := trackRepo.FindLike(name, count)

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
