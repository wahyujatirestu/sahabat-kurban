package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PekurbanController struct {
	pService service.PekurbanService
}

func NewPekurbanController(pService service.PekurbanService) *PekurbanController {
	return &PekurbanController{pService: pService}
}

func (c *PekurbanController) Create(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	var req dto.CreatePekurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if currentUser.Role == "user" {
		if req.UserID == nil {
			ctx.JSON(403, gin.H{"error": "You must include your own user_id"})
			return
		}
		uid, err := uuid.Parse(*req.UserID)
		if err != nil || uid != currentUser.ID {
			ctx.JSON(403, gin.H{"error": "You can only register yourself as pekurban"})
			return
		}
	}


	if currentUser.Role == "user" {
		existing, _ := c.pService.GetByUserId(ctx.Request.Context(), currentUser.ID)
		if existing != nil {
			ctx.JSON(409, gin.H{"error": "You have already registered as pekurban"})
			return
		}
	}

	res, err := c.pService.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"data": res,
		"message": "Pekurban created successfully",
	})
}

func (c *PekurbanController) GetAll(ctx *gin.Context) {
	data, err := c.pService.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": data})
}

func (c *PekurbanController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	data, err := c.pService.GetById(ctx.Request.Context(), id)
	if err != nil || data == nil {
		ctx.JSON(404, gin.H{"error": "Pekurban not found"})
		return
	}

	ctx.JSON(200, gin.H{"data": data})
}

func (c *PekurbanController) Update(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
	}

	p, err := c.pService.GetById(ctx.Request.Context(), id)
	if err != nil || p == nil {
		ctx.JSON(404, gin.H{"error": "Pekurban not found"})
		return
	}

	if currentUser.Role != "admin" {
		if currentUser.Role == "user" {
			if p.UserID == nil || *p.UserID != currentUser.ID.String() {
				ctx.JSON(403, gin.H{"error": "You can only update your own kurban data"})
				return
			}
		}
		if currentUser.Role == "panitia" {
			if !(p.UserID == nil || *p.UserID == currentUser.ID.String()) {
				ctx.JSON(403, gin.H{"error": "Panitia can only update offline data or their own"})
			}
		}
	}

	var req dto.UpdatePekurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.pService.Update(ctx.Request.Context(), id, req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Pekurban updated successfully"})
}

func (c *PekurbanController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.pService.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Pekurban deleted successfully"})
}