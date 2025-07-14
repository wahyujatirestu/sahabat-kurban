package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PenerimaDagingController struct {
	service service.PenerimaDagingService
}

func NewPenerimaDagingController(service service.PenerimaDagingService) *PenerimaDagingController {
	return &PenerimaDagingController{service: service}
}

func (c *PenerimaDagingController) Create(ctx *gin.Context) {
	var req dto.CreatePenerimaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"data": res,
		"message": "Penerima daging created successfully",
	})
}

func (c *PenerimaDagingController) GetAll(ctx *gin.Context) {
	res, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": res})
}

func (c *PenerimaDagingController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": res})
}

func (c *PenerimaDagingController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	var req dto.UpdatePenerimaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(ctx.Request.Context(), id, req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Penerima daging updated successfully"})
}

func (c *PenerimaDagingController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Penerima daging deleted successfully"})
}