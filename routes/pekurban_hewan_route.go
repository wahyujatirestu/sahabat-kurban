package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func PekurbanHewanRoute(rg *gin.RouterGroup, c *controller.PekurbanHewanController, authMw middleware.AuthMiddleware) {
	r := rg.Group("/patungan")
	{
		r.POST("/", authMw.RequireToken(), c.Create)
		r.GET("/", authMw.RequireToken("admin", "panitia"), c.GetAll)
		r.GET("/hewan/:hewan_id", authMw.RequireToken(), c.GetByHewanID)
		r.GET("/pekurban/:pekurban_id", authMw.RequireToken(), c.GetByPekurbanID)
		r.PUT("/:id", authMw.RequireToken(), c.Update)
		r.DELETE("/:pekurban_id/:hewan_id", authMw.RequireToken("admin"), c.Delete)
	}
}
