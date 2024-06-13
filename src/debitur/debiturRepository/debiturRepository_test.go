package debiturRepository_test

import (
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/src/debitur/debiturRepository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddPengajuanPinjaman(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	query := regexp.QuoteMeta("INSERT INTO pinjaman (user_id, jumlah_pinjaman, tenor, description, bunga_per_bulan, status_pengajuan) VALUES ($1, $2, $3, $4, $5, 'pending')")
	mock.ExpectExec(query).WithArgs(1, 1000.0, 6, "Test description", 0.09).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddPengajuanPinjaman(1, 1000.0, 6, "Test description")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPengajuanPinjaman(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	query := regexp.QuoteMeta("SELECT id, jumlah_pinjaman, tenor, description, bunga_per_bulan, status_pengajuan, created_at, updated_at FROM pinjaman WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC")
	rows := sqlmock.NewRows([]string{"id", "jumlah_pinjaman", "tenor", "description", "bunga_per_bulan", "status_pengajuan", "created_at", "updated_at"}).
		AddRow(1, 1000.0, 6, "Test description", 0.09, "pending", time.Now(), time.Now())

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	result, err := repo.GetPengajuanPinjaman(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCicilan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	queryCount := regexp.QuoteMeta("SELECT COUNT(*) FROM cicilan WHERE pinjaman_id = $1")
	querySelect := regexp.QuoteMeta("SELECT id, pinjaman_id, tanggal_jatuh_tempo, jumlah_bayar, status FROM cicilan WHERE pinjaman_id = $1")

	mock.ExpectQuery(queryCount).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	rows := sqlmock.NewRows([]string{"id", "pinjaman_id", "tanggal_jatuh_tempo", "jumlah_bayar", "status"}).
		AddRow(1, 1, time.Now(), 1000.0, "unpaid")
	mock.ExpectQuery(querySelect).WithArgs("1").WillReturnRows(rows)

	result, paging, err := repo.GetCicilan(1, 10, 0, "1", "")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, paging.TotalData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCicilanPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	querySelect := regexp.QuoteMeta("SELECT id, jumlah_bayar FROM cicilan WHERE pinjaman_id = $1 AND status = 'unpaid'")
	mock.ExpectQuery(querySelect).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "jumlah_bayar"}).AddRow(1, 1000.0))

	queryCustomer := regexp.QuoteMeta("SELECT d.fullname FROM pinjaman p JOIN detail_users d ON p.user_id = d.user_id WHERE p.id = $1")
	mock.ExpectQuery(queryCustomer).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"fullname"}).AddRow("John Doe"))

	mock.ExpectExec("INSERT INTO midtrans_tx").WillReturnResult(sqlmock.NewResult(1, 1))

	data, err := repo.CicilanPayment(1, 1000.0)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// MockMidtransService adalah mock untuk midtrans.MidtransService
type MockMidtransService struct {
	mock.Mock
}

func (m *MockMidtransService) Pay(payload dto.MidtransSnapRequest) (dto.MidtransSnapResponse, error) {
	args := m.Called(payload)
	return args.Get(0).(dto.MidtransSnapResponse), args.Error(1)
}

func (m *MockMidtransService) VerifyPayment(orderId int) (bool, error) {
	args := m.Called(orderId)
	return args.Bool(0), args.Error(1)
}

func TestCicilanVerify_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock midtrans service
	mockMidtrans := new(MockMidtransService)
	mockMidtrans.On("VerifyPayment", 1).Return(true, nil)

	repo := debiturRepository.NewDebiturRepository(db)
	repo.SetMidtransService(mockMidtrans)

	// Mock database queries
	mock.ExpectExec("UPDATE cicilan SET status = 'paid', tanggal_selesai_bayar = NOW\\(\\) WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE midtrans_tx SET status = 'success' WHERE cicilan_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT pinjaman_id FROM cicilan WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"pinjaman_id"}).AddRow(1))

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM cicilan WHERE pinjaman_id = \\$1 AND status = 'unpaid'\\)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("UPDATE pinjaman SET status_pengajuan = 'completed', updated_at = NOW\\(\\) WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.CicilanVerify(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCicilanVerify_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock midtrans service
	mockMidtrans := new(MockMidtransService)
	mockMidtrans.On("VerifyPayment", 1).Return(false, nil)

	repo := debiturRepository.NewDebiturRepository(db)
	repo.SetMidtransService(mockMidtrans)

	// Call the method
	err = repo.CicilanVerify(1)
	assert.Error(t, err)
	assert.Equal(t, "payment not success", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePinjamanStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	querySelect := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM cicilan WHERE pinjaman_id = $1 AND status = 'unpaid')")
	mock.ExpectQuery(querySelect).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	err = repo.UpdatePinjamanStatus(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePinjamanStatus_Completed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := debiturRepository.NewDebiturRepository(db)

	querySelect := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM cicilan WHERE pinjaman_id = $1 AND status = 'unpaid')")
	mock.ExpectQuery(querySelect).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("UPDATE pinjaman SET status_pengajuan = 'completed'").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdatePinjamanStatus(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
