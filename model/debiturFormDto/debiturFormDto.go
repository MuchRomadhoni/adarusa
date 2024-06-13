package debiturFormDto

import "fp_pinjaman_online/model/dto/json"

type ( 
    Debitur struct {
    DetailUser      DetailDebitur   `json:"detail_user"`
    UserJobs        UserJobs        `json:"user_job_detail"`
    EmergencyContact EmergencyContact `json:"emergency"`
    }

    DetailDebitur struct {
        UserID      int    `json:"user_id,omitempty"`
        LimitID     int    `json:"limit_id" binding:"omitempty"`
        Nik         string `json:"nik" binding:"required,numeric,len=16"`
        Fullname    string `json:"fullname" binding:"required"`
        PhoneNumber string `json:"phone_number" binding:"required,numeric"`
        Address     string `json:"address" binding:"required"`
        City        string `json:"city" binding:"required"`
        FotoKtp     string `json:"foto_ktp"`
        FotoSelfie  string `json:"foto_selfie"`
    }

    UserJobs struct {
        UserID      int    `json:"user_id,omitempty"`
        JobName      string `json:"job_name" binding:"required"`
        Salary       float64 `json:"salary" binding:"required,number"`
        OfficeName   string `json:"office_name" binding:"required"`
        OfficeContact string `json:"office_contact" binding:"required"`
        OfficeAddress string `json:"office_address" binding:"required"`
    }

    EmergencyContact struct {
        UserID      int    `json:"user_id,omitempty"`
        Name        string `json:"name" binding:"required"`
        PhoneNumber string `json:"phone_number" binding:"required,numeric"`
    }

    Response struct {
        ResponseCode int `json:"responseCode"`
        Data []DetailDebitur `json:"data"`
        Paging json.Paging `json:"paging"`
    }
)