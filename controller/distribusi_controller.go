package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type DistribusiDagingController struct {
	service service.DistribusiDagingService
}

func NewDistribusiDagingController(service service.DistribusiDagingService) *DistribusiDagingController {
	return &DistribusiDagingController{service: service}
}

func (c *DistribusiDagingController) Create(ctx *gin.Context) {
	var req dto.CreateDistribusiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(422, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"data": res,
		"message": "Distribusi daging created successfully",
	})
}

func (c *DistribusiDagingController) GetAll(ctx *gin.Context) {
	res, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": res})
}

func (c *DistribusiDagingController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": res})
}

func (c *DistribusiDagingController) Delete(ctx *gin.Context) {
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

	ctx.JSON(200, gin.H{"message": "Distribusi daging deleted successfully"})
}