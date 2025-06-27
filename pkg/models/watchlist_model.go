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
	WatchlistID int        `json:"watchlist_id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	ReleaseYear int        `json:"release_year" binding:"required"`
	Genre       string     `json:"genre" binding:"required"`
	Director    string     `json:"director" binding:"required"`
	Status      string     `json:"status" binding:"required"`
	AddedDate   *time.Time `json:"added_date"`
}

// Example for swagger :)

type WatchListAddRequestExample struct {
	Title       string     `json:"title" example:"Coco" binding:"required"`
	ReleaseYear int        `json:"release_year" example:"2017" binding:"required"`
	Genre       string     `json:"genre" example:"Animation" binding:"required"`
	Director    string     `json:"director" example:"Lee Unkrich" binding:"required"`
	Status      string     `json:"status" example:"not watched" binding:"required"`
	AddedDate   *time.Time `json:"added_date" example:"2025-06-20T00:00:00Z"`
}

type WatchListUpdateRequestExample struct {
	WatchlistID int        `json:"watchlist_id" binding:"required" example:"7"`
	Title       string     `json:"title" binding:"required" example:"Coco"`
	ReleaseYear int        `json:"release_year" binding:"required" example:"2017"`
	Genre       string     `json:"genre" binding:"required" example:"Animation"`
	Director    string     `json:"director" binding:"required" example:"Lee Unkrich"`
	Status      string     `json:"status" binding:"required" example:"watching"`
	AddedDate   *time.Time `json:"added_date,omitempty" example:"2025-06-20T00:00:00Z"`
}
