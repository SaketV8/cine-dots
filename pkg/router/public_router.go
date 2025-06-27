// https://stackoverflow.com/a/62608670
// https://stackoverflow.com/questions/62608429/how-to-combine-group-of-routes-in-gin

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/saketV8/cine-dots/pkg/utils"

	docs "github.com/saketV8/cine-dots/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupPublicRouter(app *App, superRouterGroup *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	routerGroup := superRouterGroup.Group(utils.ROUTER_PREFIX)

	{
		v1 := routerGroup.Group(utils.ROUTER_PREFIX_VERSION)
		{
			v1.GET("/watchlist/all", app.WatchListHandler.GetAllWatchListHandler)
			v1.GET("/watchlist/watched", app.WatchListHandler.GetWatchedListHandler)
			v1.GET("/watchlist/watching", app.WatchListHandler.GetWatchingListHandler)
			v1.GET("/watchlist/notwatched", app.WatchListHandler.GetNotWatchedListHandler)
			v1.GET("/watchlist/:watchlist_id", app.WatchListHandler.GetWatchListByIdHandler)

			v1.POST("/watchlist/add", app.WatchListHandler.AddWatchListHandler)
			v1.DELETE("/watchlist/delete", app.WatchListHandler.DeleteWatchListHandler)
			v1.PATCH("/watchlist/update", app.WatchListHandler.UpdateWatchListHandler)
		}
	}

	// routerGroup.GET()
	// Swagger endpoint
	superRouterGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// routerGroup.GET("/health", HealthCheck)
}
