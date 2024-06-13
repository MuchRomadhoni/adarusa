package debitur

import (
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/model/dto/debiturDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/src/midtrans"
)

type DebiturUsecase interface {
	PengajuanPinjaman(id int, jumlahPinjaman float64, tenor int, description string) error
	GetPengajuanPinjaman(id int) ([]debiturDto.GetPengajuanResponse, error)
	GetCicilan(page, limit, offset int, id string, status string) ([]debiturDto.GetCicilanResponse, json.Paging, error)
	CicilanPayment(pinjamanId int, jumlahBayar float64) (dto.MidtransSnapResponse, error)
	CicilanVerify(id int) error
}

type DebiturRepository interface {
	AddPengajuanPinjaman(id int, jumlahPinjaman float64, tenor int, description string) error
	GetPengajuanPinjaman(id int) ([]debiturDto.GetPengajuanResponse, error)
	GetCicilan(page, limit, offset int, id string, status string) ([]debiturDto.GetCicilanResponse, json.Paging, error)
	CicilanPayment(pinjamanId int, jumlahBayar float64) (dto.MidtransSnapResponse, error)
	CicilanVerify(id int) error
	UpdatePinjamanStatus(id int) error
	SetMidtransService(service midtrans.MidtransService)
}
