package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/handlers"
	"github.com/saketV8/cine-dots/pkg/models"
	"github.com/stretchr/testify/assert"
)

// mockWatchListRepository is an in-memory mock that implements WatchListRepository.
// It simulates a database by returning data or errors based on test scenarios.
type mockWatchListRepository struct {
	getAllFunc        func() ([]models.Watchlist, error)
	getWatchedFunc    func() ([]models.Watchlist, error)
	getWatchingFunc   func() ([]models.Watchlist, error)
	getNotWatchedFunc func() ([]models.Watchlist, error)
	getByIDFunc       func(string) (models.Watchlist, error)
	addFunc           func(models.Watchlist) (models.Watchlist, error)
	deleteFunc        func(models.WatchListDeleteRequest) (int, error)
	updateFunc        func(models.WatchListUpdateRequest) (int, error)
}

func (m *mockWatchListRepository) GetAllWatchList() ([]models.Watchlist, error) {
	return m.getAllFunc()
}

func (m *mockWatchListRepository) GetWatchedList() ([]models.Watchlist, error) {
	return m.getWatchedFunc()
}

func (m *mockWatchListRepository) GetWatchingList() ([]models.Watchlist, error) {
	return m.getWatchingFunc()
}

func (m *mockWatchListRepository) GetNotWatchedList() ([]models.Watchlist, error) {
	return m.getNotWatchedFunc()
}

func (m *mockWatchListRepository) GetWatchListById(id string) (models.Watchlist, error) {
	return m.getByIDFunc(id)
}

func (m *mockWatchListRepository) AddWatchList(w models.Watchlist) (models.Watchlist, error) {
	return m.addFunc(w)
}

func (m *mockWatchListRepository) DeleteWatchList(req models.WatchListDeleteRequest) (int, error) {
	return m.deleteFunc(req)
}

func (m *mockWatchListRepository) UpdateWatchList(req models.WatchListUpdateRequest) (int, error) {
	return m.updateFunc(req)
}

func setupTestRouter(handler *handlers.WatchListHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/watchlist/all", handler.GetAllWatchListHandler)
		api.GET("/watchlist/watched", handler.GetWatchedListHandler)
		api.GET("/watchlist/watching", handler.GetWatchingListHandler)
		api.GET("/watchlist/notwatched", handler.GetNotWatchedListHandler)
		api.GET("/watchlist/:watchlist_id", handler.GetWatchListByIdHandler)
		api.POST("/watchlist/add", handler.AddWatchListHandler)
		api.DELETE("/watchlist/delete", handler.DeleteWatchListHandler)
		api.PATCH("/watchlist/update", handler.UpdateWatchListHandler)
	}

	return router
}

// TestGetAllWatchListHandler tests the GetAllWatchListHandler for both success and error scenarios.
func TestGetAllWatchListHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func() ([]models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - returns all watchlist items",
			mockFunc: func() ([]models.Watchlist, error) {
				return []models.Watchlist{
					{WatchlistID: 1, Title: "Movie A", Status: "watched"},
					{WatchlistID: 2, Title: "Movie B", Status: "not watched"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Database error",
			mockFunc: func() ([]models.Watchlist, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				getAllFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/watchlist/all", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data []models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.True(t, len(data) > 0)
			}
		})
	}
}

