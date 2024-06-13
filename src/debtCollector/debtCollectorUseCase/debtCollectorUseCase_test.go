package debtCollectorUseCase

import (
	"errors"
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	"fp_pinjaman_online/model/entity/usersEntity"
	"fp_pinjaman_online/src/debtCollector"
	mymock "fp_pinjaman_online/src/debtCollector/debtCollectorRepository/mock"
	"fp_pinjaman_online/src/users/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DebtCollectorUseCaseSuite struct {
	suite.Suite
	dcRepoMock   *mymock.DebtCollectorRepositoryMock
	userRepoMock *mocks.UserRepository
	usecase      debtCollector.DebtCollectorUseCase
}

func (s *DebtCollectorUseCaseSuite) SetupTest() {
	s.dcRepoMock = &mymock.DebtCollectorRepositoryMock{Mock: mock.Mock{}}
	s.userRepoMock = &mocks.UserRepository{}
	s.usecase = NewDebtCollectorUseCase(s.dcRepoMock, s.userRepoMock)
}

func TestUsersUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(DebtCollectorUseCaseSuite))
}

func (s *DebtCollectorUseCaseSuite) TestLogTugasAuthorizationCheck_Success() {
	logTugasIdMock := "1"

	expLog := debtCollectorEntity.LogTugas{}

	s.dcRepoMock.Mock.On("SelectLogTugasById", logTugasIdMock).Return(expLog, nil)

	log, err := s.usecase.LogTugasAuthorizationCheck(logTugasIdMock)

	s.NoError(err)
	s.Equal(expLog, log)
}

func (s *DebtCollectorUseCaseSuite) TestGetAllLateDebtorByCity_Success() {
	dcIdMock := "1"
	pageMock := 1
	sizeMock := 5

	expLoggedDc := debtCollectorEntity.DebtCollector{
		City: "malang",
	}
	expLateDebtorList := []debtCollectorEntity.LateDebtor{}
	expPagingMock := json.Paging{}

	s.dcRepoMock.Mock.On("SelectDebtCollectorById", dcIdMock).Return(expLoggedDc, nil)
	s.dcRepoMock.Mock.On("SelectAllLateDebitur", expLoggedDc.City, pageMock, sizeMock).Return(expLateDebtorList, expPagingMock, nil)

	lateDebtorList, paging, err := s.usecase.GetAllLateDebtorByCity(dcIdMock, pageMock, sizeMock)

	s.NoError(err)
	s.Equal(expPagingMock, paging)
	s.Equal(expLateDebtorList, lateDebtorList)
}

func (s *DebtCollectorUseCaseSuite) TestCreateLogTugas_Success() {
	mockPayload := debtCollectorDto.NewLogTugasPayload{
		TugasId:     "1",
		Description: "tes isi log",
	}
	s.dcRepoMock.Mock.On("SelectTugasById", mockPayload.TugasId).Return(debtCollectorEntity.Tugas{}, nil)
	s.dcRepoMock.Mock.On("InsertLogTugas", mockPayload).Return(nil)

	err := s.usecase.CreateLogTugas(mockPayload)

	s.NoError(err)
}

func (s *DebtCollectorUseCaseSuite) TestCreateLogTugas_Fail() {
	s.Run("Fail tugas not found", func() {
		mockPayloadNotFound := debtCollectorDto.NewLogTugasPayload{
			TugasId:     "1",
			Description: "tes isi log",
		}

		expectedError := errors.New("error not found")
		s.dcRepoMock.Mock.On("SelectTugasById", mockPayloadNotFound.TugasId).Return(debtCollectorEntity.Tugas{}, expectedError)

		err := s.usecase.CreateLogTugas(mockPayloadNotFound)
		s.Error(err)
		s.Equal(expectedError, err)
	})

	s.Run("Fail internal server error", func() {
		mockPayloadInternalServerError := debtCollectorDto.NewLogTugasPayload{
			TugasId:     "2",
			Description: "tes isi log",
		}

		expectedErrorServer := errors.New("error internal server")
		s.dcRepoMock.Mock.On("SelectTugasById", mockPayloadInternalServerError.TugasId).Return(debtCollectorEntity.Tugas{}, nil)
		s.dcRepoMock.Mock.On("InsertLogTugas", mockPayloadInternalServerError).Return(expectedErrorServer)

		err := s.usecase.CreateLogTugas(mockPayloadInternalServerError)
		s.Error(err)
		s.Equal(expectedErrorServer, err)
	})
}

