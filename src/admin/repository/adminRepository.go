package adminRepository

import (
	"database/sql"
	"fmt"
	adminEntity "fp_pinjaman_online/model/entity/admin"
	adminInterface "fp_pinjaman_online/src/admin"
	"log"
	"time"
)

type adminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) adminInterface.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) RetrieveUserStatusById(userID int) (*adminEntity.UserCompleteInfo, error) {
	query := `
    SELECT 
		u.id AS user_id,
		u.email,
		u.status,
		u.role_id,
		r.account_number,
		r.bank_name,
		k.name AS emergency_contact_name,
		k.phone_number AS emergency_contact_phone,
		j.job_name,
		j.gaji AS gaji,
		j.office_name,
		d.NIK,
		d.fullname,
		d.phone_number AS personal_phone_number,
		d.address AS personal_address,
		d.city,
		d.foto_ktp,
		d.foto_selfie,
		u.created_at,
		u.updated_at,
		u.verified_at,
		u.deleted_at
	FROM 
    	users u
	LEFT JOIN rekening r ON u.id = r.user_id
	LEFT JOIN kontak_darurat k ON u.id = k.user_id
	LEFT JOIN users_job_detail j ON u.id = j.user_id
	LEFT JOIN detail_users d ON u.id = d.user_id
	WHERE u.id = $1;
    `

	uci := adminEntity.UserCompleteInfo{}
	err := r.db.QueryRow(query, userID).Scan(
		&uci.UserID,
		&uci.Email,
		&uci.Status,
		&uci.RoleID,
		&uci.AccountNumber,
		&uci.BankName,
		&uci.EmergencyContact,
		&uci.EmergencyPhone,
		&uci.JobName,
		&uci.Gaji,
		&uci.OfficeName,
		&uci.NIK,
		&uci.FullName,
		&uci.PersonalPhoneNumber,
		&uci.PersonalAddress,
		&uci.City,
		&uci.FotoKTP,
		&uci.FotoSelfie,
		&uci.CreatedAt,
		&uci.UpdatedAt,
		&uci.VerifiedAt,
		&uci.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with ID: %d", userID)
			return nil, nil
		}
		log.Printf("Error fetching user details: %v", err)
		return nil, err
	}

	return &uci, nil
}

func (r *adminRepository) UpdateUserStatus(userID int, status string) error {
	query := `
	UPDATE users
		SET status = $2,
    		verified_at = NOW(),
    		updated_at = NOW()
		WHERE id = $1;
    `
	_, err := r.db.Exec(query, userID, status)
	if err != nil {
		log.Printf("Error updating user status: %v", err)
		return err
	}

	return nil
}

func (r *adminRepository) InsertLimitId(limitID, userID int) error {
	query := `UPDATE detail_users SET limit_id= $1 WHERE user_id = $2`
	_, err := r.db.Exec(query, limitID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepository) RetrieveUserLimitByAdmin(userID int) (*adminEntity.UserCompleteInfoLoanLimit, error) {
	query := `
    SELECT 
        u.id AS user_id,
        u.email,
        u.status,
        uj.job_name,
        uj.gaji AS gaji,
        uj.office_name,
        du.NIK,
        du.fullname,
        du.phone_number AS personal_phone_number,
        du.address AS personal_address,
        du.city,
        du.foto_ktp,
        du.foto_selfie,
        lp.max_pinjaman
    FROM 
        users u
    LEFT JOIN users_job_detail uj ON u.id = uj.user_id
    LEFT JOIN detail_users du ON u.id = du.user_id
    LEFT JOIN limit_pinjaman lp ON du.limit_id = lp.id
    WHERE u.id = $1;
    `

	uci := adminEntity.UserCompleteInfoLoanLimit{}
	err := r.db.QueryRow(query, userID).Scan(
		&uci.UserID,
		&uci.Email,
		&uci.Status,
		&uci.JobName,
		&uci.Gaji,
		&uci.OfficeName,
		&uci.NIK,
		&uci.FullName,
		&uci.PersonalPhoneNumber,
		&uci.PersonalAddress,
		&uci.City,
		&uci.FotoKTP,
		&uci.FotoSelfie,
		&uci.MaxPinjaman,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with ID: %d", userID)
			return nil, nil
		}
		log.Printf("Error fetching user details: %v", err)
		return nil, err
	}

	return &uci, nil
}

func (r *adminRepository) RetrievePinjamanById(loanID int) (*adminEntity.Pinjaman, error) {
	query := `
        SELECT id, user_id, admin_id, jumlah_pinjaman, tenor, bunga_per_bulan, description, 
               status_pengajuan, created_at, updated_at, deleted_at
        FROM pinjaman
        WHERE id = $1;
    `

	var pinjaman adminEntity.Pinjaman

	row := r.db.QueryRow(query, loanID)
	err := row.Scan(
		&pinjaman.ID, &pinjaman.UserID, &pinjaman.AdminID, &pinjaman.JumlahPinjaman,
		&pinjaman.Tenor, &pinjaman.BungaPerBulan, &pinjaman.Description,
		&pinjaman.StatusPengajuan, &pinjaman.CreatedAt, &pinjaman.UpdatedAt, &pinjaman.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No Pinjaman found with ID %d", loanID)
			return nil, err
		}
		log.Printf("Error querying for Pinjaman with ID %d: %v", loanID, err)
		return nil, err
	}

	return &pinjaman, nil
}

