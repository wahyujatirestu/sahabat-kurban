package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PembayaranController struct {
	service service.PembayaranKurbanService
	pekurbanServ service.PekurbanService
}

func NewPembayaranController(s service.PembayaranKurbanService, p service.PekurbanService) *PembayaranController {
	return &PembayaranController{service: s, pekurbanServ: p}
}

func (c *PembayaranController) Create(ctx *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	if currentUser.Role == "user" {
		pekurban, err := c.pekurbanServ.GetByUserId(ctx.Request.Context(), currentUser.ID)
		if err != nil {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You are not registered as a pekurban"})
			return
		}
		if pekurban == nil || pekurban.ID != req.PekurbanID.String() {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You can only pay for your own kurban"})
			return
		}
	} else if currentUser.Role == "panitia" {
		pekurban, err := c.pekurbanServ.GetById(ctx.Request.Context(), req.PekurbanID)
		if err != nil || pekurban == nil {
			ctx.JSON(404, gin.H{
				"status": 404,
				"error": "Pekurban not found"})
			return
		}
		if pekurban.UserID != nil {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "Panitia cannot pay for registered users"})
			return
		}
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res, 
		"message": "Pembayaran successfully created",
	})
}

func (c *PembayaranController) GetAll(ctx *gin.Context) {
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
		"message": "Pembayaran retrieved successfully",
	})
}

func (c *PembayaranController) GetByID(ctx *gin.Context) {
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
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	if res == nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": "Data not found"})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Pembayaran retrieved successfully",
	})
}

func (c *PembayaranController) GetByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	if orderID == "" {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Order ID is required"})
		return
	}

	res, err := c.service.GetByOrderID(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	if res == nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": "Order ID not found"})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Pembayaran retrieved successfully",
	})
}

func (c *PembayaranController) GetRekapDanaPerHewan(ctx *gin.Context) {
	res, err := c.service.GetRekapDanaPerHewan(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Pembayaran retrieved successfully",
	})
}

func (c *PembayaranController) GetProgressPembayaran(ctx *gin.Context) {
	res, err := c.service.GetProgressPembayaran(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Pembayaran retrieved successfully",
	})
}
