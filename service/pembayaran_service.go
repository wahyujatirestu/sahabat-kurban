// service/pembayaran_service.go
package service

import (
	"context"
	"errors"
	"math"
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
	GetRekapDanaPerHewan(ctx context.Context) ([]dto.RekapDanaHewanResponse, error)
	GetProgressPembayaran(ctx context.Context) ([]dto.ProgressPembayaranPekurban, error)
}

type pembayaranKurbanService struct {
	repo            repository.PembayaranKurbanRepository
	midtransService payserv.MidtransService
	pRepo 			repository.PekurbanHewanRepository
	hRepo			repository.HewanKurbanRepository
	pekurbanRepo	repository.PekurbanRepository
}

func NewPembayaranKurbanService(repo repository.PembayaranKurbanRepository, mid payserv.MidtransService, pRepo repository.PekurbanHewanRepository, hRepo repository.HewanKurbanRepository, pekurbanRepo repository.PekurbanRepository) PembayaranKurbanService {
	return &pembayaranKurbanService{
		repo: repo,
		midtransService: mid,
		pRepo: pRepo,
		hRepo: hRepo,
		pekurbanRepo: pekurbanRepo,
	}
}

func (s *pembayaranKurbanService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	orderID := utils.GenerateOrderID()

	pekurban, err := s.pekurbanRepo.FindById(ctx, req.PekurbanID)
	if err != nil {
		return nil, err
	}

	if pekurban.Name == nil || pekurban.Email == nil || pekurban.Phone == nil {
		return nil, errors.New("data pekurban tidak lengkap (nama, email, atau phone kosong)")
	}

	patunganList, err := s.pRepo.GetByPekurbanId(ctx, req.PekurbanID)
	if err != nil {
		return nil, err
	}
	if len(patunganList) == 0 {
		return nil, errors.New("pekurban tidak memiliki relasi dengan hewan kurban")
	}

	var total float64
	for _, r := range patunganList {
		hewanId, _ := uuid.Parse(r.HewanID)
		hewan, err := s.hRepo.GetById(ctx, hewanId)
		if err != nil || hewan == nil {
			return nil, errors.New("data hewan kurban not found")
		}
		total += r.Porsi * hewan.Harga
		total = math.Round(total*100) / 100
	}

	payload := dto.ToMidtransChargeRequest(orderID, total, *pekurban.Name, *pekurban.Email, *pekurban.Phone, req)
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
		ID:               	uuid.New(),
		OrderID:          	orderID,
		TransactionID:    	midResp.TransactionID,
		PekurbanID:       	req.PekurbanID,
		Metode:           	req.Metode,
		PaymentType:      	&midResp.PaymentType,
		VANumber:         	vaNumber,
		Status:           	midResp.TransactionStatus,
		FraudStatus:      	midResp.FraudStatus,
		ApprovalCode:     	midResp.ApprovalCode,
		TransactionTime:  	&trxTime,
		Jumlah: 			total,
		TanggalPembayaran: 	time.Now(),
		Created_At:        	time.Now(),
		Updated_At:        	time.Now(),
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	res := dto.ToPaymentResponse(payment, total, midResp)
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
	res := dto.ToPaymentResponse(p, p.Jumlah, nil)
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
	res := dto.ToPaymentResponse(p, p.Jumlah, nil)
	return &res, nil
}

func (s *pembayaranKurbanService) GetAll(ctx context.Context) ([]dto.PaymentResponse, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []dto.PaymentResponse
	for _, p := range list {
		result = append(result, dto.ToPaymentResponse(p, p.Jumlah, nil))
	}
	return result, nil
}

func (s *pembayaranKurbanService) GetRekapDanaPerHewan(ctx context.Context) ([]dto.RekapDanaHewanResponse, error) {
	list, err := s.repo.GetTotalPembayaranPerHewan(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.RekapDanaHewanResponse
	for _, h := range list {
		status := "belum lunas"
		if h.TotalMasuk >= h.HargaTarget {
			if h.TotalMasuk > h.HargaTarget {
				status = "melebihi target"
			} else {
				status = "lunas"
			}
		}
		result = append(result, dto.RekapDanaHewanResponse{
			HewanID:     h.HewanID,
			Jenis:       h.Jenis,
			HargaTarget: h.HargaTarget,
			TotalMasuk:  h.TotalMasuk,
			Status:      status,
		})
	}
	return result, nil
}

func (s *pembayaranKurbanService) GetProgressPembayaran(ctx context.Context) ([]dto.ProgressPembayaranPekurban, error) {
	list, err := s.repo.GetProgressPembayaranPekurban(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.ProgressPembayaranPekurban
	for _, d := range list {
		progress := 0.0
		status := "belum bayar"

		if d.TotalTagihan > 0 {
			progress = d.TotalBayar / d.TotalTagihan
			switch {
			case d.TotalBayar == 0:
				status = "belum bayar"
			case d.TotalBayar < d.TotalTagihan:
				status = "sebagian"
			case d.TotalBayar == d.TotalTagihan:
				status = "lunas"
			case d.TotalBayar > d.TotalTagihan:
				status = "lebih"
			}
		}

		result = append(result, dto.ProgressPembayaranPekurban{
			PekurbanID:    d.PekurbanID,
			NamaPekurban:  d.NamaPekurban,
			JumlahPorsi:   d.PorsiTotal,
			TotalTagihan:  d.TotalTagihan,
			TotalBayar:    d.TotalBayar,
			Progress:      progress,
			Status:        status,
		})
	}
	return result, nil
}
