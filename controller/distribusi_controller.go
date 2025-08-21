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
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(422, gin.H{
			"status": 422,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res,
		"message": "Distribusi daging created successfully",
	})
}

func (c *DistribusiDagingController) GetAll(ctx *gin.Context) {
	res, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Distribusi daging retrieved successfully",
	})
}

func (c *DistribusiDagingController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Distribusi daging retrieved successfully",
	})
}

func (c *DistribusiDagingController) GetTotalPaket(ctx *gin.Context) {
	total, err := c.service.GetTotalDistribusiPaket(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"total_distribusi_paket": total,
		"message": "Distribusi daging retrieved successfully",
	})
}

func (c *DistribusiDagingController) GetPenerimaBelumDistribusi(ctx *gin.Context) {
	data, err := c.service.GetPenerimaBelumTerdistribusi(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Distribusi daging retrieved successfully",
	})
}


func (c *DistribusiDagingController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Distribusi daging deleted successfully",
	})
}