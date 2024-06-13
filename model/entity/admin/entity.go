package adminEntity

import (
	"database/sql"
	"time"
)

type (
	Users struct {
		ID         int        `json:"id"`
		Email      string     `json:"email"`
		Password   string     `json:"password"`
		RoleID     int        `json:"role_id"`
		Status     string     `json:"status"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		VerifiedAt *time.Time `json:"verified_at"`
		DeletedAt  *time.Time `json:"deleted_at"`
	}

	Rekening struct {
		ID            int    `json:"id"`
		UserID        int    `json:"user_id"`
		AccountNumber string `json:"account_number"`
		BankName      string `json:"bank_name"`
	}

	KontakDarurat struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	}

	UserJobDetail struct {
		ID            int     `json:"id"`
		UserID        int     `json:"user_id"`
		JobName       string  `json:"job_name"`
		Salary        float64 `json:"salary"`
		OfficeName    string  `json:"office_name"`
		OfficeContact string  `json:"office_contact"`
		Address       string  `json:"address"`
	}

	UserDetail struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		LimitID     int    `json:"limit_id"`
		NIK         string `json:"nik"`
		FullName    string `json:"fullname"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		City        string `json:"city"`
		FotoKTP     string `json:"foto_ktp"`
		FotoSelfie  string `json:"foto_selfie"`
	}

	LimitPinjaman struct {
		ID          int     `json:"id"`
		MaxPinjaman float64 `json:"max_pinjaman"`
	}

	Pinjaman struct {
		ID              int           `json:"id"`
		UserID          int           `json:"user_id"`
		AdminID         sql.NullInt64 `json:"admin_id"`
		JumlahPinjaman  float64       `json:"jumlah_pinjaman"`
		Tenor           int           `json:"tenor"`
		BungaPerBulan   float64       `json:"bunga_per_bulan"`
		Description     string        `json:"description"`
		StatusPengajuan string        `json:"status_pengajuan"`
		CreatedAt       time.Time     `json:"created_at"`
		UpdatedAt       time.Time     `json:"updated_at"`
		DeletedAt       *time.Time    `json:"deleted_at"`
	}

	UserCompleteInfo struct {
		UserID              int             `json:"user_id"`
		Email               string          `json:"email"`
		Status              string          `json:"status"`
		RoleID              int             `json:"role_id"`
		AccountNumber       sql.NullString  `json:"account_number"`
		BankName            sql.NullString  `json:"bank_name"`
		EmergencyContact    sql.NullString  `json:"emergency_contact_name"`
		EmergencyPhone      sql.NullString  `json:"emergency_contact_phone"`
		JobName             sql.NullString  `json:"job_name"`
		Gaji                sql.NullFloat64 `json:"gaji"`
		OfficeName          sql.NullString  `json:"office_name"`
		NIK                 sql.NullString  `json:"nik"`
		FullName            sql.NullString  `json:"fullname"`
		PersonalPhoneNumber sql.NullString  `json:"personal_phone_number"`
		PersonalAddress     sql.NullString  `json:"personal_address"`
		City                sql.NullString  `json:"city"`
		FotoKTP             sql.NullString  `json:"foto_ktp"`
		FotoSelfie          sql.NullString  `json:"foto_selfie"`
		CreatedAt           time.Time       `json:"created_at"`
		UpdatedAt           time.Time       `json:"updated_at"`
		VerifiedAt          *time.Time      `json:"verified_at"`
		DeletedAt           *time.Time      `json:"deleted_at"`
	}

	UserCompleteInfoLoanLimit struct {
		UserID              int             `json:"user_id"`
		Email               string          `json:"email"`
		Status              string          `json:"status"`
		JobName             sql.NullString  `json:"job_name"`
		Gaji                sql.NullFloat64 `json:"gaji"`
		OfficeName          sql.NullString  `json:"office_name"`
		NIK                 sql.NullString  `json:"nik"`
		FullName            sql.NullString  `json:"fullname"`
		PersonalPhoneNumber sql.NullString  `json:"personal_phone_number"`
		PersonalAddress     sql.NullString  `json:"personal_address"`
		City                sql.NullString  `json:"city"`
		FotoKTP             sql.NullString  `json:"foto_ktp"`
		FotoSelfie          sql.NullString  `json:"foto_selfie"`
		MaxPinjaman         sql.NullFloat64 `json:"max_pinjaman"`
		CreatedAt           time.Time       `json:"created_at"`
		UpdatedAt           time.Time       `json:"updated_at"`
		VerifiedAt          *time.Time      `json:"verified_at"`
		DeletedAt           *time.Time      `json:"deleted_at"`
	}

	ClaimTugas struct {
		TugasID     int        `json:"id"`
		UserID      int        `json:"user_id"`
		CollectorID int        `json:"collector_id"`
		StatusTugas string     `json:"status"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		DeletedAt   *time.Time `json:"deleted_at"`
	}

	Balance struct {
		ID     int     `json:"id"`
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	Withdrawal struct {
		ID        int        `json:"id"`
		UserID    int        `json:"user_id"`
		Amount    float64    `json:"amount"`
		Status    string     `json:"status"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}
)
