package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func PekurbanRoute(rg *gin.RouterGroup, c *controller.PekurbanController, auth middleware.AuthMiddleware)  {
	p := rg.Group("/pekurban")
	{
		p.GET("/", auth.RequireToken("admin", "panitia"), c.GetAll)
		p.GET("/:id", auth.RequireToken("admin", "panitia"), c.GetById)
		p.POST("/", auth.RequireToken(), c.Create)
		p.PUT("/:id", auth.RequireToken(), c.Update)
		p.DELETE("/:id", auth.RequireToken("admin"), c.Delete)
	}
}