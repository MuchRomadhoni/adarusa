package debiturRepository

import (
	"database/sql"
	"errors"
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/model/dto/debiturDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/src/debitur"
	"fp_pinjaman_online/src/midtrans"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type debiturRepository struct {
	db              *sql.DB
	midtransService midtrans.MidtransService
}

func NewDebiturRepository(db *sql.DB) debitur.DebiturRepository {
	return &debiturRepository{db, midtrans.NewMidtransService(resty.New())}
}

func (r *debiturRepository) SetMidtransService(service midtrans.MidtransService) {
	r.midtransService = service
}

func (d *debiturRepository) AddPengajuanPinjaman(id int, jumlahPinjaman float64, tenor int, description string) error {
	var sukuBunga float64
	if tenor < 3 {
		sukuBunga = 0.05
	} else if tenor >= 3 && tenor < 6 {
		sukuBunga = 0.07
	} else if tenor >= 6 && tenor < 9 {
		sukuBunga = 0.09
	} else if tenor >= 9 && tenor <= 12 {
		sukuBunga = 0.11
	}

	_, err := d.db.Exec("INSERT INTO pinjaman (user_id, jumlah_pinjaman, tenor, description, bunga_per_bulan, status_pengajuan) VALUES ($1, $2, $3, $4, $5, 'pending')", id, jumlahPinjaman, tenor, description, sukuBunga)
	if err != nil {
		return err
	}
	return nil
}

func (d *debiturRepository) GetPengajuanPinjaman(id int) ([]debiturDto.GetPengajuanResponse, error) {
	var data []debiturDto.GetPengajuanResponse
	rows, err := d.db.Query("SELECT id, jumlah_pinjaman, tenor, description, bunga_per_bulan, status_pengajuan, created_at, updated_at FROM pinjaman WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pengajuan debiturDto.GetPengajuanResponse
		if err := rows.Scan(&pengajuan.Id, &pengajuan.JumlahPinjaman, &pengajuan.Tenor, &pengajuan.Description, &pengajuan.BungaPerBulan, &pengajuan.StatusPengajuan, &pengajuan.CreatedAt, &pengajuan.UpdatedAt); err != nil {
			return nil, err
		}
		data = append(data, pengajuan)
	}
	return data, nil
}

func (d *debiturRepository) GetCicilan(page, limit, offset int, id string, status string) ([]debiturDto.GetCicilanResponse, json.Paging, error) {
	var paging json.Paging
	var totalData int
	err := d.db.QueryRow("SELECT COUNT(*) FROM cicilan WHERE pinjaman_id = $1", id).Scan(&totalData)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, json.Paging{}, err
	}

	paging.Page = page
	paging.TotalData = totalData

	var rows *sql.Rows
	if status == "" {
		rows, err = d.db.Query("SELECT id, pinjaman_id, tanggal_jatuh_tempo, jumlah_bayar, status FROM cicilan WHERE pinjaman_id = $1", id)
	} else {
		rows, err = d.db.Query("SELECT id, pinjaman_id, tanggal_jatuh_tempo, jumlah_bayar, status FROM cicilan WHERE pinjaman_id = $1 AND status = $2", id, status)
	}
	var data []debiturDto.GetCicilanResponse
	if err != nil {
		return nil, json.Paging{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var cicilan debiturDto.GetCicilanResponse
		if err := rows.Scan(&cicilan.Id, &cicilan.PinjamanId, &cicilan.TanggalJatuhTempo, &cicilan.JumlahBayar, &cicilan.Status); err != nil {
			return nil, json.Paging{}, err
		}
		data = append(data, cicilan)
	}

	return data, paging, nil
}

func (d *debiturRepository) CicilanPayment(pinjamanId int, totalBayar float64) (dto.MidtransSnapResponse, error) {

	var cicilanId int
	var jumlahBayar float64
	err := d.db.QueryRow("SELECT id, jumlah_bayar FROM cicilan WHERE pinjaman_id = $1 AND status = 'unpaid'", pinjamanId).Scan(&cicilanId, &jumlahBayar)
	if err != nil {
		return dto.MidtransSnapResponse{}, errors.New("anda tidak mempunyai cicilan")
	}

	if totalBayar != jumlahBayar {
		return dto.MidtransSnapResponse{}, errors.New("total bayar tidak sesuai")
	}

	//get customer name
	var customerName string
	err = d.db.QueryRow("SELECT d.fullname FROM pinjaman p JOIN detail_users d ON p.user_id = d.user_id WHERE p.id = $1", pinjamanId).Scan(&customerName)
	if err != nil {
		log.Error().Msg(err.Error())
		return dto.MidtransSnapResponse{}, err
	}

	client := resty.New()
	midtransService := midtrans.NewMidtransService(client)

	payload := dto.MidtransSnapRequest{
		TransactionDetails: struct {
			OrderID  int     `json:"order_id"`
			GrossAmt float64 `json:"gross_amount"`
		}{OrderID: cicilanId, GrossAmt: totalBayar},
		PaymentType: "gopay",
		Customer:    customerName,
	}

	data, err := midtransService.Pay(payload)
	if err != nil {
		log.Error().Msg(err.Error())
		return dto.MidtransSnapResponse{}, err
	}

	_, err = d.db.Exec("INSERT INTO midtrans_tx (cicilan_id, amount, snap_url) VALUES ($1, $2, $3)", cicilanId, totalBayar, data.RedirectUrl)
	if err != nil {
		log.Error().Msg(err.Error())
		return dto.MidtransSnapResponse{}, err
	}

	return data, nil
}

func (d *debiturRepository) CicilanVerify(id int) error {

	// client := resty.New()
	midtransService := d.midtransService
	success, _ := midtransService.VerifyPayment(id)
	if !success {
		return errors.New("payment not success")
	}

	_, err := d.db.Exec("UPDATE cicilan SET status = 'paid', tanggal_selesai_bayar = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}
	_, err = d.db.Exec("UPDATE midtrans_tx SET status = 'success' WHERE cicilan_id = $1", id)
	if err != nil {
		return err
	}

	var pinjamanId int
	err = d.db.QueryRow("SELECT pinjaman_id FROM cicilan WHERE id = $1", id).Scan(&pinjamanId)
	if err != nil {
		return err
	}

	err = d.UpdatePinjamanStatus(pinjamanId)
	if err != nil {
		return err
	}
	return nil
}

func (d *debiturRepository) UpdatePinjamanStatus(id int) error {

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM cicilan WHERE pinjaman_id = $1 AND status = 'unpaid')"
	err := d.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = d.db.Exec("UPDATE pinjaman SET status_pengajuan = 'completed', updated_at = NOW() WHERE id = $1", id)
		if err != nil {
			return err
		}
	}

	return nil
}
