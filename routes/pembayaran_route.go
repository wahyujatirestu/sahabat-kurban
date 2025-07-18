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
		p.GET("/", auth.RequireToken("admin", "panitia"), c.GetAll)
		p.GET("/:id", auth.RequireToken("admin", "panitia"), c.GetByID)
		p.GET("/order/:order_id", auth.RequireToken("admin", "panitia"), c.GetByOrderID)
		p.GET("/rekap/hewan", auth.RequireToken("admin", "panitia"), c.GetRekapDanaPerHewan)
		p.GET("/rekap/pekurban", auth.RequireToken("admin", "panitia"), c.GetProgressPembayaran)
	}
}
