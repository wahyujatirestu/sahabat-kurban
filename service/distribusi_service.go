package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
)

type DistribusiDagingService interface {
	Create(ctx context.Context, req dto.CreateDistribusiRequest) (*dto.DistribusiResponse, error)
	GetAll(ctx context.Context) ([]dto.DistribusiResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.DistribusiResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalDistribusiPaket(ctx context.Context) (int, error)
	GetPenerimaBelumTerdistribusi(ctx context.Context) ([]dto.PenerimaResponse, error)
}

type distribusiDagingService struct {
	repo repository.DistribusiDagingRepository
	penerimaRepo repository.PenerimaDagingRepository
	hewanRepo repository.HewanKurbanRepository
}

func NewDistribusiDagingService(repo repository.DistribusiDagingRepository, penerimaRepo repository.PenerimaDagingRepository, hewanRepo repository.HewanKurbanRepository) DistribusiDagingService {
	return &distribusiDagingService{repo: repo, penerimaRepo: penerimaRepo, hewanRepo: hewanRepo}
}

func (s *distribusiDagingService) Create(ctx context.Context, req dto.CreateDistribusiRequest) (*dto.DistribusiResponse, error) {
	penerimaID, err := uuid.Parse(req.PenerimaID)
	if err != nil {
		return nil, errors.New("Invalid Penerima ID")
	}
	hewanID, err := uuid.Parse(req.HewanID)
	if err != nil {
		return nil, errors.New("Invalid Hewan ID")
	}

	tanggalDistribusi, err := time.Parse("2006-01-02", req.TanggalDistribusi)
	if err != nil {
		return nil, errors.New("Invalid date format. Use YYYY-MM-DD")
	}

	penerima, err := s.penerimaRepo.GetByID(ctx, penerimaID)
	if err != nil {
		return nil, err
	}
	if penerima == nil {
		return nil, errors.New("Penerima not found")
	}

	hewan, err := s.hewanRepo.GetById(ctx, hewanID)
	if err != nil {
		return nil, err
	}
	if hewan == nil {
		return nil, errors.New("Hewan not found")
	}

	existing, err := s.repo.FindByPenerimaID(ctx, penerimaID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("Penerima already received distribution")
	}

	dis := &model.DistribusiDaging{
		ID: uuid.New(),
		PenerimaID: penerimaID,
		HewanID: hewanID,
		JumlahPaket: req.JumlahPaket,
		TanggalDistribusi: tanggalDistribusi,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.repo.Create(ctx, dis); err != nil {
		return nil, err
	}

	res := dto.ToDistribusiResponse(dis)
	return &res, nil
}

func (s *distribusiDagingService) GetAll(ctx context.Context) ([]dto.DistribusiResponse, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.DistribusiResponse
	for _, d := range list{
		res = append(res, dto.ToDistribusiResponse(d))
	}

	return res, nil
}

func (s *distribusiDagingService) GetByID(ctx context.Context, id uuid.UUID) (*dto.DistribusiResponse, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil || d == nil {
		return nil, err
	}

	res := dto.ToDistribusiResponse(d)
	return &res, nil
}

func (s *distribusiDagingService) GetTotalDistribusiPaket(ctx context.Context) (int, error) {
	return s.repo.CountTotalPaket(ctx)
}

func (s *distribusiDagingService) GetPenerimaBelumTerdistribusi(ctx context.Context) ([]dto.PenerimaResponse, error) {
	list, err := s.penerimaRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.PenerimaResponse
	for _, penerima := range list {
		exist, err := s.repo.FindByPenerimaID(ctx, penerima.ID)
		if err != nil {
			return nil, err
		}
		if exist == nil {
			result = append(result, dto.ToPenerimaResponse(penerima))
		}
	}

	return result, nil
}


func (s *distribusiDagingService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}