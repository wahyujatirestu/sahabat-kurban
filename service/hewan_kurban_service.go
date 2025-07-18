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
	repo 	repository.HewanKurbanRepository
	pRepo	repository.PenyembelihanRepository
}

func NewHewanKurbanService(r repository.HewanKurbanRepository, pr repository.PenyembelihanRepository) HewanKurbanService {
	return &hewanKurbanService{repo: r, pRepo: pr}
}

func (s *hewanKurbanService) Create(ctx context.Context, req dto.CreateHewanKurbanRequest) (*dto.HewanKurbanResponse, error) {
	tanggal, err := time.Parse("2006-01-02", req.TglPendaftaran)
	if err != nil {
		return nil, errors.New("Invalid date format, must be YYYY-MM-DD")
	}

	isPrivate := false
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}

	var harga float64

	if isPrivate {
		harga = 0
	} else {
		if req.Harga == nil || *req.Harga <= 0 {
			return nil, errors.New("Field harga must be filled and greater than 0 if hewan is not private")
		}
		harga = *req.Harga
	}

	h := &model.HewanKurban{
		ID: 				uuid.New(),
		Jenis: 				model.JenisHewan(req.Jenis),
		Berat: 				req.Berat,
		Harga:   			harga,
		IsPrivate:          isPrivate,
		TanggalPendaftaran: tanggal,
		Created_At: 		time.Now(),
		Updated_At: 		time.Now(),
	}

	if err := s.repo.Create(ctx, h); err != nil {
		return nil, err
	}

	res := dto.ToHewanKurbanResponse(h, false)

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

	_, err = s.pRepo.GetByHewanID(ctx, id)
	isDisembelih := err == nil

	res := dto.ToHewanKurbanResponse(data, isDisembelih)
	return &res, nil
}

func (s *hewanKurbanService) GetAll(ctx context.Context) ([]dto.HewanKurbanResponse, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	var result []dto.HewanKurbanResponse
	for _, h := range data{
		_, err := s.pRepo.GetByHewanID(ctx, h.ID)
		isDisembelih := err == nil
		result = append(result, dto.ToHewanKurbanResponse(h, isDisembelih))
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
		existing.Jenis = model.JenisHewan(req.Jenis)
	}
	if req.Berat > 0 {
		existing.Berat = req.Berat
	}
	if req.Harga > 0 {
    existing.Harga = req.Harga
	}
	if req.IsPrivate != nil {
		existing.IsPrivate = *req.IsPrivate
	}

	if req.TglPendaftaran != "" {
		tgl, err := time.Parse("2006-01-02", req.TglPendaftaran)
		if err != nil {
			return errors.New("Invalid date format")
		}
		existing.TanggalPendaftaran = tgl
	}

	return s.repo.Update(ctx, existing)
}

func (s *hewanKurbanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}


