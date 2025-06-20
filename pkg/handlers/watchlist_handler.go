package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/models"
	"github.com/saketV8/cine-dots/pkg/repositories"
)

type WatchListHandler struct {
	WatchListModel *repositories.WatchListModel
}

// Get
// =====================================================================================
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
