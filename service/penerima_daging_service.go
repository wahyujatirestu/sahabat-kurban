package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
)

type PenerimaDagingService interface {
	Create(ctx context.Context, req dto.CreatePenerimaRequest) (*dto.PenerimaResponse, error)
	GetAll(ctx context.Context) ([]dto.PenerimaResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PenerimaResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdatePenerimaRequest) (*dto.PenerimaResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type penerimaDagingService struct {
	repo repository.PenerimaDagingRepository
	repoPk	repository.PekurbanRepository
}

func NewPenerimaDagingService(repo repository.PenerimaDagingRepository , repoPk repository.PekurbanRepository) PenerimaDagingService {
	return &penerimaDagingService{repo: repo, repoPk: repoPk}
}

func (s *penerimaDagingService) Create(ctx context.Context, req dto.CreatePenerimaRequest) (*dto.PenerimaResponse, error) {
	var pekurbanID *uuid.UUID
	if req.PekurbanID != nil {
		id, err := uuid.Parse(*req.PekurbanID)
		if err == nil {
			pekurbanID = &id

			pekurban, err := s.repoPk.FindById(ctx, id)
			if err != nil || pekurban == nil {
				return nil, errors.New("Pekurban not found")
			}

			if strings.TrimSpace(*pekurban.Name) == "" {
				return nil, errors.New("Pekurban name is required")
			}

			if req.Name == "" && pekurban.Name != nil {
				req.Name = *pekurban.Name
			}

			if req.Alamat == nil && pekurban.Alamat != nil {
				req.Alamat = pekurban.Alamat
			}

			if req.Phone == nil && pekurban.Phone != nil {
				req.Phone = pekurban.Phone
			}
		}
	}

	p := &model.PenerimaDaging{
		ID: uuid.New(),
		Name: req.Name,
		Alamat: req.Alamat,
		Phone: req.Phone,
		Status: req.Status,
		PekurbanID: pekurbanID,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	res := dto.ToPenerimaResponse(p)
	return &res, nil
}

func (s *penerimaDagingService) GetAll(ctx context.Context) ([]dto.PenerimaResponse, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	var res []dto.PenerimaResponse
	for _, p := range list {
		res = append(res, dto.ToPenerimaResponse(p))
	}
	return res, nil
}

func (s *penerimaDagingService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PenerimaResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, err
	}

	res := dto.ToPenerimaResponse(p)
	return &res, nil
}

func (s *penerimaDagingService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePenerimaRequest) (*dto.PenerimaResponse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil || existing == nil {
		return nil, err
	}

	if strings.TrimSpace(*req.Name) != "" {
		existing.Name = *req.Name
	}
	if strings.TrimSpace(*req.Alamat) != "" {
		existing.Alamat = req.Alamat
	}
	if strings.TrimSpace(*req.Phone) != "" {
		existing.Phone = req.Phone
	}
	existing.Status = req.Status

	if req.PekurbanID != nil {
		uid, err := uuid.Parse(*req.PekurbanID)
		if err == nil {
			existing.PekurbanID = &uid
		}
	}

	if err = s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	res := dto.ToPenerimaResponse(existing)
	return &res, nil
}

func (s *penerimaDagingService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)	
}