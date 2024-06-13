package userUseCase

import (
	"fp_pinjaman_online/model/dcFormDto"
	"fp_pinjaman_online/model/debiturFormDto"
	"fp_pinjaman_online/model/entity/usersEntity"
	"fp_pinjaman_online/model/userDto"
	"fp_pinjaman_online/src/users/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Hasher interface {
	HashedPassword(password string) (string, error)
}

func TestLogin_Success(t *testing.T) {
	// Initialize the mock repository
	mockRepo := new(mocks.UserRepository)

	// Initialize the use case with the mock repository
	useCase := &userUC{
		userRepo: mockRepo,
	}

	// Define the input and expected output
	req := userDto.LoginRequest{Email: "test@example.com", Password: "password"}
	expectedUser := userDto.User{Id: "1", Email: "test@example.com", Roles: "user"}

	// Setup the expectations
	mockRepo.On("Login", req).Return("hashedPassword", nil)
	mockRepo.On("GetUserByEmail", req.Email).Return(expectedUser, nil)

	// Call the function
	token, err := useCase.Login(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and epected output
	req := userDto.CreateRequest{Name:"",  Email: "test@example.com", Password: "password", Roles: 1}
	roleId := 1

	// setup the expectations
	mockRepo.On("UserExists", req.Email).Return(false, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("userDto.CreateRequest"), roleId).Return(nil)

	// call the function
	err := useCase.CreateUser(req, roleId)
	
	// assertions
	assert.NoError(t, err)

	// check thta the expectation were met
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and expected output
	email := "test@example.com"
	expectedUser := userDto.User{Id: "1", Email: "test@example.com", Roles: "user"}

	// setup the expectations
	mockRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

	// call the function
	user, err := useCase.GetUserByEmail(email)

	// assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateDetailDebitur_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input
	debitur := debiturFormDto.Debitur{
		DetailUser: debiturFormDto.DetailDebitur{
			UserID:      1,
			LimitID:     1,
			Nik:         "1234567890",
			Fullname:    "Test User",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
			City:        "Test City",
			FotoKtp:     "test_ktp.jpg",
			FotoSelfie:  "test_selfie.jpg",
		},
		UserJobs: debiturFormDto.UserJobs{
			UserID:        1,
			JobName:       "Test Job",
			Salary:        5000.0,
			OfficeName:    "Test Office",
			OfficeContact: "1234567890",
			OfficeAddress: "Test Office Address",
		},
		EmergencyContact: debiturFormDto.EmergencyContact{
			UserID:      1,
			Name:        "Test Emergency Contact",
			PhoneNumber: "1234567890",
		},
	}

	// setup the expectations
	mockRepo.On("CreateDetailDebitur", debitur).Return(nil)

	// call the function
	err := useCase.CreateDetailDebitur(debitur)

	// assertions
	assert.NoError(t, err)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateDetailDc_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input
	dc := dcFormDto.DetailDC{
		UserID:      1,
		LimitID:     1,
		Nik:         "1234567890",
		Fullname:    "Test User",
		PhoneNumber: "1234567890",
		Address:     "Test Address",
		City:        "Test City",
		FotoKtp:     "test_ktp.jpg",
		FotoSelfie:  "test_selfie.jpg",
	}

	// setup the expectations
	mockRepo.On("CreateDetailDc", dc).Return(nil)

	// call the function
	err := useCase.CreateDetailDc(dc)

	// assertions
	assert.NoError(t, err)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetFullname_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and expected output
	userId := 1
	expectedFullname := "Test User"

	// setup the expectations
	mockRepo.On("GetFullname", userId).Return(expectedFullname, nil)

	// call the function
	fullname, err := useCase.GetFullname(userId)

	// assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedFullname, fullname)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetUserDataById_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and expected output
	userId := "1"
	expectedRoleId := "1"
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

	// setup the expectations
	mockRepo.On("GetRolesById", userId).Return(expectedRoleId, nil)
	mockRepo.On("GetUserDetailByUserId", userId).Return(expectedDetailData, nil)
	mockRepo.On("GetUserJobDetailByUserId", userId).Return(expectedJobData, nil)
	mockRepo.On("GetEmergencyContactByUserId", userId).Return(expectedEmergencyContact, nil)

	// call the function
	userData, err := useCase.GetUserDataById(userId)

	// assertions
	assert.NoError(t, err)

	detailUserData, ok := userData.(usersEntity.DetailedUserData)
	assert.True(t, ok, "Expected userData to be of type DetailedUserData")
	
	assert.Equal(t, expectedDetailData, detailUserData.PersonalData)
	assert.Equal(t, expectedJobData, detailUserData.EmploymentData)
	assert.Equal(t, expectedEmergencyContact, detailUserData.EmergencyContact)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdatePhotoPaths_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input
	userId := 1
	fotoKTP := "path/to/ktp.jpg"
	fotoSelfie := "path/to/selfie.jpg"

	// setup the expectations
	mockRepo.On("UpdatePhotoPaths", userId, fotoKTP, fotoSelfie).Return(nil)

	// call the function
	err := useCase.UpdatePhotoPaths(userId, fotoKTP, fotoSelfie)

	// assertions
	assert.NoError(t, err)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetDataByRole_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and expected output
	role := "testRole"
	status := "testStatus"
	page := 1
	size := 10
	expectedDebitur := []debiturFormDto.DetailDebitur{
		{
			UserID:      1,
			LimitID:     1,
			Nik:         "1234567890",
			Fullname:    "Test User",
			PhoneNumber: "1234567890",
			Address:     "Test Address",
			City:        "Test City",
			FotoKtp:     "test_ktp.jpg",
			FotoSelfie:  "test_selfie.jpg",
		},
	}
	expectedTotalData := 100

	// setup the expectations
	mockRepo.On("GetDataByRole", role, status, size, (page-1)*size).Return(expectedDebitur, expectedTotalData, nil)

	// call the function
	debitur, totalData, err := useCase.GetDataByRole(role, status, page, size)

	// assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedDebitur, debitur)
	assert.Equal(t, expectedTotalData, totalData)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdateBankAccount_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input
	userId := 1
	accountNumber := "1234567890"
	bankName := "Test Bank"

	// setup the expectations
	mockRepo.On("IsBankAccExist", userId, accountNumber).Return(false, nil)
	mockRepo.On("UpdateBankAccount", userId, accountNumber, bankName).Return(nil)

	// call the function
	err := useCase.UpdateBankAccount(userId, accountNumber, bankName)

	// assertions
	assert.NoError(t, err)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestIsBankAccExist_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	useCase := &userUC{
		userRepo: mockRepo,
	}

	// define the input and expected output
	userId := 1
	accountNumber := "1234567890"
	expectedExists := true

	// setup the expectations
	mockRepo.On("IsBankAccExist", userId, accountNumber).Return(expectedExists, nil)

	// call the function
	exists, err := useCase.IsBankAccExist(userId, accountNumber)

	// assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedExists, exists)

	// check that the expectations were met
	mockRepo.AssertExpectations(t)
}