package models

import "time"

// Watchlist represents a single watchlist entry for a movie
type Watchlist struct {
	WatchlistID int       `json:"watchlist_id"`
	Title       string    `json:"title" binding:"required"`
	ReleaseYear int       `json:"release_year" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	Director    string    `json:"director" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	AddedDate   time.Time `json:"added_date" binding:"required"`
}

type WatchListDeleteRequest struct {
	WatchlistID int `json:"watchlist_id" binding:"required"`
}

// type WatchListUpdateRequest struct {
// 	Watchlist

// 	// overriding the WatchlistID from <Watchlist> struct
// 	WatchlistID int       `json:"watchlist_id" binding:"required"`
// 	AddedDate   time.Time `json:"added_date,omitempty"`
// }

type WatchListUpdateRequest struct {
	WatchlistID int       `json:"watchlist_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	ReleaseYear int       `json:"release_year" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	Director    string    `json:"director" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	AddedDate   time.Time `json:"added_date"`
}
