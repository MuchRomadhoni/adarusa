package mymock

import (
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"

	"github.com/stretchr/testify/mock"
)

type DebtCollectorRepositoryMock struct {
	Mock mock.Mock
}

func (dm *DebtCollectorRepositoryMock) SelectTugasById(tugasId string) (debtCollectorEntity.Tugas, error) {
	args := dm.Mock.Called(tugasId)
	if args.Get(1) != nil {
		return debtCollectorEntity.Tugas{}, args.Error(1)
	}
	return args.Get(0).(debtCollectorEntity.Tugas), nil
}

func (dm *DebtCollectorRepositoryMock) InsertLogTugas(newLogPayload debtCollectorDto.NewLogTugasPayload) error {
	args := dm.Mock.Called(newLogPayload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (dm *DebtCollectorRepositoryMock) UpdateLogTugasById(storedLog debtCollectorEntity.LogTugas, updateLogPayload debtCollectorDto.UpdateLogTugasPayload) error {
	args := dm.Mock.Called(storedLog, updateLogPayload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (dm *DebtCollectorRepositoryMock) SelectLogTugasById(logTugasId string) (debtCollectorEntity.LogTugas, error) {
	args := dm.Mock.Called(logTugasId)
	if args.Get(1) != nil {
		return debtCollectorEntity.LogTugas{}, args.Error(1)
	}
	return args.Get(0).(debtCollectorEntity.LogTugas), nil
}

func (dm *DebtCollectorRepositoryMock) SoftDeleteLogTugasById(logTugasId string) error {
	args := dm.Mock.Called(logTugasId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (dm *DebtCollectorRepositoryMock) SelectAllLogByTugasId(tugasId string, page, size int) ([]debtCollectorEntity.LogTugas, json.Paging, error) {
	args := dm.Mock.Called(tugasId, page, size)
	if args.Get(2) != nil {
		return []debtCollectorEntity.LogTugas{}, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.LogTugas), args.Get(1).(json.Paging), nil
}

func (dm *DebtCollectorRepositoryMock) SelectAllLateDebitur(dcCity string, page, size int) ([]debtCollectorEntity.LateDebtor, json.Paging, error) {
	args := dm.Mock.Called(dcCity, page, size)
	if args.Get(2) != nil {
		return nil, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.LateDebtor), args.Get(1).(json.Paging), nil
}

func (dm *DebtCollectorRepositoryMock) SelectDebtCollectorById(id string) (debtCollectorEntity.DebtCollector, error) {
	args := dm.Mock.Called(id)
	if args.Get(1) != nil {
		return debtCollectorEntity.DebtCollector{}, args.Error(1)
	}
	return args.Get(0).(debtCollectorEntity.DebtCollector), nil
}

func (dm *DebtCollectorRepositoryMock) SelectLateDebiturById(userId, dcCity string) (string, error) {
	args := dm.Mock.Called(userId, dcCity)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}
	return args.String(0), nil
}

func (dm *DebtCollectorRepositoryMock) CreateClaimTugas(dcId, userId string) error {
	args := dm.Mock.Called(dcId, userId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (dm *DebtCollectorRepositoryMock) SelectAllTugas(dcId, status string, page, size int) ([]debtCollectorEntity.Tugas, json.Paging, error) {
	args := dm.Mock.Called(dcId, status, page, size)
	if args.Get(2) != nil {
		return nil, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.Tugas), args.Get(1).(json.Paging), nil
}

func (dm *DebtCollectorRepositoryMock) CountOngoingTugas(dcId string) (int, error) {
	args := dm.Mock.Called(dcId)
	if args.Get(1) != nil {
		return 0, args.Error(1)
	}
	return args.Int(0), nil
}

func (dm *DebtCollectorRepositoryMock) SelectBalanceByUserId(userId string) (float64, error) {
	args := dm.Mock.Called(userId)
	if args.Get(1) != nil {
		return 0, args.Error(1)
	}
	return float64(args.Int(0)), nil
}

func (dm *DebtCollectorRepositoryMock) CreateWithdrawRequest(userId string, amount float64) error {
	args := dm.Mock.Called(userId, amount)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (dm *DebtCollectorRepositoryMock) SelectDebtorFromTugas(dcId, userId string) (string, error) {
	args := dm.Mock.Called(dcId, userId)
	if args.Get(1) != nil {
		return "", nil
	}
	return args.String(0), nil
}
