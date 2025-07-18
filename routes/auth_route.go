package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
)

func AuthRoute(rg *gin.RouterGroup, ac *controller.AuthController) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", ac.Register)
		authGroup.POST("/login", ac.Login)
		authGroup.POST("/refresh", ac.RefreshToken)
		authGroup.POST("/logout", ac.Logout)
		authGroup.GET("/verify-email", ac.VerifyEmail)
		authGroup.POST("/resend-verification", ac.ResendVerification)
		authGroup.POST("/forgot-password", ac.ForgotPassword)
		authGroup.POST("/reset-password", ac.ResetPassword)
	}
}