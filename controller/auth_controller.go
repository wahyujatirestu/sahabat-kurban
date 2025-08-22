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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "Register request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res,
		"message": "User registered successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": err.Error()})
		return
	}

	ctx.JSON(200,gin.H{
		"status": 200,
		"data": res,
		"message": "Login successfully",
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh JWT access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh request"
// @Success 200 {object} dto.TokenOnlyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.authService.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": res})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token and logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Logout request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), req.RefreshToken); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Logout successfully",
	})
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify user email using token sent by email
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/verify-email [get]
func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "token is required"})
		return
	}

	err := c.authService.VerifyEmail(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Email verified successfully",
	})
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resend email verification link
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/resend-verification [post]
func (c *AuthController) ResendVerification(ctx *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	err := c.authService.ResendVerification(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Verification email has been resent successfully"})
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send reset password link to user email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}
	err := c.authService.ForgotPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Email reset password has been sent successfully",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "New password and confirm password must be the same"})
		return
	}

	err := c.authService.ResetPassword(ctx.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Password reset successfully",
	})
}
