package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/models"
	"github.com/saketV8/cine-dots/pkg/repositories"
)

type WatchListHandler struct {
	// WatchListModel *repositories.WatchListModel
	WatchListModel repositories.WatchListModelInterface // Interface type
}

// GET Methods
// =====================================================================================

// GetAllWatchListHandler godoc
// @Summary      Get all Watchlists
// @Description  Retrieves all watchlists from the database.
// @Tags         watchlists
// @Produce      json
// @Success      200  {array}  models.Watchlist
// @Failure      500  {object} gin.H  "Failed to get All WatchList"
// @Router       /watchlist/all [get]
func (watchListHandler *WatchListHandler) GetAllWatchListHandler(ctx *gin.Context) {
	watchLists, err := watchListHandler.WatchListModel.GetAllWatchList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get All WatchList",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, watchLists)
}

// GetWatchedListHandler godoc
// @Summary      Retrieve watched watchlists
// @Description  Fetches all watchlists with a "watched" status from the database
// @Tags         watchlists
// @Produce      json
// @Success      200  {array}  models.Watchlist
// @Failure      500  {object} gin.H  "Failed to get Watched List"
// @Router       /watchlist/watched [get]
func (watchListHandler *WatchListHandler) GetWatchedListHandler(ctx *gin.Context) {
	watchLists, err := watchListHandler.WatchListModel.GetWatchedList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Watched List",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, watchLists)
}

// GetWatchingListHandler godoc
// @Summary      Retrieve watchlists with "watching" status
// @Description  Returns all watchlists with a "watching" status from the database
// @Tags         watchlists
// @Produce      json
// @Success      200  {array}   models.Watchlist
// @Failure      500  {object}  gin.H  "Failed to get Watching List"
// @Router       /watchlist/watching [get]
func (watchListHandler *WatchListHandler) GetWatchingListHandler(ctx *gin.Context) {
	watchLists, err := watchListHandler.WatchListModel.GetWatchingList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Watching List",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, watchLists)
}

// GetNotWatchedListHandler godoc
// @Summary      Retrieve watchlists that are not watched
// @Description  Returns all watchlists with a "not watched" status from the database
// @Tags         watchlists
// @Produce      json
// @Success      200  {array}   models.Watchlist
// @Failure      500  {object}  gin.H  "Failed to get Watching List"
// @Router       /watchlist/notwatched [get]
func (watchListHandler *WatchListHandler) GetNotWatchedListHandler(ctx *gin.Context) {
	watchLists, err := watchListHandler.WatchListModel.GetNotWatchedList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get Watching List",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, watchLists)
}

// GetWatchListByIdHandler godoc
// @Summary      Retrieve a watchlist by ID
// @Description  Fetches the watchlist whose ID is provided in the path
// @Tags         watchlists
// @Produce      json
// @Param        watchlist_id  path      string  true  "Watchlist ID"
// @Success      200           {object}  models.Watchlist
// @Failure      500           {object}  gin.H  "Failed to get WatchList by ID"
// @Router       /watchlist/{watchlist_id} [get]
func (watchListHandler *WatchListHandler) GetWatchListByIdHandler(ctx *gin.Context) {
	watchlist_id_param := ctx.Param("watchlist_id")
	watchLists, err := watchListHandler.WatchListModel.GetWatchListById(watchlist_id_param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get WatchList by ID",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, watchLists)
}

// =====================================================================================
// =====================================================================================

// Other methods
// =====================================================================================

// AddWatchListHandler godoc
// @Summary      Create a new watchlist item
// @Description  Adds a new watchlist entry to the database
// @Tags         watchlists
// @Accept       json
// @Produce      json
// @Param        watchlist  body      models.WatchListAddRequestExample  true  "Watchlist Data"
// @Success      200        {object}  models.Watchlist
// @Failure      400        {object}  gin.H  "Invalid WatchList Data"
// @Failure      500        {object}  gin.H  "Failed to add WatchList data"
// @Router       /watchlist/add [post]
func (watchListHandler *WatchListHandler) AddWatchListHandler(ctx *gin.Context) {
	//getting param from POST request body
	var body models.Watchlist

	// this will bind data coming from POST request
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid WatchList Data",
			"details": err.Error(),
		})
		return
	}

	watchListAdded, err := watchListHandler.WatchListModel.AddWatchList(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add WatchList data",
			"details": err.Error(),
			"body":    body,
		})
		return
	}

	ctx.JSON(http.StatusOK, watchListAdded)
}

// DeleteWatchListHandler godoc
// @Summary      Delete a watchlist entry
// @Description  Removes a watchlist from the database based on the provided watchlist ID
// @Tags         watchlists
// @Accept       json
// @Produce      json
// @Param        request  body      models.WatchListDeleteRequest  true  "Delete Request (watchlist_id)"
// @Success      200      {object}  gin.H  "WatchList deleted successfully"
// @Failure      400      {object}  gin.H  "Invalid WatchList ID"
// @Failure      500      {object}  gin.H  "Failed to delete WatchList"
// @Router       /watchlist/delete [delete]
func (watchListHandler *WatchListHandler) DeleteWatchListHandler(ctx *gin.Context) {
	//getting param from POST request body
	var body models.WatchListDeleteRequest

	// this will bind data coming from POST request
	err := ctx.BindJSON(&body)
	if err != nil {
		// If binding fails, return a 400 error with the error message
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid WatchList ID",
			"details": err.Error(),
		})
		return
	}

	rowAffected, err := watchListHandler.WatchListModel.DeleteWatchList(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete WatchList",
			"details": err.Error(),
			"body":    body,
		})
		return
	}

	// Add option to check if rowAffected == 0 then return Already Deleted or DATA DNE
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "WatchList deleted successfully",
		"row-affected": rowAffected,
		"body":         body,
	})
}

// UpdateWatchListHandler godoc
// @Summary      Update an existing watchlist entry
// @Description  Updates an existing watchlist with new data
// @Tags         watchlists
// @Accept       json
// @Produce      json
// @Param        request  body      models.WatchListUpdateRequestExample  true  "Updated WatchList Data"
// @Success      200      {object}  gin.H  "WatchList updated successfully"
// @Failure      400      {object}  gin.H  "Invalid WatchList Data"
// @Failure      500      {object}  gin.H  "Failed to update WatchList"
// @Router       /watchlist/update [patch]
func (watchListHandler *WatchListHandler) UpdateWatchListHandler(ctx *gin.Context) {
	//getting param from POST request body
	var body models.WatchListUpdateRequest

	// this will bind data coming from POST request
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid WatchList Data",
			"details": err.Error(),
		})
		return
	}

	rowAffected, err := watchListHandler.WatchListModel.UpdateWatchList(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update WatchList",
			"details": err.Error(),
			"body":    body,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "WatchList updated successfully",
		"row-affected": rowAffected,
		"body":         body,
	})
}
