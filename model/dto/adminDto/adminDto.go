package adminDto

type (
	RequestUpdateStatusUser struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		LimitID int    `json:"limit_id"`
	}

	RequestVerifyLoan struct {
		LoanID          int    `json:"loan_id"`
		UserID          int    `json:"user_id"`
		AdminID         int    `json:"admin_id"`
		StatusPengajuan string `json:"status_pengajuan"`
	}

	RequestUpdateClaimTugas struct {
		ID     int    `json:"tugas_id"`
		Status string `json:"status"`
	}
	AdminResponse struct {
		ID     int    `json:"id"`
		Email  string `json:"email"`
		Status string `json:"status"`
	}

	LoanResponse struct {
		LoanID  int    `json:"loan_id"`
		UserID  int    `json:"user_id"`
		Message string `json:"message,omitempty"`
	}

	ClaimTugasResponse struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}

	RequestWithdrawal struct {
		ID     int    `json:"id"`
		UserID int    `json:"user_id"`
		Status string `json:"status"`
	}

	WithdrawalResponse struct {
		ID     int     `json:"id"`
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
		Status string  `json:"status"`
	}
)
