package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/database"
	"github.com/saketV8/cine-dots/pkg/handlers"
	"github.com/saketV8/cine-dots/pkg/models"
	"github.com/saketV8/cine-dots/pkg/repositories"
	"github.com/saketV8/cine-dots/pkg/utils"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestAPI initializes a test API server with a real database connection
func setupTestAPI(t *testing.T) (*gin.Engine, *database.Database) {
	// Use an in-memory SQLite database for testing
	db, err := database.InitializeDatabase("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create the watchlist table in the in-memory database
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

	// Setup the router with the handlers
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create app with handlers
	watchListHandler := &handlers.WatchListHandler{
		WatchListModel: &repositories.WatchListModel{
			DB: db.DB,
		},
	}

	// Setup routes
	routerGroup := r.Group(utils.ROUTER_PREFIX)
	v1 := routerGroup.Group(utils.ROUTER_PREFIX_VERSION)
	{
		v1.GET("/watchlist/all", watchListHandler.GetAllWatchListHandler)
		v1.GET("/watchlist/watched", watchListHandler.GetWatchedListHandler)
		v1.GET("/watchlist/watching", watchListHandler.GetWatchingListHandler)
		v1.GET("/watchlist/notwatched", watchListHandler.GetNotWatchedListHandler)
		v1.GET("/watchlist/:watchlist_id", watchListHandler.GetWatchListByIdHandler)
		v1.POST("/watchlist/add", watchListHandler.AddWatchListHandler)
		v1.DELETE("/watchlist/delete", watchListHandler.DeleteWatchListHandler)
		v1.PATCH("/watchlist/update", watchListHandler.UpdateWatchListHandler)
	}

	return r, db
}

// insertTestAPIData adds sample data for API testing
func insertTestAPIData(t *testing.T, db *database.Database) {
	// Insert test data
	insertSQL := `
    INSERT INTO Watchlist (title, release_year, genre, director, status, added_date)
    VALUES 
    ('API Test Movie 1', 2021, 'Action', 'Director 1', 'watched', ?),
    ('API Test Movie 2', 2022, 'Comedy', 'Director 2', 'watching', ?),
    ('API Test Movie 3', 2023, 'Drama', 'Director 3', 'not watched', ?);
    `

	now := time.Now().Format(time.RFC3339)
	_, err := db.DB.Exec(insertSQL, now, now, now)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

func TestAPIGetAllWatchList(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	insertTestAPIData(t, db)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/watchlist/all", nil)
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var watchlists []models.Watchlist
	err := json.Unmarshal(resp.Body.Bytes(), &watchlists)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(watchlists))
	assert.Equal(t, "API Test Movie 1", watchlists[0].Title)
}

func TestAPIGetWatchListById(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	insertTestAPIData(t, db)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/watchlist/1", nil)
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var watchlist models.Watchlist
	err := json.Unmarshal(resp.Body.Bytes(), &watchlist)
	assert.NoError(t, err)
	assert.Equal(t, 1, watchlist.WatchlistID)
	assert.Equal(t, "API Test Movie 1", watchlist.Title)

	// Test non-existent ID
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/999", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestAPIAddWatchList(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	// Create request body
	newWatchlist := models.Watchlist{
		Title:       "API New Test Movie",
		ReleaseYear: 2024,
		Genre:       "Sci-Fi",
		Director:    "API Test Director",
		Status:      "not watched",
		AddedDate:   time.Now(),
	}

	body, _ := json.Marshal(newWatchlist)
	req, _ := http.NewRequest("POST", "/api/v1/watchlist/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var addedWatchlist models.Watchlist
	err := json.Unmarshal(resp.Body.Bytes(), &addedWatchlist)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, addedWatchlist.WatchlistID)
	assert.Equal(t, "API New Test Movie", addedWatchlist.Title)

	// Verify by getting all watchlists
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/all", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var watchlists []models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &watchlists)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchlists))
}

func TestAPIUpdateWatchList(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	insertTestAPIData(t, db)

	// Create request body
	updateRequest := models.WatchListUpdateRequest{
		WatchlistID: 1,
		Title:       "API Updated Movie",
		ReleaseYear: 2025,
		Genre:       "Updated Genre",
		Director:    "Updated Director",
		Status:      "watching",
	}

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PATCH", "/api/v1/watchlist/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "row-affected")

	// Verify update
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var watchlist models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &watchlist)
	assert.NoError(t, err)
	assert.Equal(t, "API Updated Movie", watchlist.Title)
	assert.Equal(t, "watching", watchlist.Status)
}

