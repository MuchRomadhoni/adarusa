package debtCollectorEntity

type (
	Tugas struct {
		ID          string `json:"id"`
		UserId      string `json:"userId"`
		CollectorId string `json:"collectorId,omitempty"`
		Status      string `json:"status"`
	}

	LogTugas struct {
		ID          string `json:"id"`
		TugasId     string `json:"tugasId,omitempty"`
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt,omitempty"`
		UpdatedAt   string `json:"updatedAt,omitempty"`
	}

	LateDebtor struct {
		ID           string  `json:"id"`
		FullName     string  `json:"fullName"`
		Address      string  `json:"address"`
		UnpaidAmount float64 `json:"unpaidAmount"`
	}

	DebtCollector struct {
		ID       string `json:"id"`
		FullName string `json:"fullName"`
		City     string `json:"city"`
	}
)
