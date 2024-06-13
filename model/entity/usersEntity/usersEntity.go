package usersEntity

type (
	DetailedUserData struct {
		PersonalData     DetailUser       `json:"personalData"`
		EmploymentData   UserJobDetail    `json:"employmentData"`
		EmergencyContact EmergencyContact `json:"emergencyContact"`
	}

	DetailUser struct {
		NIK         string `json:"nik"`
		FullName    string `json:"fullName"`
		PhoneNumber string `json:"phoneNumber"`
		Address     string `json:"address"`
		City        string `json:"city"`
		FotoKtp     string `json:"fotoKTP"`
		FotoSelfie  string `json:"fotoSelfie"`
	}

	UserJobDetail struct {
		JobName       string  `json:"jobName"`
		Salary        float64 `json:"salary"`
		OfficeName    string  `json:"officeName"`
		OfficeContact string  `json:"officeContact"`
		OfficeAddress string  `json:"officeAddress"`
	}

	EmergencyContact struct {
		ContactName string `json:"contactName"`
		PhoneNumber string `json:"emergencyNumber"`
	}
)
