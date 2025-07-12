package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func PenyembelihanRoute(rg *gin.RouterGroup, c *controller.PenyembelihanController, auth middleware.AuthMiddleware) {
	pr := rg.Group("/penyembelihan")
	{
		pr.POST("", auth.RequireToken("admin", "panitia"), c.Create)
		pr.PUT("/:id", auth.RequireToken("admin", "panitia"), c.Update)
		pr.DELETE("/:id", auth.RequireToken("admin"), c.Delete)
		pr.GET("/", auth.RequireToken("admin", "panitia", "user"), c.GetAll)
		pr.GET("/:id", auth.RequireToken("admin", "panitia", "user"), c.GetById)
	}
}
