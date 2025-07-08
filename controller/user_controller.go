package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.userService.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, users)
}

func (c *UserController) GetById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.userService.GetById(ctx.Request.Context(), id)
	if err != nil || user == nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(200, user)
}

func (c *UserController) GetMyProfile(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	user, err := c.userService.GetById(ctx.Request.Context(), currentUser.ID)
	if err != nil || user == nil {
		ctx.JSON(404, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(200, user)
}


func (c *UserController) Update(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if currentUser.Role != "admin" && currentUser.ID != id {
		ctx.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.Update(ctx.Request.Context(), id, req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User updated successfully"})
}

func (c *UserController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.UpdateRole(ctx.Request.Context(), id, req.Role); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Role updated successfully"})
}

func (c *UserController) CreateAdmin(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "admin" {
		ctx.JSON(400, gin.H{"error": "Only admin role is allowed here"})
		return
	}

	res, err := c.userService.CreateWithRole(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"data": res,
		"massage": "Admin created successfully",
	})
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	user := userRaw.(model.User)

	paramID := ctx.Param("id")
	if user.ID.String() != paramID {
		ctx.JSON(403, gin.H{"error": "you can only change your own password"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uid, err := uuid.Parse(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.userService.ChangePassword(ctx.Request.Context(), uid, req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Password changed successfully"})
}

func (c *UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.userService.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}

