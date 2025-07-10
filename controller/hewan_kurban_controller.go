package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type HewanKurbanController struct {
	service service.HewanKurbanService
}

func NewHewanKurbanController(service service.HewanKurbanService) *HewanKurbanController {
	return &HewanKurbanController{service: service}
}

func (c *HewanKurbanController) Create(ctx *gin.Context) {
	var req dto.CreateHewanKurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"data": res,
		"message": "Hewan kurban created successfully",
	})
}

func (c *HewanKurbanController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	data, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": data})
}

func (c *HewanKurbanController) GetAll(ctx *gin.Context) {
	list, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": list})
}

func (c *HewanKurbanController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	var req dto.UpdateHewanKurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(ctx.Request.Context(), id, req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Hewan kuban updated successfully"})
}

func (c *HewanKurbanController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Hewan kurban deleted successfully"})
}