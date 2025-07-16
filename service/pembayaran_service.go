// service/pembayaran_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	payserv "github.com/wahyujatirestu/sahabat-kurban/payments/service"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	"github.com/wahyujatirestu/sahabat-kurban/utils"
)

type PembayaranKurbanService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error)
	GetByOrderID(ctx context.Context, orderID string) (*dto.PaymentResponse, error)
	GetAll(ctx context.Context) ([]dto.PaymentResponse, error)
}

type pembayaranKurbanService struct {
	repo            repository.PembayaranKurbanRepository
	midtransService payserv.MidtransService
}

func NewPembayaranKurbanService(repo repository.PembayaranKurbanRepository, mid payserv.MidtransService) PembayaranKurbanService {
	return &pembayaranKurbanService{
		repo: repo,
		midtransService: mid,
	}
}

func (s *pembayaranKurbanService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	orderID := utils.GenerateOrderID()

	payload := dto.ToMidtransChargeRequest(orderID, req)
	midResp, err := s.midtransService.Charge(payload)
	if err != nil {
		return nil, err
	}

	var vaNumber *string
	if len(midResp.VANumbers) > 0 {
		vaNumber = &midResp.VANumbers[0].VANumber
	}

	trxTime, _ := time.Parse("2006-01-02 15:04:05", midResp.TransactionTime)

	payment := &model.PembayaranKurban{
		ID:               uuid.New(),
		OrderID:          orderID,
		TransactionID:    midResp.TransactionID,
		PekurbanID:       req.PekurbanID,
		Metode:           req.Metode,
		PaymentType:      &midResp.PaymentType,
		VANumber:         vaNumber,
		Jumlah:           req.Jumlah,
		Status:           midResp.TransactionStatus,
		FraudStatus:      midResp.FraudStatus,
		ApprovalCode:     midResp.ApprovalCode,
		TransactionTime:  &trxTime,
		TanggalPembayaran: time.Now(),
		Created_At:        time.Now(),
		Updated_At:        time.Now(),
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	res := dto.ToPaymentResponse(payment, midResp)
	return &res, nil
}

func (s *pembayaranKurbanService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("pembayaran not found")
	}
	res := dto.ToPaymentResponse(p, nil)
	return &res, nil
}

func (s *pembayaranKurbanService) GetByOrderID(ctx context.Context, orderID string) (*dto.PaymentResponse, error) {
	p, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("order ID not found")
	}
	res := dto.ToPaymentResponse(p, nil)
	return &res, nil
}

func (s *pembayaranKurbanService) GetAll(ctx context.Context) ([]dto.PaymentResponse, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []dto.PaymentResponse
	for _, p := range list {
		result = append(result, dto.ToPaymentResponse(p, nil))
	}
	return result, nil
}
