package integration

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saketV8/cine-dots/pkg/database"
	"github.com/saketV8/cine-dots/pkg/models"
	"github.com/saketV8/cine-dots/pkg/repositories"
	"github.com/stretchr/testify/assert"
)

// setupTestDB initializes a test database connection
func setupTestDB(t *testing.T) *database.Database {
	db, err := database.InitializeDatabase("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS Watchlist (
        watchlist_id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        release_year INTEGER NOT NULL,
        genre TEXT NOT NULL,
        director TEXT NOT NULL,
        status TEXT NOT NULL,
        added_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.DB.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

// insertTestData inserts sample data for testing
func insertTestData(t *testing.T, db *sql.DB) {
	insertSQL := `
    INSERT INTO Watchlist (title, release_year, genre, director, status, added_date)
    VALUES 
    ('Test Movie 1', 2021, 'Action', 'Director 1', 'watched', ?),
    ('Test Movie 2', 2022, 'Comedy', 'Director 2', 'watching', ?),
    ('Test Movie 3', 2023, 'Drama', 'Director 3', 'not watched', ?);
    `

	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(insertSQL, now, now, now)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

func TestGetAllWatchList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	watchlists, err := repo.GetAllWatchList()

	assert.NoError(t, err)
	assert.Equal(t, 3, len(watchlists))
	assert.Equal(t, "Test Movie 1", watchlists[0].Title)
	assert.Equal(t, "Test Movie 2", watchlists[1].Title)
	assert.Equal(t, "Test Movie 3", watchlists[2].Title)
}

func TestGetWatchedList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	watchlists, err := repo.GetWatchedList()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchlists))
	assert.Equal(t, "Test Movie 1", watchlists[0].Title)
	assert.Equal(t, "watched", watchlists[0].Status)
}

func TestGetWatchingList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	watchlists, err := repo.GetWatchingList()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchlists))
	assert.Equal(t, "Test Movie 2", watchlists[0].Title)
	assert.Equal(t, "watching", watchlists[0].Status)
}

func TestGetNotWatchedList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	watchlists, err := repo.GetNotWatchedList()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchlists))
	assert.Equal(t, "Test Movie 3", watchlists[0].Title)
	assert.Equal(t, "not watched", watchlists[0].Status)
}

func TestGetWatchListById(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	watchlist, err := repo.GetWatchListById("1")
	assert.NoError(t, err)
	assert.Equal(t, 1, watchlist.WatchlistID)
	assert.Equal(t, "Test Movie 1", watchlist.Title)
	assert.Equal(t, "watched", watchlist.Status)

	_, err = repo.GetWatchListById("999")
	assert.Error(t, err)
}

func TestAddWatchList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	newWatchlist := models.Watchlist{
		Title:       "New Test Movie",
		ReleaseYear: 2024,
		Genre:       "Sci-Fi",
		Director:    "New Director",
		Status:      "not watched",
		AddedDate:   time.Now(),
	}

	added, err := repo.AddWatchList(newWatchlist)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, added.WatchlistID)
	assert.Equal(t, "New Test Movie", added.Title)

	watchlists, err := repo.GetAllWatchList()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchlists))
	assert.Equal(t, "New Test Movie", watchlists[0].Title)
}

func TestUpdateWatchList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	updateRequest := models.WatchListUpdateRequest{
		WatchlistID: 1,
		Title:       "Updated Movie",
		ReleaseYear: 2025,
		Genre:       "Updated Genre",
		Director:    "Updated Director",
		Status:      "watching",
	}

	rowsAffected, err := repo.UpdateWatchList(updateRequest)

	assert.NoError(t, err)
	assert.Equal(t, 1, rowsAffected)

	updatedWatchlist, err := repo.GetWatchListById("1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Movie", updatedWatchlist.Title)
	assert.Equal(t, 2025, updatedWatchlist.ReleaseYear)
	assert.Equal(t, "Updated Genre", updatedWatchlist.Genre)
	assert.Equal(t, "Updated Director", updatedWatchlist.Director)
	assert.Equal(t, "watching", updatedWatchlist.Status)

	nonExistentUpdate := models.WatchListUpdateRequest{
		WatchlistID: 999,
		Title:       "This should not update",
		Status:      "watched",
	}

	rowsAffected, err = repo.UpdateWatchList(nonExistentUpdate)
	assert.NoError(t, err)
	assert.Equal(t, 0, rowsAffected)
}

func TestDeleteWatchList(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	insertTestData(t, db.DB)

	repo := &repositories.WatchListModel{
		DB: db.DB,
	}

	deleteRequest := models.WatchListDeleteRequest{
		WatchlistID: 1,
	}

	rowsAffected, err := repo.DeleteWatchList(deleteRequest)

	assert.NoError(t, err)
	assert.Equal(t, 1, rowsAffected)

	watchlists, err := repo.GetAllWatchList()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(watchlists))

	for _, watchlist := range watchlists {
		assert.NotEqual(t, 1, watchlist.WatchlistID)
	}

	rowsAffected, err = repo.DeleteWatchList(deleteRequest)
	assert.NoError(t, err)
	assert.Equal(t, 0, rowsAffected)
}
