package debiturDto

import "time"

type (
	DebiturRequest struct {
		JumlahPinjaman float64 `json:"jumlah_pinjaman" binding:"required,number"`
		Tenor          int     `json:"tenor" binding:"required,number"`
		Description    string  `json:"description" binding:"required"`
	}

	GetPengajuanResponse struct {
		Id              int       `json:"id"`
		JumlahPinjaman  float64   `json:"jumlahPinjaman"`
		Tenor           int       `json:"tenor"`
		BungaPerBulan   float64   `json:"bungaPerBulan"`
		Description     string    `json:"description"`
		StatusPengajuan string    `json:"statusPengajuan"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"endDate"`
	}

	GetCicilanResponse struct {
		Id                  int       `json:"id"`
		PinjamanId          int       `json:"pinjamanId"`
		TanggalJatuhTempo   time.Time `json:"tanggalJatuhTempo"`
		TanggalBayarSelesai time.Time `json:"tanggalBayarSelesai"`
		JumlahBayar         float64   `json:"jumlahBayar"`
		Status              string    `json:"status"`
	}

	CicilanPaymentRequest struct {
		PinjamanId int     `json:"pinjamanId" binding:"required,number"`
		TotalBayar float64 `json:"totalBayar" binding:"required,number"`
	}

	CicilanPaymentResponse struct {
		CicilanId           int       `json:"cicilanId"`
		PinjamanId          int       `json:"pinjamanId"`
		TanggalJatuhTempo   time.Time `json:"tanggalJatuhTempo"`
		TanggalBayarSelesai time.Time `json:"tanggalBayarSelesai"`
		JumlahBayar         float64   `json:"jumlahBayar"`
		Status              string    `json:"status"`
	}
)
