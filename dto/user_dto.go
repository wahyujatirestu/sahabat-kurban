package dto

import "github.com/wahyujatirestu/sahabat-kurban/model"

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type ChangePasswordRequest struct {
	OldPassword 	string `json:"old_password" binding:"required"`
	NewPassword 	string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}



func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}
}

func ToUserResponseList(users []*model.User) []UserResponse {
	var result []UserResponse
	for _, u := range users {
		result = append(result, ToUserResponse(u))
	}
	return result
}