func TestAPIDeleteWatchList(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	insertTestAPIData(t, db)

	// Create request body
	deleteRequest := models.WatchListDeleteRequest{
		WatchlistID: 1,
	}

	body, _ := json.Marshal(deleteRequest)
	req, _ := http.NewRequest("DELETE", "/api/v1/watchlist/delete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")

	// Verify deletion
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/all", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var watchlists []models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &watchlists)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(watchlists))

	// Verify the specific item was deleted
	for _, watchlist := range watchlists {
		assert.NotEqual(t, 1, watchlist.WatchlistID)
	}
}

func TestAPIFilteredLists(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	insertTestAPIData(t, db)

	// Test watched list
	req, _ := http.NewRequest("GET", "/api/v1/watchlist/watched", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var watchedList []models.Watchlist
	err := json.Unmarshal(resp.Body.Bytes(), &watchedList)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchedList))
	assert.Equal(t, "watched", watchedList[0].Status)

	// Test watching list
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/watching", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var watchingList []models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &watchingList)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(watchingList))
	assert.Equal(t, "watching", watchingList[0].Status)

	// Test not watched list
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/notwatched", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var notWatchedList []models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &notWatchedList)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notWatchedList))
	assert.Equal(t, "not watched", notWatchedList[0].Status)
}

func TestAPIFullCycle(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	// 1. Add a new watchlist item
	ti, err := time.Parse(time.RFC3339, "2025-06-23T15:24:10Z")
	if err != nil {
		// handle the error (e.g., log or panic)
		log.Fatal(err)
	}
	newWatchlist := models.Watchlist{
		Title:       "Full Cycle Test Movie",
		ReleaseYear: 2024,
		Genre:       "Action",
		Director:    "Full Cycle Director",
		Status:      "not watched",
		// AddedDate:   time.Parse(time.RFC3339, "2025-06-23T15:24:10Z"),
		AddedDate: ti,
	}

	body, _ := json.Marshal(newWatchlist)
	req, _ := http.NewRequest("POST", "/api/v1/watchlist/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var addedWatchlist models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &addedWatchlist)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, addedWatchlist.WatchlistID)

	watchlistID := addedWatchlist.WatchlistID

	// 2. Get the added item to verify
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/"+strconv.Itoa(watchlistID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var retrievedWatchlist models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &retrievedWatchlist)
	assert.NoError(t, err)
	assert.Equal(t, addedWatchlist.WatchlistID, retrievedWatchlist.WatchlistID)
	assert.Equal(t, "Full Cycle Test Movie", retrievedWatchlist.Title)

	// 3. Update the item
	updateRequest := models.WatchListUpdateRequest{
		WatchlistID: watchlistID,
		Title:       "Updated Full Cycle Movie",
		ReleaseYear: 2025,
		Genre:       "Drama",
		Director:    "Updated Director",
		Status:      "watching",
	}

	body, _ = json.Marshal(updateRequest)
	req, _ = http.NewRequest("PATCH", "/api/v1/watchlist/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 4. Get the updated item to verify
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/"+strconv.Itoa(watchlistID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var updatedWatchlist models.Watchlist
	err = json.Unmarshal(resp.Body.Bytes(), &updatedWatchlist)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Full Cycle Movie", updatedWatchlist.Title)
	assert.Equal(t, "watching", updatedWatchlist.Status)

	// 5. Delete the item
	deleteRequest := models.WatchListDeleteRequest{
		WatchlistID: watchlistID,
	}

	body, _ = json.Marshal(deleteRequest)
	req, _ = http.NewRequest("DELETE", "/api/v1/watchlist/delete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 6. Verify deletion by trying to get the item
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/"+strconv.Itoa(watchlistID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code) // Should get an error as item is deleted
}

func TestAPIErrorHandling(t *testing.T) {
	// Setup
	router, db := setupTestAPI(t)
	defer db.DB.Close()

	// Test with invalid JSON in request body
	invalidJSON := []byte(`{"watchlist_id": 1, "title": "Invalid JSON`)
	req, _ := http.NewRequest("PATCH", "/api/v1/watchlist/update", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// Test with missing required fields
	missingFields := models.Watchlist{
		Title: "Missing Fields",
		// Missing other required fields
	}
	body, _ := json.Marshal(missingFields)
	req, _ = http.NewRequest("POST", "/api/v1/watchlist/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// Test with non-existent ID
	req, _ = http.NewRequest("GET", "/api/v1/watchlist/999", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
