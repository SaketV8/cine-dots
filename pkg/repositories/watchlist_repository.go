package repositories

import (
	"database/sql"
	"time"

	"github.com/saketV8/cine-dots/pkg/models"
)

type WatchListModelInterface interface {
	GetAllWatchList() ([]models.Watchlist, error)
	GetWatchedList() ([]models.Watchlist, error)
	GetWatchingList() ([]models.Watchlist, error)
	GetNotWatchedList() ([]models.Watchlist, error)
	GetWatchListById(watchlist_id string) (models.Watchlist, error)

	AddWatchList(watchList models.Watchlist) (models.Watchlist, error)
	DeleteWatchList(watchList models.WatchListDeleteRequest) (int, error)
	UpdateWatchList(watchList models.WatchListUpdateRequest) (int, error)
}

type WatchListModel struct {
	DB *sql.DB
}

func (watchListModel *WatchListModel) GetAllWatchList() ([]models.Watchlist, error) {
	statement := `SELECT watchlist_id, title, release_year, genre, director, status, added_date FROM Watchlist;`

	rows, err := watchListModel.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	// empty watchList model slice
	watchLists := []models.Watchlist{}
	for rows.Next() {
		// empty watchList model
		watchList := models.Watchlist{}

		// setting data to the watchList model from the rows
		err := rows.Scan(
			&watchList.WatchlistID,
			&watchList.Title,
			&watchList.ReleaseYear,
			&watchList.Genre,
			&watchList.Director,
			&watchList.Status,
			&watchList.AddedDate,
		)
		if err != nil {
			return nil, err
		}

		// adding watchList data to watchLists slice
		watchLists = append(watchLists, watchList)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return watchLists, nil
}

func (watchListModel *WatchListModel) GetWatchedList() ([]models.Watchlist, error) {
	statement := `SELECT watchlist_id, title, release_year, genre, director, status, added_date FROM Watchlist WHERE status = "watched";`

	rows, err := watchListModel.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	// empty watchList model slice
	watchLists := []models.Watchlist{}
	for rows.Next() {
		// empty watchList model
		watchList := models.Watchlist{}

		// setting data to the watchList model from the rows
		err := rows.Scan(
			&watchList.WatchlistID,
			&watchList.Title,
			&watchList.ReleaseYear,
			&watchList.Genre,
			&watchList.Director,
			&watchList.Status,
			&watchList.AddedDate,
		)
		if err != nil {
			return nil, err
		}

		// adding watchList data to watchLists slice
		watchLists = append(watchLists, watchList)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return watchLists, nil
}

func (watchListModel *WatchListModel) GetWatchingList() ([]models.Watchlist, error) {
	statement := `SELECT watchlist_id, title, release_year, genre, director, status, added_date FROM Watchlist WHERE status = "watching";`

	rows, err := watchListModel.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	// empty watchList model slice
	watchLists := []models.Watchlist{}
	for rows.Next() {
		// empty watchList model
		watchList := models.Watchlist{}

		// setting data to the watchList model from the rows
		err := rows.Scan(
			&watchList.WatchlistID,
			&watchList.Title,
			&watchList.ReleaseYear,
			&watchList.Genre,
			&watchList.Director,
			&watchList.Status,
			&watchList.AddedDate,
		)
		if err != nil {
			return nil, err
		}

		// adding watchList data to watchLists slice
		watchLists = append(watchLists, watchList)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return watchLists, nil
}

func (watchListModel *WatchListModel) GetNotWatchedList() ([]models.Watchlist, error) {
	statement := `SELECT watchlist_id, title, release_year, genre, director, status, added_date FROM Watchlist WHERE status = "not watched";`

	rows, err := watchListModel.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	// empty watchList model slice
	watchLists := []models.Watchlist{}
	for rows.Next() {
		// empty watchList model
		watchList := models.Watchlist{}

		// setting data to the watchList model from the rows
		err := rows.Scan(
			&watchList.WatchlistID,
			&watchList.Title,
			&watchList.ReleaseYear,
			&watchList.Genre,
			&watchList.Director,
			&watchList.Status,
			&watchList.AddedDate,
		)
		if err != nil {
			return nil, err
		}

		// adding watchList data to watchLists slice
		watchLists = append(watchLists, watchList)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return watchLists, nil
}

// GetWatchListById
func (watchListModel *WatchListModel) GetWatchListById(watchlist_id string) (models.Watchlist, error) {
	statement := `SELECT watchlist_id, title, release_year, genre, director, status, added_date FROM Watchlist WHERE watchlist_id = ?;`
	// empty watchList model
	watchList := models.Watchlist{}

	err := watchListModel.DB.QueryRow(statement, watchlist_id).Scan(
		&watchList.WatchlistID,
		&watchList.Title,
		&watchList.ReleaseYear,
		&watchList.Genre,
		&watchList.Director,
		&watchList.Status,
		&watchList.AddedDate,
	)
	if err != nil {
		return watchList, err
	}

	return watchList, nil
}

func (watchListModel *WatchListModel) AddWatchList(watchList models.Watchlist) (models.Watchlist, error) {
	statement := `INSERT INTO Watchlist (title, release_year, genre, director, status) VALUES (?, ?, ?, ?, ?);`

	watchListResult := models.Watchlist{}

	result, err := watchListModel.DB.Exec(statement, watchList.Title, watchList.ReleaseYear, watchList.Genre, watchList.Director, watchList.Status)
	if err != nil {
		return models.Watchlist{}, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return models.Watchlist{}, err
	}

	// adding Id to model inserted autmatically by sqlite before returning to end user
	watchListResult.WatchlistID = int(lastInsertedId)
	watchListResult.Title = watchList.Title
	watchListResult.ReleaseYear = watchList.ReleaseYear
	watchListResult.Genre = watchList.Genre
	watchListResult.Director = watchList.Director
	watchListResult.Status = watchList.Status
	watchListResult.AddedDate = watchList.AddedDate

	return watchListResult, nil
}

func (watchListModel *WatchListModel) DeleteWatchList(watchList models.WatchListDeleteRequest) (int, error) {
	statement := `DELETE FROM Watchlist WHERE watchlist_id = ?;`

	result, err := watchListModel.DB.Exec(statement, watchList.WatchlistID)
	if err != nil {
		return 0, err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowAffected), nil
}

func (watchListModel *WatchListModel) UpdateWatchList(watchList models.WatchListUpdateRequest) (int, error) {
	statement := `UPDATE Watchlist SET title = ?, release_year = ?, genre = ?, director = ?, status = ?, added_date = ? WHERE watchlist_id = ?;`

	result, err := watchListModel.DB.Exec(statement, watchList.Title, watchList.ReleaseYear, watchList.Genre, watchList.Director, watchList.Status, time.Now(), watchList.WatchlistID)
	if err != nil {
		return 0, err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowAffected), nil
}