func (r *adminRepository) UpdateLoanStatus(loanID int, status string) error {
	query := `
	UPDATE pinjaman
		SET status_pengajuan = $2,
    		updated_at = NOW()
		WHERE id = $1;
    `
	_, err := r.db.Exec(query, loanID, status)
	if err != nil {
		log.Printf("Error updating pinjaman status: %v", err)
		return err
	}

	return nil
}

func (r *adminRepository) InsertCicilan(loanID int, dueDate time.Time, amount float64, status string) error {
	query := `
        INSERT INTO cicilan (pinjaman_id, tanggal_jatuh_tempo, jumlah_bayar, status)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.Exec(query, loanID, dueDate, amount, status)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepository) RetrieveTugasById(tugasID int) (adminEntity.ClaimTugas, error) {
	var tugas adminEntity.ClaimTugas
	query := `SELECT id, user_id, collector_id, status, created_at, updated_at, deleted_at 
              FROM claim_tugas 
              WHERE id = $1`

	err := r.db.QueryRow(query, tugasID).Scan(
		&tugas.TugasID,
		&tugas.UserID,
		&tugas.CollectorID,
		&tugas.StatusTugas,
		&tugas.CreatedAt,
		&tugas.UpdatedAt,
		&tugas.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return adminEntity.ClaimTugas{}, fmt.Errorf("no claim task found with ID %d", tugasID)
		}
		return adminEntity.ClaimTugas{}, fmt.Errorf("error retrieving claim task by ID %d: %v", tugasID, err)
	}

	return tugas, nil
}

func (r *adminRepository) UpdateClaimTugas(tugasID int, status string) error {
	sql := "UPDATE claim_tugas SET status = $2, updated_at = NOW() WHERE id = $1"
	_, err := r.db.Exec(sql, tugasID, status)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepository) RetrieveBalanceDCById(id int) (adminEntity.Balance, error) {
	var balance adminEntity.Balance
	err := r.db.QueryRow("SELECT id, user_id,amount FROM balance WHERE id = $1", id).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Amount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return adminEntity.Balance{}, fmt.Errorf("no balance found with ID %d", id)
		}
		return adminEntity.Balance{}, fmt.Errorf("error retrieving withdrawal by ID %d: %v", id, err)
	}
	return balance, nil
}

func (r *adminRepository) RetrieveWithdrawalById(withdrawalID int) (adminEntity.Withdrawal, error) {
	var withdraw adminEntity.Withdrawal
	err := r.db.QueryRow("SELECT id, user_id, amount, status, created_at, updated_at, deleted_at FROM withdrawal WHERE id = $1", withdrawalID).Scan(
		&withdraw.ID,
		&withdraw.UserID,
		&withdraw.Amount,
		&withdraw.Status,
		&withdraw.CreatedAt,
		&withdraw.UpdatedAt,
		&withdraw.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return adminEntity.Withdrawal{}, fmt.Errorf("no withdrawal found with ID %d", withdrawalID)
		}
		return adminEntity.Withdrawal{}, fmt.Errorf("error retrieving withdrawal by ID %d: %v", withdrawalID, err)
	}
	return withdraw, nil
}

func (r *adminRepository) UpdateWithdrawalStatus(withdrawalID int, newStatus string) error {
	currentTime := time.Now()

	query := "UPDATE withdrawal SET status = $1, updated_at = $2 WHERE id = $3"

	result, err := r.db.Exec(query, newStatus, currentTime, withdrawalID)
	if err != nil {

		return fmt.Errorf("error updating withdrawal status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf("error checking affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no withdrawal found with ID %d", withdrawalID)
	}

	return nil
}

func (r *adminRepository) UpdateBalance(userID int, amount float64) error {
	var storedId int
	queryInnit := `SELECT user_id FROM balance WHERE user_id = $1`
	err := r.db.QueryRow(queryInnit, userID).Scan(&storedId)
	if err != nil {
		if err == sql.ErrNoRows {
			queryInsert := `INSERT INTO balance (user_id, amount) VALUES ($1, $2)`
			_, err := r.db.Exec(queryInsert, userID, amount)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("error updating balance for user %d: %v", userID, err)
	}

	query := `UPDATE balance SET amount = amount + $1 WHERE user_id = $2`
	_, err = r.db.Exec(query, amount, userID)
	if err != nil {
		return fmt.Errorf("error updating balance for user %d: %v", userID, err)
	}
	return nil
}
