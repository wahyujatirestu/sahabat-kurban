package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
)

type PekurbanHewanService interface {
	Create(ctx context.Context, req dto.CreatePekurbanHewanRequest) (*dto.PekurbanHewanResponse, error)
	GetAll(ctx context.Context) ([]dto.PekurbanHewanResponse, error)
	GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
	GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]dto.PekurbanHewanResponse, error)
	Update(ctx context.Context, pekurbanID, hewanID uuid.UUID, req dto.UpdatePekurbanHewanRequest) (*dto.PekurbanHewanResponse, error)
	Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error
}

type pekurbanHewanService struct {
	repo 		 repository.PekurbanHewanRepository
	pRepo		 repository.PekurbanRepository
	hRepo 		 repository.HewanKurbanRepository
}

func NewPekurbanHewanService(repo repository.PekurbanHewanRepository, pRepo repository.PekurbanRepository, hRepo repository.HewanKurbanRepository) PekurbanHewanService {
	return &pekurbanHewanService{repo: repo, pRepo: pRepo, hRepo: hRepo}
}

func (s *pekurbanHewanService) Create(ctx context.Context, req dto.CreatePekurbanHewanRequest) (*dto.PekurbanHewanResponse, error) {
	pekurbanID, err := uuid.Parse(req.PekurbanID)
	if err != nil {
		return nil, errors.New("Invalid pekurban ID")
	}

	pekurban, err := s.pRepo.FindById(ctx, pekurbanID)
	if err != nil {
		return nil, err
	}
	if pekurban == nil {
		return nil, errors.New("Pekurban not found")
	}

	hewanID, err := uuid.Parse(req.HewanID)
	if err != nil {
		return nil, errors.New("Invalid hewan ID")
	}

	hewan, err := s.hRepo.GetById(ctx, hewanID)
	if err != nil {
		return nil, err
	}
	if hewan == nil {
		return nil, errors.New("Hewan kurban not found")
	}

	var porsi float64
	switch hewan.Jenis {
	case "sapi":
		if req.JumlahOrang > 7 {
			return nil, errors.New("The maximum number of people for a sapi is 7")
		}
		porsi = float64(req.JumlahOrang) / 7.0
	case "kambing", "domba":
		if req.JumlahOrang != 1 {
			return nil, errors.New("Kambing/domba can only be for one person")
		}
		porsi = 1.0
	default:
		return nil, errors.New("Invalid jenis hewan")
	}

	existing, err := s.repo.GetByHewanId(ctx, hewanID)
	if err != nil {
		return nil, err
	}

	if hewan.IsPrivate {
		if len(existing) > 0 {
			return nil, errors.New("Hewan private can only be owned by one pekurban")
		}
		if porsi != 1.0 {
			return nil, errors.New("Hewan private must be fully owned (porsi 1.0)")
		}
	}

	var totalPorsi float64
	for _, ph := range existing {
		totalPorsi += ph.Porsi
	}
	if totalPorsi+porsi > 1.0 {
		return nil, errors.New("Total portion exceeds the maximum limit")
	}

	data := &model.PekurbanHewan{
		PekurbanID: pekurbanID,
		HewanID:    hewanID,
		Porsi:      porsi,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	resp := &dto.PekurbanHewanResponse{
		PekurbanID:  data.PekurbanID.String(),
		Pekurban:    *pekurban.Name,
		Hewan:       string(hewan.Jenis),
		HewanID:     data.HewanID.String(),
		Porsi:       data.Porsi,
	}

	return resp, nil
}

func (s *pekurbanHewanService) GetAll(ctx context.Context) ([]dto.PekurbanHewanResponse, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.PekurbanHewanResponse
	for _, rel := range list {
		jumlahOrang := 0
		
		switch rel.Hewan {
		case "sapi":
			jumlahOrang = int(math.Round(rel.Porsi * 7.0))
		case "kambing", "domba":
			jumlahOrang = 1
		default:
			return nil, errors.New("Invalid jenis hewan")
		}
		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID: rel.PekurbanID,
			Pekurban:   rel.Pekurban,
			HewanID:    rel.HewanID,
			Hewan:      rel.Hewan,
			Porsi:      rel.Porsi,
			JumlahOrang: jumlahOrang,
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
		jumlahOrang := 0			
		
		switch ph.Hewan {
		case "sapi":
			jumlahOrang = int(math.Round(ph.Porsi * 7.0))
		case "kambing", "domba":
			jumlahOrang = 1
		default:
			return nil, errors.New("Invalid jenis hewan")
		}

		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID:  ph.PekurbanID,
			Pekurban:    ph.Pekurban,
			HewanID:     ph.HewanID,
			Hewan:       ph.Hewan,
			Porsi:       ph.Porsi,
			JumlahOrang: jumlahOrang,
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
		jumlahOrang := 0			
		
		switch ph.Hewan {
		case "sapi":
			jumlahOrang = int(math.Round(ph.Porsi * 7.0))
		case "kambing", "domba":
			jumlahOrang = 1
		default:
			return nil, errors.New("Invalid jenis hewan")
		}

		res = append(res, dto.PekurbanHewanResponse{
			PekurbanID: ph.PekurbanID,
			Pekurban:   ph.Pekurban,
			HewanID:    ph.HewanID,
			Hewan:      ph.Hewan,
			Porsi:      ph.Porsi,
			JumlahOrang: jumlahOrang,
		})
	}

	return res, nil
}

func (s *pekurbanHewanService) Update(ctx context.Context, pekurbanID, hewanID uuid.UUID, req dto.UpdatePekurbanHewanRequest) (*dto.PekurbanHewanResponse, error) {
	if req.JumlahOrang <= 0 || req.JumlahOrang > 7 {
		return nil, errors.New("jumlah orang harus antara 1 dan 7")
	}

	hewan, err := s.hRepo.GetById(ctx, hewanID)
	if err != nil {
		return nil, err
	}
	if hewan == nil {
		return nil, errors.New("hewan tidak ditemukan")
	}

	// Validasi kambing/domba hanya 1 orang
	if !hewan.IsPrivate && (hewan.Jenis == "kambing" || hewan.Jenis == "domba") && req.JumlahOrang != 1 {
		return nil, errors.New("kambing dan domba hanya boleh 1 orang")
	}

	// Hitung porsi otomatis
	porsi := float64(req.JumlahOrang) / 7.0
	if hewan.IsPrivate || hewan.Jenis == "kambing" || hewan.Jenis == "domba" {
		porsi = 1.0
	}

	// Validasi total porsi
	existing, err := s.repo.GetByHewanId(ctx, hewanID)
	if err != nil {
		return nil, err
	}

	pekurbanId, err := s.pRepo.FindById(ctx, pekurbanID)
	if err != nil {
		return nil, err
	}
	if pekurbanId == nil {
		return nil, errors.New("pekurban tidak ditemukan")
	}

	var total float64
	for _, ph := range existing {
		if ph.PekurbanID != pekurbanID.String() {
			total += ph.Porsi
		}
	}
	if total+porsi > 1.0 {
		return nil, errors.New("total porsi melebihi batas maksimal")
	}

	// Build struct model dan update
	data := &model.PekurbanHewan{
		PekurbanID: pekurbanID,
		HewanID:    hewanID,
		Porsi:      porsi,
	}

	err = s.repo.Update(ctx, data)

	if err != nil {
		return nil, err
	}

	return &dto.PekurbanHewanResponse{
		PekurbanID: data.PekurbanID.String(),
		Pekurban:   *pekurbanId.Name,
		HewanID:    data.HewanID.String(),
		Hewan:      string(hewan.Jenis),
		Porsi:      data.Porsi,
		JumlahOrang: req.JumlahOrang,
	}, nil
}

func (s *pekurbanHewanService) Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error {
	return s.repo.Delete(ctx, pekurbanID, hewanID)
}