func (s *DebtCollectorUseCaseSuite) TestGetLogTugasById_Success() {
	mockLogTugasId := "1"
	expectedLogTugas := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "tugas",
	}

	s.dcRepoMock.Mock.On("SelectLogTugasById", mockLogTugasId).Return(expectedLogTugas, nil)

	log, err := s.usecase.GetLogTugasById(mockLogTugasId)

	s.NoError(err)
	s.Equal(expectedLogTugas, log)
}

func (s *DebtCollectorUseCaseSuite) TestEditLogTugasById_Success() {
	mockLogTugasId := "1"
	expectedLogTugas := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "tugas",
	}
	mockPayload := debtCollectorDto.UpdateLogTugasPayload{
		Description: "tugas updated",
	}

	s.dcRepoMock.Mock.On("SelectLogTugasById", mockLogTugasId).Return(expectedLogTugas, nil)
	s.dcRepoMock.Mock.On("UpdateLogTugasById", expectedLogTugas, mockPayload).Return(nil)

	err := s.usecase.EditLogTugasById(mockLogTugasId, mockPayload)

	s.NoError(err)
}

func (s *DebtCollectorUseCaseSuite) TestDeleteLogTugasById_Success() {
	mockLogTugasId := "1"
	expectedLogTugas := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "tugas",
	}

	s.dcRepoMock.Mock.On("SelectLogTugasById", mockLogTugasId).Return(expectedLogTugas, nil)
	s.dcRepoMock.Mock.On("SoftDeleteLogTugasById", mockLogTugasId).Return(nil)

	err := s.usecase.DeleteLogTugasById(mockLogTugasId)

	s.NoError(err)
}

func (s *DebtCollectorUseCaseSuite) TestGetAllLogTugas_Success() {
	mockTugasId := "1"
	expLogList := []debtCollectorEntity.LogTugas{{}}
	expPaging := json.Paging{}

	s.dcRepoMock.Mock.On("SelectAllLogByTugasId", mockTugasId, 0, 0).Return(expLogList, expPaging, nil)

	logList, paging, err := s.usecase.GetAllLogTugas(mockTugasId, 0, 0)

	s.NoError(err)
	s.Equal(expLogList, logList)
	s.Equal(expPaging, paging)
}

func (s *DebtCollectorUseCaseSuite) TestClaimTugas_Success() {
	dcIdMock := "1"
	mockPayload := debtCollectorDto.NewTugasPayload{
		UserId: "1",
	}
	expLoggedDc := debtCollectorEntity.DebtCollector{
		City: "malang",
	}

	s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(2, nil)
	s.dcRepoMock.Mock.On("SelectDebtCollectorById", dcIdMock).Return(expLoggedDc, nil)
	s.dcRepoMock.Mock.On("SelectLateDebiturById", mockPayload.UserId, expLoggedDc.City).Return("", nil)
	s.dcRepoMock.Mock.On("CreateClaimTugas", dcIdMock, mockPayload.UserId).Return(nil)

	err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

	s.NoError(err)
}

func (s *DebtCollectorUseCaseSuite) TestClaimTugas_Fail() {
	dcIdMock := "1"
	mockPayload := debtCollectorDto.NewTugasPayload{
		UserId: "1",
	}
	s.Run("Error on CountOngoingTugas", func() {
		expCountError := errors.New("count error")
		s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(0, expCountError)

		err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

		s.Error(err)
		s.Equal(expCountError.Error(), err.Error())
		s.dcRepoMock.Mock.AssertExpectations(s.T())
		s.dcRepoMock.Mock.ExpectedCalls = nil
	})

	s.Run("Error max tugas exceed", func() {
		expMaxError := errors.New("maximum ongoing tax is 3")
		s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(3, nil)

		err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

		s.Error(err)
		s.Equal(expMaxError.Error(), err.Error())
		s.dcRepoMock.Mock.AssertExpectations(s.T())
		s.dcRepoMock.Mock.ExpectedCalls = nil
	})

	s.Run("Error on SelectDebtCollectorById", func() {
		expError := errors.New("dc not found")
		s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(2, nil)
		s.dcRepoMock.Mock.On("SelectDebtCollectorById", dcIdMock).Return(debtCollectorEntity.DebtCollector{}, expError)

		err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

		s.Error(err)
		s.Equal(expError.Error(), err.Error())
		s.dcRepoMock.Mock.AssertExpectations(s.T())
		s.dcRepoMock.Mock.ExpectedCalls = nil
	})

	s.Run("Error on SelectLateDebiturById", func() {
		expLoggedDc := debtCollectorEntity.DebtCollector{
			City: "malang",
		}
		expError := errors.New("debitur not found")
		s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(2, nil)
		s.dcRepoMock.Mock.On("SelectDebtCollectorById", dcIdMock).Return(expLoggedDc, nil)
		s.dcRepoMock.Mock.On("SelectLateDebiturById", mockPayload.UserId, expLoggedDc.City).Return("", expError)

		err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

		s.Error(err)
		s.Equal(expError.Error(), err.Error())
		s.dcRepoMock.Mock.AssertExpectations(s.T())
		s.dcRepoMock.Mock.ExpectedCalls = nil
	})

	s.Run("Error on CreateClaimTugas", func() {
		expLoggedDc := debtCollectorEntity.DebtCollector{
			City: "malang",
		}
		expError := errors.New("error sql")
		s.dcRepoMock.Mock.On("CountOngoingTugas", dcIdMock).Return(2, nil)
		s.dcRepoMock.Mock.On("SelectDebtCollectorById", dcIdMock).Return(expLoggedDc, nil)
		s.dcRepoMock.Mock.On("SelectLateDebiturById", mockPayload.UserId, expLoggedDc.City).Return("", nil)
		s.dcRepoMock.Mock.On("CreateClaimTugas", dcIdMock, mockPayload.UserId).Return(expError)

		err := s.usecase.ClaimTugas(dcIdMock, mockPayload)

		s.Error(err)
		s.Equal(expError.Error(), err.Error())
		s.dcRepoMock.Mock.AssertExpectations(s.T())
	})
}

