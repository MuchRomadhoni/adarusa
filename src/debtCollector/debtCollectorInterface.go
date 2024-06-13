package debtCollector

import (
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	"fp_pinjaman_online/model/entity/usersEntity"
)

type DebtCollectorUseCase interface {
	CreateLogTugas(newLogPayload debtCollectorDto.NewLogTugasPayload) error
	GetLogTugasById(logTugasId string) (debtCollectorEntity.LogTugas, error)
	EditLogTugasById(logTugasId string, payload debtCollectorDto.UpdateLogTugasPayload) error
	LogTugasAuthorizationCheck(logTugasId string) (debtCollectorEntity.LogTugas, error)
	DeleteLogTugasById(logTugasId string) error
	GetAllLogTugas(tugasId string, page, size int) ([]debtCollectorEntity.LogTugas, json.Paging, error)
	GetAllLateDebtorByCity(dcId string, page, size int) ([]debtCollectorEntity.LateDebtor, json.Paging, error)
	ClaimTugas(dcId string, payload debtCollectorDto.NewTugasPayload) error
	GetAllTugas(dcId, status string, page, size int) ([]debtCollectorEntity.Tugas, json.Paging, error)
	GetBalanceByUserId(userId string) (float64, error)
	CreateWithdrawRequest(userId string, amount float64) error
	GetDebtorData(userId, dcId string) (usersEntity.DetailedUserData, error)
}

type DebtCollectorRepository interface {
	SelectTugasById(tugasId string) (debtCollectorEntity.Tugas, error)
	InsertLogTugas(newLogPayload debtCollectorDto.NewLogTugasPayload) error
	UpdateLogTugasById(storedLog debtCollectorEntity.LogTugas, updateLogPayload debtCollectorDto.UpdateLogTugasPayload) error
	SelectLogTugasById(logTugasId string) (debtCollectorEntity.LogTugas, error)
	SoftDeleteLogTugasById(logTugasId string) error
	SelectAllLogByTugasId(tugasId string, page, size int) ([]debtCollectorEntity.LogTugas, json.Paging, error)
	SelectAllLateDebitur(dcCity string, page, size int) ([]debtCollectorEntity.LateDebtor, json.Paging, error)
	SelectLateDebiturById(userId, dcCity string) (string, error)
	CreateClaimTugas(dcId, userId string) error
	SelectAllTugas(dcId, status string, page, size int) ([]debtCollectorEntity.Tugas, json.Paging, error)
	CountOngoingTugas(dcId string) (int, error)
	SelectBalanceByUserId(userId string) (float64, error)
	CreateWithdrawRequest(userId string, amount float64) error
	SelectDebtorFromTugas(dcId, userId string) (string, error)
	SelectDebtCollectorById(id string) (debtCollectorEntity.DebtCollector, error) // TODO : it should be in users repository
}
