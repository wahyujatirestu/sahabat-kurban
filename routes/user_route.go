package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
)

func UserRoute(rg *gin.RouterGroup, uc *controller.UserController, auth middleware.AuthMiddleware)  {
	users := rg.Group("/users")
	{
		users.GET("/", auth.RequireToken("admin"), uc.GetAll)                   
		users.GET("/:id", auth.RequireToken("admin"), uc.GetById)
		users.POST("/register", auth.RequireToken("admin"), uc.CreateAdmin)             
		users.PATCH("/:id/role", auth.RequireToken("admin"), uc.UpdateRole)
		users.PUT("/:id", auth.RequireToken("admin", "user", "panitia"), uc.Update)
		users.PUT("/:id/password", auth.RequireToken(), uc.ChangePassword)    
		users.GET("/me", auth.RequireToken(), uc.GetMyProfile)
		users.DELETE("/:id", auth.RequireToken("admin"), uc.Delete)       
	}
}