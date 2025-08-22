package service

import (
	"context"
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
)

type ReportService interface {
	GetConsolidatedReport(ctx context.Context, f model.ReportFilter) (*dto.LaporanResponse, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetConsolidatedReport(ctx context.Context, f model.ReportFilter) (*dto.LaporanResponse, error) {
	// parallel fetch bisa dipertimbangkan; di sini sequential untuk kesederhanaan
	pekurbanAgg, err := s.repo.GetPekurbanAggregates(ctx, f)
	if err != nil {
		return nil, err
	}
	hewanDetails, err := s.repo.GetPekurbanHewanDetails(ctx, f)
	if err != nil {
		return nil, err
	}
	rekapHewan, err := s.repo.GetHewanAggregate(ctx, f)
	if err != nil {
		return nil, err
	}
	rekapDistribusi, err := s.repo.GetDistribusiAggregate(ctx, f)
	if err != nil {
		return nil, err
	}
	rekapPembayaran, err := s.repo.GetPembayaranAggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	// summary
	totalPekurban, err := s.repo.CountPekurban(ctx)
	if err != nil {
		return nil, err
	}
	totalHewan, err := s.repo.CountHewan(ctx, f)
	if err != nil {
		return nil, err
	}
	totalPenerima, err := s.repo.CountPenerima(ctx)
	if err != nil {
		return nil, err
	}
	totalPaket, err := s.repo.SumPaketDistribusi(ctx, f)
	if err != nil {
		return nil, err
	}
	totalPembayaran, err := s.repo.SumPembayaranSettlement(ctx, f)
	if err != nil {
		return nil, err
	}

	// map hewan detail per pekurban
	hewanByPekurban := map[string][]model.PekurbanHewanDetail{}
	for _, h := range hewanDetails {
		hewanByPekurban[h.PekurbanID] = append(hewanByPekurban[h.PekurbanID], h)
	}

	// build DTO pekurban + status "lunas" vs "belum"
	pekurbanDTOs := make([]dto.PekurbanDTO, 0, len(pekurbanAgg))
	for _, a := range pekurbanAgg {
		status := "belum"
		// toleransi floating kecil
		if a.TotalBayar+1e-6 >= a.TotalKewajiban {
			status = "lunas"
		}
		hd := hewanByPekurban[a.PekurbanID]
		hewanDTO := make([]dto.PekurbanHewanDTO, 0, len(hd))
		for _, h := range hd {
			hewanDTO = append(hewanDTO, dto.PekurbanHewanDTO{
				HewanID: h.HewanID,
				Jenis:   h.Jenis,
				Berat:   h.Berat,
				Harga:   h.Harga,
				Porsi:   h.Porsi,
			})
		}
		pekurbanDTOs = append(pekurbanDTOs, dto.PekurbanDTO{
			PekurbanID:       a.PekurbanID,
			Nama:             a.Name,
			Phone:            a.Phone,
			Email:            a.Email,
			TotalHewan:       a.TotalHewan,
			TotalBayar:       a.TotalBayar,
			TotalKewajiban:   a.TotalKewajiban,
			StatusPembayaran: status,
			Hewan:            hewanDTO,
		})
	}

	// map rekap hewan
	hewanDTOs := make([]dto.HewanDTO, 0, len(rekapHewan))
	for _, h := range rekapHewan {
		hewanDTOs = append(hewanDTOs, dto.HewanDTO{
			Jenis: h.Jenis,
			Total: h.Total,
			Berat: h.Berat,
			Harga: h.Harga,
		})
	}

	// map rekap distribusi
	distribusiDTOs := make([]dto.DistribusiDTO, 0, len(rekapDistribusi))
	for _, d := range rekapDistribusi {
		distribusiDTOs = append(distribusiDTOs, dto.DistribusiDTO{
			Kategori:      d.Status,
			TotalPenerima: d.TotalPenerima,
			TotalPaket:    d.TotalPaket,
		})
	}

	// map rekap pembayaran
	paymentDTOs := make([]dto.PembayaranDTO, 0, len(rekapPembayaran))
	for _, p := range rekapPembayaran {
		paymentDTOs = append(paymentDTOs, dto.PembayaranDTO{
			Status: p.Status,
			Total:  p.Total,
			Jumlah: p.Jumlah,
		})
	}

	now := time.Now()
	resp := &dto.LaporanResponse{
		Tanggal: now,
		Ringkasan: dto.RingkasanDTO{
			TotalPekurban:             totalPekurban,
			TotalHewan:                totalHewan,
			TotalPenerima:             totalPenerima,
			TotalPaketDistribusi:      totalPaket,
			TotalPembayaranSettlement: totalPembayaran,
		},
		Pekurban:         pekurbanDTOs,
		RekapHewan:       hewanDTOs,
		RekapDistribusi:  distribusiDTOs,
		RekapPembayaran:  paymentDTOs,
	}
	return resp, nil
}
