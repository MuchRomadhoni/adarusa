package mymock

import (
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	"fp_pinjaman_online/model/entity/usersEntity"

	"github.com/stretchr/testify/mock"
)

type DebtCollectorUseCaseMock struct {
	Mock mock.Mock
}

func (um *DebtCollectorUseCaseMock) CreateLogTugas(newLogPayload debtCollectorDto.NewLogTugasPayload) error {
	args := um.Mock.Called(newLogPayload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (um *DebtCollectorUseCaseMock) GetLogTugasById(logTugasId string) (debtCollectorEntity.LogTugas, error) {
	args := um.Mock.Called(logTugasId)
	if args.Get(1) != nil {
		return debtCollectorEntity.LogTugas{}, args.Error(1)
	}
	return args.Get(0).(debtCollectorEntity.LogTugas), nil
}

func (um *DebtCollectorUseCaseMock) EditLogTugasById(logTugasId string, payload debtCollectorDto.UpdateLogTugasPayload) error {
	args := um.Mock.Called(logTugasId, payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (um *DebtCollectorUseCaseMock) LogTugasAuthorizationCheck(logTugasId string) (debtCollectorEntity.LogTugas, error) {
	args := um.Mock.Called(logTugasId)
	if args.Get(1) != nil {
		return debtCollectorEntity.LogTugas{}, args.Error(1)
	}
	return args.Get(0).(debtCollectorEntity.LogTugas), nil
}

func (um *DebtCollectorUseCaseMock) DeleteLogTugasById(logTugasId string) error {
	args := um.Mock.Called(logTugasId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (um *DebtCollectorUseCaseMock) GetAllLogTugas(tugasId string, page, size int) ([]debtCollectorEntity.LogTugas, json.Paging, error) {
	args := um.Mock.Called(tugasId, page, size)
	if args.Get(2) != nil {
		return nil, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.LogTugas), args.Get(1).(json.Paging), nil
}

func (um *DebtCollectorUseCaseMock) GetAllLateDebtorByCity(dcId string, page, size int) ([]debtCollectorEntity.LateDebtor, json.Paging, error) {
	args := um.Mock.Called(dcId, page, size)
	if args.Get(2) != nil {
		return nil, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.LateDebtor), args.Get(1).(json.Paging), nil
}

func (um *DebtCollectorUseCaseMock) ClaimTugas(dcId string, payload debtCollectorDto.NewTugasPayload) error {
	args := um.Mock.Called(dcId, payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (um *DebtCollectorUseCaseMock) GetAllTugas(dcId, status string, page, size int) ([]debtCollectorEntity.Tugas, json.Paging, error) {
	args := um.Mock.Called(dcId, status, page, size)
	if args.Get(2) != nil {
		return nil, json.Paging{}, args.Error(2)
	}
	return args.Get(0).([]debtCollectorEntity.Tugas), args.Get(1).(json.Paging), nil
}

func (um *DebtCollectorUseCaseMock) GetBalanceByUserId(userId string) (float64, error) {
	args := um.Mock.Called(userId)
	if args.Get(1) != nil {
		return 0, args.Error(1)
	}
	return float64(args.Int(0)), nil
}

func (um *DebtCollectorUseCaseMock) CreateWithdrawRequest(userId string, amount float64) error {
	args := um.Mock.Called(userId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (um *DebtCollectorUseCaseMock) GetDebtorData(userId, dcId string) (usersEntity.DetailedUserData, error) {
	args := um.Mock.Called(userId, dcId)
	if args.Get(1) != nil {
		return usersEntity.DetailedUserData{}, args.Error(1)
	}
	return args.Get(0).(usersEntity.DetailedUserData), nil
}
