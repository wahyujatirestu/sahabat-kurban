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

type PenyembelihanService interface {
	Create(ctx context.Context, req dto.CreatePenyembelihanRequest) (*dto.PenyembelihanResponse, error)
	GetAll(ctx context.Context) ([]dto.PenyembelihanResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.PenyembelihanResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdatePenyembelihanRequest) (*dto.PenyembelihanResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type penyembelihanService struct {
	repo  repository.PenyembelihanRepository
	pRepo repository.PembayaranKurbanRepository
}

func NewPenyembelihanService(repo repository.PenyembelihanRepository, pRepo repository.PembayaranKurbanRepository) PenyembelihanService {
	return &penyembelihanService{repo: repo, pRepo: pRepo}
}

func (s *penyembelihanService) Create(ctx context.Context, req dto.CreatePenyembelihanRequest) (*dto.PenyembelihanResponse, error) {
	hewanID, err := uuid.Parse(req.HewanID)
	if err != nil {
		return nil, err
	}

	isLunas, err := s.pRepo.IsHewanLunas(ctx, hewanID)
	if err != nil {
		return nil, err
	}

	if !isLunas {
		return nil, errors.New("Hewan is not fully paid yet and cannot be slaughtered.")
	}

	if strings.TrimSpace(req.Lokasi) == "" {
		return nil, errors.New("Lokasi is required")
	}

	p := &model.Penyembelihan{
		HewanID: hewanID,
		TglPenyembelihan: req.TanggalPenyembelihan,
		Lokasi: req.Lokasi,
		UrutanRencana: req.UrutanRencana,
		UrutanAktual: nil,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	res := dto.ToPenyembelihanResponse(p)
	return &res, nil
}

func (s *penyembelihanService) GetAll(ctx context.Context) ([]dto.PenyembelihanResponse, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.PenyembelihanResponse
	for _, p := range list {
		res = append(res, dto.ToPenyembelihanResponse(p))
	}

	return res, nil
}

func (s *penyembelihanService) GetById(ctx context.Context, id uuid.UUID) (*dto.PenyembelihanResponse, error) {
	p, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.ToPenyembelihanResponse(p)
	return &res, nil
}

func (s *penyembelihanService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePenyembelihanRequest) (*dto.PenyembelihanResponse, error) {
	existing, err := s.repo.GetById(ctx, id)
	if err != nil || existing == nil {
		return nil, err
	}

	existing.TglPenyembelihan = req.TanggalPenyembelihan
	if strings.TrimSpace(req.Lokasi) != "" {
		existing.Lokasi = req.Lokasi
	}
	existing.UrutanRencana = req.UrutanRencana
	existing.UrutanAktual = req.UrutanAktual

	if err = s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	res := dto.ToPenyembelihanResponse(existing)
	return &res, nil
}

func (s *penyembelihanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}