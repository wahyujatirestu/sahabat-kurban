package dto

type RegisterUserRequest struct {
	Username		string		`json:"username" binding:"required,min=8"`
	Name			string		`json:"name" binding:"required"`
	Email			string		`json:"email" binding:"required,email"`
	Password		string		`json:"password" binding:"required,min=8"`
	ConfirmPassword string 		`json:"confirm_password" binding:"required,eqfield=Password"`
}

type RegisterRequest struct {
	Username		string		`json:"username" binding:"required,min=8"`
	Name			string		`json:"name" binding:"required"`
	Email			string		`json:"email" binding:"required,email"`
	Password		string		`json:"password" binding:"required,min=8"`
	Role			string		`json:"role" binding:"required,oneof=admin panitia user"`
}

type LoginRequest struct {
	Identifier		string		`json:"identifier" binding:"required"`
	Password		string		`json:"Password" binding:"required"`	
}

type AuthResponse struct {
	ID				string		`json:"id"`
	Name			string		`json:"name"`
	Email			string		`json:"email"`
	Role			string		`json:"role"`
	AccessToken		string		`json:"accessToken"`
	RefreshToken	string		`json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenOnlyResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin panitia user"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       	string  `json:"token" binding:"required"`
	NewPassword 	string  `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string 	`json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
