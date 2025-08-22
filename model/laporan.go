package model

import "time"

// Filter global laporan
type ReportFilter struct {
	StartDate *time.Time // optional
	EndDate   *time.Time // optional
}

// ===== Domain aggregates (hasil olahan di repository) =====

type PekurbanAggregate struct {
	PekurbanID      string
	Name            string
	Phone           string
	Email           string
	TotalHewan      int
	TotalBayar      float64 // settlement only
	TotalKewajiban  float64 
}

type PekurbanHewanDetail struct {
	PekurbanID string
	HewanID    string
	Jenis      string  // sapi/kambing/domba
	Berat      float64 // kg
	Harga      float64 // 0 jika private
	Porsi      float64 // 0< porsi <=1
}

type HewanAggregate struct {
	Jenis string
	Total int
	Berat float64
	Harga float64 // total harga (0 termasuk hewan private)
}

type DistribusiAggregate struct {
	Status        string // warga/dhuafa/panitia/pekurban
	TotalPenerima int
	TotalPaket    int
}

type PembayaranAggregate struct {
	Status string // pending/settlement/failed/expired/deny
	Total  int
	Jumlah float64
}

// Ringkasan angka-angka besar
type ReportSummary struct {
	TotalPekurban           int
	TotalHewan              int
	TotalPenerima           int
	TotalPaketDistribusi    int
	TotalPembayaranSettlement float64
}

// ===== Domain view untuk service → dto =====

type ConsolidatedReport struct {
	Tanggal   time.Time
	Summary   ReportSummary
	Pekurban  []PekurbanDetailView
	RekapHewan      []HewanAggregate
	RekapDistribusi []DistribusiAggregate
	RekapPembayaran []PembayaranAggregate
}

type PekurbanDetailView struct {
	PekurbanID       string
	Name             string
	Phone            string
	Email            string
	TotalHewan       int
	TotalBayar       float64
	TotalKewajiban   float64
	StatusPembayaran string // "lunas" | "belum"
	Hewan            []PekurbanHewanDetail
}
