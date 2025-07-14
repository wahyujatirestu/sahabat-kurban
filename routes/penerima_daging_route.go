package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func PenerimaDagingRoute(rg *gin.RouterGroup, c *controller.PenerimaDagingController, auth middleware.AuthMiddleware) {
	r := rg.Group("/penerima")
	{
		r.POST("/", auth.RequireToken("admin", "panitia"), c.Create)
		r.PUT("/:id", auth.RequireToken("admin", "panitia"), c.Update)
		r.DELETE("/:id", auth.RequireToken("admin"), c.Delete)
		r.GET("/", auth.RequireToken(), c.GetAll)
		r.GET("/:id", auth.RequireToken(), c.GetByID)
	}
}
