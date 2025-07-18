package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(400, gin.H{"error": "Confirm password does not match"})
		return
	}

	res, err := c.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"data": res,
		"message": "User registered successfully",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200,gin.H{
		"data": res,
		"message": "Login successfully",
	})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.authService.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": res})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), req.RefreshToken); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Logout successfully"})
}

func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(400, gin.H{"error": "token is required"})
		return
	}

	err := c.authService.VerifyEmail(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Email verified successfully"})
}

func (c *AuthController) ResendVerification(ctx *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.ResendVerification(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Verification email has been resent successfully"})
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.authService.ForgotPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Email reset password telah dikirim"})
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(400, gin.H{"error": "Konfirmasi password tidak cocok"})
		return
	}

	err := c.authService.ResetPassword(ctx.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Password berhasil direset"})
}
