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
	GetAll(ctx context.Context) ([]dto.PekurbanHewanResponse, error)
	GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
	GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
	Update(ctx context.Context, pekurbanID, hewanID uuid.UUID, req dto.UpdatePekurbanHewanRequest) error
	Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error
}

type pekurbanHewanService struct {
	repo 	repository.PekurbanHewanRepository
	hRepo 	repository.HewanKurbanRepository
}

func NewPekurbanHewanService(repo repository.PekurbanHewanRepository, hRepo repository.HewanKurbanRepository) PekurbanHewanService {
	return &pekurbanHewanService{repo: repo, hRepo: hRepo}
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

	hewan, err := s.hRepo.GetById(ctx, hewanID)
	if err != nil {
		return err
	}
	if hewan == nil {
		return errors.New("Hewan kurban not found")
	}

	var porsi float64
	switch hewan.Jenis {
	case "sapi":
		porsi = float64(req.JumlahOrang) / 7.0
		if req.JumlahOrang > 7 {
			return errors.New("The maximum number of people for a sapi is 7")
		}
	case "kambing", "domba":
		if req.JumlahOrang != 1 {
			return errors.New("Kambing/domba can only be for one person")
		}
		porsi = 1.0
	default:
		return errors.New("Invalid jenis hewan")
	}

	existing, err := s.repo.GetByHewanId(ctx, hewanID)
	if err != nil {
		return err
	}

	if hewan.IsPrivate {
		if len(existing) > 0 {
			return errors.New("Hewan private can only be owned by one pekurban")
		}
		if porsi != 1.0 {
			return errors.New("Hewan private must be fully owned (porsi 1.0)")
		}
	}

	var totalPorsi float64
	for _, ph := range existing {
		totalPorsi += ph.Porsi
	}
	if totalPorsi+porsi > 1.0 {
		return errors.New("Total portion exceeds the maximum limit")
	}

	data := &model.PekurbanHewan{
		PekurbanID: pekurbanID,
		HewanID:    hewanID,
		Porsi:      porsi,
	}

	return s.repo.Create(ctx, data)
}


func (s *pekurbanHewanService) GetAll(ctx context.Context) ([]dto.PekurbanHewanResponse, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.PekurbanHewanResponse
	for _, rel := range list {
		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID: rel.PekurbanID.String(),
			HewanID:    rel.HewanID.String(),
			Porsi:      rel.Porsi,
		})
	}
	return res, nil
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

func (s *pekurbanHewanService) Update(ctx context.Context, pekurbanID, hewanID uuid.UUID, req dto.UpdatePekurbanHewanRequest) error {
	if req.JumlahOrang <= 0 || req.JumlahOrang > 7 {
		return errors.New("jumlah orang harus antara 1 dan 7")
	}

	hewan, err := s.hRepo.GetById(ctx, hewanID)
	if err != nil {
		return err
	}
	if hewan == nil {
		return errors.New("hewan tidak ditemukan")
	}

	// Validasi kambing/domba hanya 1 orang
	if !hewan.IsPrivate && (hewan.Jenis == "kambing" || hewan.Jenis == "domba") && req.JumlahOrang != 1 {
		return errors.New("kambing dan domba hanya boleh 1 orang")
	}

	// Hitung porsi otomatis
	porsi := float64(req.JumlahOrang) / 7.0
	if hewan.IsPrivate || hewan.Jenis == "kambing" || hewan.Jenis == "domba" {
		porsi = 1.0
	}

	// Validasi total porsi
	existing, err := s.repo.GetByHewanId(ctx, hewanID)
	if err != nil {
		return err
	}

	var total float64
	for _, ph := range existing {
		if ph.PekurbanID != pekurbanID {
			total += ph.Porsi
		}
	}
	if total+porsi > 1.0 {
		return errors.New("total porsi melebihi batas maksimal")
	}

	// Build struct model dan update
	data := &model.PekurbanHewan{
		PekurbanID: pekurbanID,
		HewanID:    hewanID,
		Porsi:      porsi,
	}

	return s.repo.Update(ctx, data)
}



func (s *pekurbanHewanService) Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error {
	return s.repo.Delete(ctx, pekurbanID, hewanID)
}