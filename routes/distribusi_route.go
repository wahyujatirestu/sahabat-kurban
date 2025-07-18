package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func DistribusiDagingRoute(rg *gin.RouterGroup, c *controller.DistribusiDagingController, auth middleware.AuthMiddleware)  {
	r := rg.Group("/distribusi")
	{
		r.POST("/", auth.RequireToken("admin", "panitia"), c.Create)
		r.GET("/", auth.RequireToken("admin", "panitia"), c.GetAll)
		r.GET("/:id", auth.RequireToken("admin", "panitia"), c.GetByID)
		r.DELETE("/:id", auth.RequireToken("admin"), c.Delete)
		r.GET("/total-paket", auth.RequireToken("admin", "panitia"), c.GetTotalPaket)
		r.GET("/belum-terdistribusi", auth.RequireToken("admin", "panitia"), c.GetPenerimaBelumDistribusi)
	}
}