func (s *DebtCollectorUseCaseSuite) TestGetAllTugas_Success() {
	dcIdMock := "1"
	statusMock := ""
	pageMock := 1
	sizeMock := 5

	expListTugas := []debtCollectorEntity.Tugas{}
	expPaging := json.Paging{}

	s.dcRepoMock.Mock.On("SelectAllTugas", dcIdMock, statusMock, pageMock, sizeMock).Return(expListTugas, expPaging, nil)

	listTugas, paging, err := s.usecase.GetAllTugas(dcIdMock, statusMock, pageMock, sizeMock)

	s.NoError(err)
	s.Equal(expPaging, paging)
	s.Equal(expListTugas, listTugas)
}

func (s *DebtCollectorUseCaseSuite) TestCreateWithdrawRequest() {
	userIdMock := "1"
	amountMock := 2000000

	s.dcRepoMock.Mock.On("SelectBalanceByUserId", userIdMock).Return(3000000, nil)
	s.dcRepoMock.Mock.On("CreateWithdrawRequest", userIdMock, float64(amountMock)).Return(nil)

	err := s.usecase.CreateWithdrawRequest(userIdMock, float64(amountMock))

	s.NoError(err)
}

func (s *DebtCollectorUseCaseSuite) TestGetDebtorData_Success() {
	userIdMock := "1"
	dcIdMock := "2"
	expDebtorId := "1"
	expectedDetailData := usersEntity.DetailUser{
		NIK:         "1234567890",
		FullName:    "Test User",
		PhoneNumber: "1234567890",
		Address:     "Test Address",
		City:        "Test City",
		FotoKtp:     "test_ktp.jpg",
		FotoSelfie:  "test_selfie.jpg",
	}
	expectedJobData := usersEntity.UserJobDetail{
		JobName:       "Test Job",
		Salary:        5000.0,
		OfficeName:    "Test Office",
		OfficeContact: "1234567890",
		OfficeAddress: "Test Office Address",
	}
	expectedEmergencyContact := usersEntity.EmergencyContact{
		ContactName: "Test Emergency Contact",
		PhoneNumber: "1234567890",
	}
	expResult := usersEntity.DetailedUserData{
		PersonalData:     expectedDetailData,
		EmploymentData:   expectedJobData,
		EmergencyContact: expectedEmergencyContact,
	}

	s.dcRepoMock.Mock.On("SelectDebtorFromTugas", dcIdMock, userIdMock).Return(expDebtorId, nil)
	s.userRepoMock.On("GetUserDetailByUserId", expDebtorId).Return(expectedDetailData, nil)
	s.userRepoMock.On("GetUserJobDetailByUserId", expDebtorId).Return(expectedJobData, nil)
	s.userRepoMock.On("GetEmergencyContactByUserId", expDebtorId).Return(expectedEmergencyContact, nil)

	data, err := s.usecase.GetDebtorData(userIdMock, dcIdMock)

	s.NoError(err)
	s.Equal(expResult, data)
}
