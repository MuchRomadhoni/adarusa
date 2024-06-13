package debiturUsecase

import (
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/model/dto/debiturDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/src/debitur"
)

type DebiturUsecase struct {
	debiturRepository debitur.DebiturRepository
}

func NewDebiturUsecase(debiturRepository debitur.DebiturRepository) debitur.DebiturUsecase {
	return &DebiturUsecase{debiturRepository}
}

func (u *DebiturUsecase) PengajuanPinjaman(id int, jumlahPinjaman float64, tenor int, description string) error {

	err := u.debiturRepository.AddPengajuanPinjaman(id, jumlahPinjaman, tenor, description)
	if err != nil {
		return err
	}
	return nil
}

func (u *DebiturUsecase) GetPengajuanPinjaman(id int) ([]debiturDto.GetPengajuanResponse, error) {
	data, err := u.debiturRepository.GetPengajuanPinjaman(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *DebiturUsecase) GetCicilan(page, limit, offset int, id string, status string) ([]debiturDto.GetCicilanResponse, json.Paging, error) {
	data, paging, err := u.debiturRepository.GetCicilan(page, limit, offset, id, status)
	if err != nil {
		return nil, json.Paging{}, err
	}
	return data, paging, nil
}

func (u *DebiturUsecase) CicilanPayment(pinjamanId int, totalBayar float64) (dto.MidtransSnapResponse, error) {
	data, err := u.debiturRepository.CicilanPayment(pinjamanId, totalBayar)
	if err != nil {
		return dto.MidtransSnapResponse{}, err
	}
	return data, nil
}

func (u *DebiturUsecase) CicilanVerify(id int) error {
	err := u.debiturRepository.CicilanVerify(id)
	if err != nil {
		return err
	}
	err = u.debiturRepository.UpdatePinjamanStatus(id)
	return nil
}
