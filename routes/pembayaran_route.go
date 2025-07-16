package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func PembayaranRoute(rg *gin.RouterGroup, c *controller.PembayaranController, auth middleware.AuthMiddleware) {
	p := rg.Group("/pembayaran")
	{
		p.POST("/", auth.RequireToken(), c.Create)
		p.GET("/", auth.RequireToken("admin"), c.GetAll)
		p.GET("/:id", auth.RequireToken("admin"), c.GetByID)
		p.GET("/order/:order_id", auth.RequireToken("admin"), c.GetByOrderID)
	}
}
