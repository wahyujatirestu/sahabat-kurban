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

type PekurbanService interface {
	Create(ctx context.Context, req dto.CreatePekurbanRequest) (*dto.PekurbanResponse, error)
	GetAll(ctx context.Context) ([]dto.PekurbanResponse, error)
	GetById(ctx context.Context, id uuid.UUID)(*dto.PekurbanResponse, error)
	GetByUserId(ctx context.Context, userID uuid.UUID)(*dto.PekurbanResponse, error)	
	Update(ctx context.Context, id uuid.UUID, req dto.UpdatePekurbanRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type pekurbanService struct {
	pRepo repository.PekurbanRepository
	uRepo repository.UserRepository
}

func NewPekurbanService(pRepo repository.PekurbanRepository, uRepo repository.UserRepository) PekurbanService {
	return &pekurbanService{pRepo: pRepo, uRepo: uRepo}
}

func (s *pekurbanService) Create(ctx context.Context, req dto.CreatePekurbanRequest) (*dto.PekurbanResponse, error) {
	var userID *uuid.UUID
	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err == nil {
			userID = &uid

			if req.Name == "" || req.Email == "" {
				user, err := s.uRepo.FindById(ctx, uid)
				if err != nil || user == nil {
					return nil, errors.New("Invalid user ID")
				}

				if req.Name == "" {
					req.Name = user.Name
				}
				if req.Email == "" {
					req.Email = user.Email
				}
			}
		}
	}

	p := &model.Pekurban{
		ID: uuid.New(),
		UserId: userID,
		Name: &req.Name,
		Phone: &req.Phone,
		Email: &req.Email,
		Alamat: &req.Alamat,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.pRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	
	res := dto.ToPekurbanRespon(p)

	return &res, nil
}

func (s *pekurbanService) GetAll(ctx context.Context) ([]dto.PekurbanResponse, error) {
	data, err := s.pRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.PekurbanResponse
	for _, p := range data{
		result = append(result, dto.ToPekurbanRespon(p))
	}

	return result, nil
}

func (s *pekurbanService) GetById(ctx context.Context, id uuid.UUID)(*dto.PekurbanResponse, error) {
	p, err := s.pRepo.FindById(ctx, id)
	if err != nil || p == nil {
		return nil, err
	}

	res := dto.ToPekurbanRespon(p)

	return &res, nil
}

func (s *pekurbanService) GetByUserId(ctx context.Context, userID uuid.UUID)(*dto.PekurbanResponse, error) {
	p, err := s.pRepo.FindByUserId(ctx, userID)
	if err != nil || p == nil {
		return nil, err
	}

	res := dto.ToPekurbanRespon(p)

	return &res, nil
}

func (s *pekurbanService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePekurbanRequest) error {
	p, err := s.pRepo.FindById(ctx, id)
	if err != nil || p == nil {
		return err
	}

	p.Name = &req.Name
	p.Phone = &req.Phone
	p.Email = &req.Email
	p.Alamat = &req.Alamat

	return s.pRepo.Update(ctx, p)
}

func (s *pekurbanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.pRepo.Delete(ctx, id)
}