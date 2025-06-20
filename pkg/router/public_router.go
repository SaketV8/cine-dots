// https://stackoverflow.com/a/62608670
// https://stackoverflow.com/questions/62608429/how-to-combine-group-of-routes-in-gin

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/utils"
)

func SetupPublicRouter(app *App, superRouterGroup *gin.Engine) {

	// <ROUTER_PREFIX> = /api
	// <ROUTER_PREFIX_VERSION> = /v1
	routerGroup := superRouterGroup.Group(utils.ROUTER_PREFIX)
	{
		v1 := routerGroup.Group(utils.ROUTER_PREFIX_VERSION)
		{
			v1.GET("/watchlist/all", app.WatchList.GetAllWatchListHandler)
			v1.GET("/watchlist/watched", app.WatchList.GetWatchedListHandler)
			v1.GET("/watchlist/watching", app.WatchList.GetWatchingListHandler)
			v1.GET("/watchlist/notwatched", app.WatchList.GetNotWatchedListHandler)
			v1.GET("/watchlist/:watchlist_id", app.WatchList.GetWatchListByIdHandler)

			v1.POST("/watchlist/add", app.WatchList.AddWatchListHandler)
			v1.DELETE("/watchlist/delete", app.WatchList.DeleteWatchListHandler)
			v1.PATCH("/watchlist/update", app.WatchList.UpdateWatchListHandler)
		}
	}
}
