package dto

type RegisterUserRequest struct {
	Username		string		`json:"username" binding:"required,min=8"`
	Name			string		`json:"name" binding:"required"`
	Email			string		`json:"email" binding:"required,email"`
	Password		string		`json:"password" binding:"required,min=8"`
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

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin panitia user"`
}