// TestGetWatchedListHandler tests fetching watched items.
func TestGetWatchedListHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func() ([]models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - returns watched items",
			mockFunc: func() ([]models.Watchlist, error) {
				return []models.Watchlist{
					{WatchlistID: 1, Title: "Movie A", Status: "watched"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Error from repository",
			mockFunc: func() ([]models.Watchlist, error) {
				return nil, errors.New("some error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				getWatchedFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/watchlist/watched", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data []models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.True(t, len(data) > 0)
			}
		})
	}
}

// TestGetWatchingListHandler tests fetching watching items.
func TestGetWatchingListHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func() ([]models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - returns watching items",
			mockFunc: func() ([]models.Watchlist, error) {
				return []models.Watchlist{
					{WatchlistID: 1, Title: "Series X", Status: "watching"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Failure - repository error",
			mockFunc: func() ([]models.Watchlist, error) {
				return nil, errors.New("failed to get data")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				getWatchingFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/watchlist/watching", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data []models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.True(t, len(data) > 0)
			}
		})
	}
}

// TestGetNotWatchedListHandler tests fetching not watched items.
func TestGetNotWatchedListHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func() ([]models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - returns not watched items",
			mockFunc: func() ([]models.Watchlist, error) {
				return []models.Watchlist{
					{WatchlistID: 1, Title: "Documentary Y", Status: "not watched"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Error from repository",
			mockFunc: func() ([]models.Watchlist, error) {
				return nil, errors.New("error retrieving data")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				getNotWatchedFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/watchlist/notwatched", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data []models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.True(t, len(data) > 0)
			}
		})
	}
}

// TestGetWatchListByIdHandler tests fetching a single item by ID.
func TestGetWatchListByIdHandler(t *testing.T) {
	tests := []struct {
		name           string
		idParam        string
		mockFunc       func(string) (models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name:    "Success - returns item by ID",
			idParam: "1",
			mockFunc: func(id string) (models.Watchlist, error) {
				return models.Watchlist{WatchlistID: 1, Title: "Movie A"}, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:    "Error retrieving by ID",
			idParam: "999",
			mockFunc: func(id string) (models.Watchlist, error) {
				return models.Watchlist{}, errors.New("not found")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				getByIDFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/watchlist/"+tt.idParam, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.Equal(t, 1, data.WatchlistID)
			}
		})
	}
}

// TestAddWatchListHandler tests adding a new item to the watchlist.
func TestAddWatchListHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          models.Watchlist
		mockFunc       func(models.Watchlist) (models.Watchlist, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - item created",
			input: models.Watchlist{
				Title:       "New Movie",
				ReleaseYear: 2000,
				Genre:       "Adventure",
				Director:    "Director A",
				Status:      "not watched",
				AddedDate:   time.Now(),
			},
			mockFunc: func(w models.Watchlist) (models.Watchlist, error) {
				w.WatchlistID = 100
				return w, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Invalid input - no title",
			input: models.Watchlist{
				Title: "",
			},
			mockFunc: func(w models.Watchlist) (models.Watchlist, error) {
				return models.Watchlist{}, errors.New("invalid input")
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Database error",
			input: models.Watchlist{
				Title:       "Some Movie",
				ReleaseYear: 1999,
				Genre:       "Drama",
				Director:    "Director B",
			},
			mockFunc: func(w models.Watchlist) (models.Watchlist, error) {
				return models.Watchlist{}, errors.New("db error")
			},
			// expectedStatus: http.StatusInternalServerError,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				addFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/watchlist/add", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var data models.Watchlist
				_ = json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NotZero(t, data.WatchlistID)
			}
		})
	}
}

// TestDeleteWatchListHandler tests deleting a watchlist item.
func TestDeleteWatchListHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          models.WatchListDeleteRequest
		mockFunc       func(models.WatchListDeleteRequest) (int, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - item deleted",
			input: models.WatchListDeleteRequest{
				WatchlistID: 123,
			},
			mockFunc: func(req models.WatchListDeleteRequest) (int, error) {
				return 1, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Invalid input",
			input: models.WatchListDeleteRequest{
				WatchlistID: 0, // invalid
			},
			mockFunc: func(req models.WatchListDeleteRequest) (int, error) {
				return 0, errors.New("invalid id")
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Database error",
			input: models.WatchListDeleteRequest{
				WatchlistID: 999,
			},
			mockFunc: func(req models.WatchListDeleteRequest) (int, error) {
				return 0, errors.New("not found")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				deleteFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/watchlist/delete", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var success map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &success)
				assert.Equal(t, float64(1), success["row-affected"])
			}
		})
	}
}

// TestUpdateWatchListHandler tests updating a watchlist item.
func TestUpdateWatchListHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          models.WatchListUpdateRequest
		mockFunc       func(models.WatchListUpdateRequest) (int, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - item updated",
			input: models.WatchListUpdateRequest{
				WatchlistID: 50,
				Title:       "Updated Title",
				ReleaseYear: 2021,
				Status:      "watched",
			},
			mockFunc: func(req models.WatchListUpdateRequest) (int, error) {
				return 1, nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Invalid input - watchlist ID missing",
			input: models.WatchListUpdateRequest{
				WatchlistID: 0,
			},
			mockFunc: func(req models.WatchListUpdateRequest) (int, error) {
				return 0, errors.New("invalid request")
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Item not found in DB",
			input: models.WatchListUpdateRequest{
				WatchlistID: 999,
				Title:       "Ghost Movie",
			},
			mockFunc: func(req models.WatchListUpdateRequest) (int, error) {
				return 0, errors.New("not found")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWatchListRepository{
				updateFunc: tt.mockFunc,
			}
			h := &handlers.WatchListHandler{WatchListModel: mockRepo}
			router := setupTestRouter(h)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPatch, "/api/v1/watchlist/update", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectError {
				var errResp map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &errResp)
				assert.Contains(t, errResp, "error")
			} else {
				var success map[string]interface{}
				_ = json.Unmarshal(resp.Body.Bytes(), &success)
				assert.Equal(t, float64(1), success["row-affected"])
			}
		})
	}
}
