package debiturUsecase_test

import (
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/model/dto/debiturDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/src/debitur/debiturUsecase"
	"fp_pinjaman_online/src/midtrans"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DebiturRepositoryMock struct {
	mock.Mock
}

func (m *DebiturRepositoryMock) AddPengajuanPinjaman(id int, jumlahPinjaman float64, tenor int, description string) error {
	args := m.Called(id, jumlahPinjaman, tenor, description)
	return args.Error(0)
}

func (m *DebiturRepositoryMock) GetPengajuanPinjaman(id int) ([]debiturDto.GetPengajuanResponse, error) {
	args := m.Called(id)
	return args.Get(0).([]debiturDto.GetPengajuanResponse), args.Error(1)
}

func (m *DebiturRepositoryMock) GetCicilan(page, limit, offset int, id string, status string) ([]debiturDto.GetCicilanResponse, json.Paging, error) {
	args := m.Called(page, limit, offset, id, status)
	return args.Get(0).([]debiturDto.GetCicilanResponse), args.Get(1).(json.Paging), args.Error(2)
}

func (m *DebiturRepositoryMock) CicilanPayment(pinjamanId int, totalBayar float64) (dto.MidtransSnapResponse, error) {
	args := m.Called(pinjamanId, totalBayar)
	return args.Get(0).(dto.MidtransSnapResponse), args.Error(1)
}

func (m *DebiturRepositoryMock) CicilanVerify(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *DebiturRepositoryMock) UpdatePinjamanStatus(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *DebiturRepositoryMock) SetMidtransService(service midtrans.MidtransService) {
	m.Called(service)
}

func TestPengajuanPinjaman(t *testing.T) {
	mockRepo := new(DebiturRepositoryMock)
	usecase := debiturUsecase.NewDebiturUsecase(mockRepo)

	mockRepo.On("AddPengajuanPinjaman", 1, 1000.0, 6, "Test description").Return(nil)

	err := usecase.PengajuanPinjaman(1, 1000.0, 6, "Test description")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetPengajuanPinjaman(t *testing.T) {
	mockRepo := new(DebiturRepositoryMock)
	usecase := debiturUsecase.NewDebiturUsecase(mockRepo)

	mockData := []debiturDto.GetPengajuanResponse{
		{Id: 1, JumlahPinjaman: 1000.0, Tenor: 6, Description: "Test description", BungaPerBulan: 0.09, StatusPengajuan: "pending"},
	}
	mockRepo.On("GetPengajuanPinjaman", 1).Return(mockData, nil)

	data, err := usecase.GetPengajuanPinjaman(1)
	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	mockRepo.AssertExpectations(t)
}

func TestGetCicilan(t *testing.T) {
	mockRepo := new(DebiturRepositoryMock)
	usecase := debiturUsecase.NewDebiturUsecase(mockRepo)

	tanggalJatuhTempo, _ := time.Parse("2006-01-02", "2022-01-01")
	tanggalBayarSelesai, _ := time.Parse("2006-01-02", "2022-01-01")

	mockData := []debiturDto.GetCicilanResponse{
		{
			Id:                  1,
			PinjamanId:          1,
			TanggalJatuhTempo:   tanggalJatuhTempo,
			TanggalBayarSelesai: tanggalBayarSelesai,
			JumlahBayar:         1000.0,
			Status:              "unpaid",
		},
	}

	mockPaging := json.Paging{Page: 1, TotalData: 1}
	mockRepo.On("GetCicilan", 1, 10, 0, "1", "").Return(mockData, mockPaging, nil)

	data, paging, err := usecase.GetCicilan(1, 10, 0, "1", "")
	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, mockPaging, paging)
	mockRepo.AssertExpectations(t)
}

func TestCicilanPayment(t *testing.T) {
	mockRepo := new(DebiturRepositoryMock)
	usecase := debiturUsecase.NewDebiturUsecase(mockRepo)

	mockData := dto.MidtransSnapResponse{
		Token:        "token",
		RedirectUrl:  "https://example.com",
		ErrorMessage: nil,
	}

	mockRepo.On("CicilanPayment", 1, 1000.0).Return(mockData, nil)

	_, err := usecase.CicilanPayment(1, 1000.0)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCicilanVerify(t *testing.T) {
	mockRepo := new(DebiturRepositoryMock)
	usecase := debiturUsecase.NewDebiturUsecase(mockRepo)

	mockRepo.On("CicilanVerify", 1).Return(nil)

	mockRepo.On("UpdatePinjamanStatus", 1).Return(nil)

	err := usecase.CicilanVerify(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
