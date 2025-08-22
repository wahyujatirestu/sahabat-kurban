package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func RegisterReportRoutes(r *gin.RouterGroup, auth middleware.AuthMiddleware, rc *controller.ReportController) {
	grp := r.Group("/laporan")
	{
		grp.GET("/", auth.RequireToken("admin", "panitia"), rc.GetLaporan)
	}
}
