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

type HewanKurbanService interface {
	Create(ctx context.Context, req dto.CreateHewanKurbanRequest) (*dto.HewanKurbanResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.HewanKurbanResponse, error)
	GetAll(ctx context.Context) ([]dto.HewanKurbanResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateHewanKurbanRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type hewanKurbanService struct {
	repo repository.HewanKurbanRepository
}

func NewHewanKurbanService(r repository.HewanKurbanRepository) HewanKurbanService {
	return &hewanKurbanService{repo: r}
}

func (s *hewanKurbanService) Create(ctx context.Context, req dto.CreateHewanKurbanRequest) (*dto.HewanKurbanResponse, error) {
	tanggal, err := time.Parse("2006-01-02", req.TglPendaftaran)
	if err != nil {
		return nil, errors.New("Invalid date format, must be YYYY-MM-DD")
	}

	h := &model.HewanKurban{
		ID: uuid.New(),
		Jenis: req.Jenis,
		Berat: req.Berat,
		TglPendaftaran: tanggal,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.repo.Create(ctx, h); err != nil {
		return nil, err
	}

	res := dto.ToHewanKurbanResponse(h)

	return &res, nil
}

func (s *hewanKurbanService) GetByID(ctx context.Context, id uuid.UUID) (*dto.HewanKurbanResponse, error) {
	data, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("Hewan kurban not found")
	}

	res := dto.ToHewanKurbanResponse(data)
	return &res, nil
}

func (s *hewanKurbanService) GetAll(ctx context.Context) ([]dto.HewanKurbanResponse, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	var result []dto.HewanKurbanResponse
	for _, h := range data{
		result = append(result, dto.ToHewanKurbanResponse(h))
	}

	return result, nil
}

func (s *hewanKurbanService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateHewanKurbanRequest) error {
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Hewan kurban not found")
	}

	if req.Jenis != "" {
		existing.Jenis = req.Jenis
	}
	if req.Berat > 0 {
		existing.Berat = req.Berat
	}
	if req.TglPendaftaran != "" {
		tgl, err := time.Parse("2006-01-02", req.TglPendaftaran)
		if err != nil {
			return errors.New("Invalid date format")
		}
		existing.TglPendaftaran = tgl
	}

	return s.repo.Update(ctx, existing)
}

func (s *hewanKurbanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}


