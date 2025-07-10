package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func HewanKurbanRoute(rg *gin.RouterGroup, c *controller.HewanKurbanController, auth middleware.AuthMiddleware)  {
	hk := rg.Group("/hewan-kurban")
	{
		hk.POST("/", auth.RequireToken("admin"), c.Create)
		hk.PUT("/:id", auth.RequireToken("admin"), c.Update)
		hk.DELETE("/:id", auth.RequireToken("admin"), c.Delete)
		hk.GET("/", auth.RequireToken(), c.GetAll)
		hk.GET("/:id", auth.RequireToken(), c.GetByID)
	}
}