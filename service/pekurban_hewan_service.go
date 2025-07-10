package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
)

type PekurbanHewanService interface {
	Create(ctx context.Context, req dto.CreatePekurbanHewanRequest) error
	Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error
	GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
	GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
}

type pekurbanHewanService struct {
	repo repository.PekurbanHewanRepository
}

func NewPekurbanHewanService(repo repository.PekurbanHewanRepository) PekurbanHewanService {
	return &pekurbanHewanService{repo: repo}
}

func (s *pekurbanHewanService) Create(ctx context.Context, req dto.CreatePekurbanHewanRequest) error {
	pekurbanID, err := uuid.Parse(req.PekurbanID)
	if err != nil {
		return errors.New("Invalid pekurban ID")
	}

	hewanID, err := uuid.Parse(req.HewanID)
	if err != nil {
		return errors.New("Invalid hewan ID")
	}

	data := &model.PekurbanHewan{
		PekurbanID: pekurbanID,
		HewanID:    hewanID,
		Porsi:      req.Porsi,
	}

	return s.repo.Create(ctx, data)
}

func (s *pekurbanHewanService) GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]dto.PekurbanHewanResponse, error) {
	list, err := s.repo.GetByHewanId(ctx, hewanID)
	if err != nil {
		return nil, err
	}

	var res []dto.PekurbanHewanResponse
	for _, ph := range list{
		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID: ph.PekurbanID.String(),
			HewanID:    ph.HewanID.String(),
			Porsi:      ph.Porsi,
		})
	}

	return res, nil
}

func (s *pekurbanHewanService) GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]dto.PekurbanHewanResponse, error) {
	list, err := s.repo.GetByPekurbanId(ctx, pekurbanID)
	if err != nil {
		return nil, err
	}

	var res []dto.PekurbanHewanResponse
	for _, ph := range list{
		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID: ph.PekurbanID.String(),
			HewanID:    ph.HewanID.String(),
			Porsi:      ph.Porsi,
		})
	}

	return res, nil
}

func (s *pekurbanHewanService) Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error {
	return s.repo.Delete(ctx, pekurbanID, hewanID)
}