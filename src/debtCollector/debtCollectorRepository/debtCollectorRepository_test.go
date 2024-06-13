package debtCollectorRepository

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	"fp_pinjaman_online/src/debtCollector"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type DebtCollectorRepoTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo debtCollector.DebtCollectorRepository
}

func (s *DebtCollectorRepoTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.repo = NewDebtCollectorRepository(db)
}

func TestDebtCollectorRepoTestSuite(t *testing.T) {
	suite.Run(t, new(DebtCollectorRepoTestSuite))
}

func (s *DebtCollectorRepoTestSuite) TestSelectDebtCollectorById_Success() {
	dcIdMock := "1"
	expDC := debtCollectorEntity.DebtCollector{
		ID:       "1",
		FullName: "ini fullname",
		City:     "malang",
	}

	rows := sqlmock.NewRows([]string{"id", "fullname", "city"}).
		AddRow(expDC.ID, expDC.FullName, expDC.City)
	query := regexp.QuoteMeta(`SELECT u.id,du.fullname,du.city FROM users u JOIN detail_users du ON du.user_id = u.id WHERE u.id = $1`)

	s.mock.ExpectQuery(query).WithArgs(dcIdMock).WillReturnRows(rows)

	dc, err := s.repo.SelectDebtCollectorById(dcIdMock)

	s.NoError(err)
	s.Equal(expDC, dc)
}

func (s *DebtCollectorRepoTestSuite) TestSelectTugasById_Success() {
	tugasIdMock := "1"
	expTugas := debtCollectorEntity.Tugas{ID: "1", UserId: "2", CollectorId: "1", Status: "ongoing"}

	rows := sqlmock.NewRows([]string{"id", "user_id", "collector_id", "status"}).
		AddRow(expTugas.ID, expTugas.UserId, expTugas.CollectorId, expTugas.Status)
	query := regexp.QuoteMeta(`SELECT id,user_id,collector_id,status FROM claim_tugas WHERE id = $1`)

	s.mock.ExpectQuery(query).WithArgs(tugasIdMock).WillReturnRows(rows)

	tugas, err := s.repo.SelectTugasById(tugasIdMock)

	s.NoError(err)
	s.Equal(expTugas, tugas)
}

func (s *DebtCollectorRepoTestSuite) TestSelectAllLateDebitur_Success() {
	dcCityMock := "Malang"
	pageMock := 1
	sizeMock := 0
	rows := sqlmock.NewRows([]string{"id", "fullname", "address", "unpaid"}).
		AddRow("1", "user satu", "address 1", 1000000)

	query := regexp.QuoteMeta(`SELECT u.id,du.fullname,du.address,SUM(c.jumlah_bayar) AS unpaid
		FROM cicilan c
		INNER JOIN pinjaman p ON c.pinjaman_id = p.id
		INNER JOIN users u ON p.user_id = u.id
		INNER JOIN detail_users du ON u.id = du.user_id
		LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
		WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
		AND du.city ILIKE '%' || $2 || '%'
		AND ct.user_id IS NULL
		GROUP BY u.id, du.fullname, du.address
		LIMIT $3 OFFSET $4`)
	countQuery := regexp.QuoteMeta(`SELECT COUNT(*) FROM (SELECT DISTINCT ON (u.id) u.id
		FROM cicilan c
		INNER JOIN pinjaman p ON c.pinjaman_id = p.id
		INNER JOIN users u ON p.user_id = u.id
		INNER JOIN detail_users du ON u.id = du.user_id
		LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
		WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
		AND du.city ILIKE '%' || $2 || '%'
		AND ct.user_id IS NULL);`)

	offset := (pageMock - 1) * sizeMock
	defaultSize := 10
	s.mock.ExpectQuery(query).WithArgs(AnyTime{}, dcCityMock, defaultSize, offset).WillReturnRows(rows)
	s.mock.ExpectQuery(countQuery).WithArgs(AnyTime{}, dcCityMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	expData := []debtCollectorEntity.LateDebtor{
		{
			ID:           "1",
			FullName:     "user satu",
			Address:      "address 1",
			UnpaidAmount: 1000000,
		},
	}
	expPaging := json.Paging{
		Page:      1,
		TotalData: 3,
	}

	data, paging, err := s.repo.SelectAllLateDebitur(dcCityMock, pageMock, sizeMock)
	s.NoError(err)
	s.Equal(expPaging, paging)
	s.Equal(expData, data)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestInsertLogTugas_Success() {
	mockPayload := debtCollectorDto.NewLogTugasPayload{
		TugasId:     "1",
		Description: "test log",
	}

	query := regexp.QuoteMeta(`INSERT INTO log_tugas(tugas_id,description) VALUES ($1, $2)`)
	s.mock.ExpectExec(query).WithArgs(mockPayload.TugasId, mockPayload.Description).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.InsertLogTugas(mockPayload)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestInsertLogTugas_Fail() {
	mockPayload := debtCollectorDto.NewLogTugasPayload{
		TugasId:     "1",
		Description: "test log",
	}

	query := regexp.QuoteMeta(`INSERT INTO log_tugas(tugas_id,description) VALUES ($1, $2)`)
	expError := fmt.Errorf("internal server error")
	s.mock.ExpectExec(query).WithArgs(mockPayload.TugasId, mockPayload.Description).WillReturnError(expError)

	err := s.repo.InsertLogTugas(mockPayload)
	s.Error(err)
	s.Equal(expError, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestUpdateLogTugasById_Success() {
	mockStoredLog := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "test log",
	}
	mockPayload := debtCollectorDto.UpdateLogTugasPayload{
		Description: "test log update",
	}

	query := regexp.QuoteMeta(`UPDATE log_tugas SET description = $1,updated_at = $2 WHERE id = $3`)
	s.mock.ExpectExec(query).WithArgs(mockPayload.Description, AnyTime{}, mockStoredLog.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.UpdateLogTugasById(mockStoredLog, mockPayload)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestUpdateLogTugasById_Fail() {
	mockStoredLog := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "test log",
	}
	mockPayload := debtCollectorDto.UpdateLogTugasPayload{
		Description: "test log update",
	}

	query := regexp.QuoteMeta(`UPDATE log_tugas SET description = $1,updated_at = $2 WHERE id = $3`)
	expError := fmt.Errorf("internal server error")
	s.mock.ExpectExec(query).WithArgs(mockPayload.Description, AnyTime{}, mockStoredLog.ID).WillReturnError(expError)

	err := s.repo.UpdateLogTugasById(mockStoredLog, mockPayload)
	s.Error(err)
	s.Equal(expError, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestSoftDeleteLogTugasById_Success() {
	mockLogTugasId := "1"

	query := regexp.QuoteMeta(`UPDATE log_tugas SET deleted_at = $1 WHERE id = $2`)
	s.mock.ExpectExec(query).WithArgs(AnyTime{}, mockLogTugasId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.SoftDeleteLogTugasById(mockLogTugasId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestSelectLogTugasById_Success() {
	logTugasIdMock := "1"
	expLogTugas := debtCollectorEntity.LogTugas{
		ID:          "1",
		TugasId:     "1",
		Description: "ini log",
		CreatedAt:   "tanggal dibuat",
		UpdatedAt:   "tanggal diupdate",
	}

	rows := sqlmock.NewRows([]string{"id", "tugas_id", "description", "created_at", "updated_at"}).
		AddRow(expLogTugas.ID, expLogTugas.TugasId, expLogTugas.Description, expLogTugas.CreatedAt, expLogTugas.UpdatedAt)
	query := regexp.QuoteMeta(`SELECT id,tugas_id,description,created_at,updated_at FROM log_tugas WHERE id = $1 AND deleted_at IS NULL`)

	s.mock.ExpectQuery(query).WithArgs(logTugasIdMock).WillReturnRows(rows)

	logTugas, err := s.repo.SelectLogTugasById(logTugasIdMock)

	s.NoError(err)
	s.Equal(expLogTugas, logTugas)
}

func (s *DebtCollectorRepoTestSuite) TestSelectAllLogByTugasId_Success() {
	tugasIdMock := "1"
	pageMock := 1
	sizeMock := 0
	rows := sqlmock.NewRows([]string{"id", "description", "created_at", "updated_at"}).
		AddRow("1", "log satu", "tanggal dibuat", "tanggal diupdate")

	query := regexp.QuoteMeta(`SELECT id,description,created_at,updated_at
	FROM log_tugas
	WHERE tugas_id = $1 AND deleted_at IS NULL
	ORDER BY created_at ASC LIMIT $2 OFFSET $3`)
	countQuery := regexp.QuoteMeta(`SELECT COUNT(*)
	FROM log_tugas
	WHERE tugas_id = $1 AND deleted_at IS NULL`)

	offset := (pageMock - 1) * sizeMock
	defaultSize := 10
	s.mock.ExpectQuery(query).WithArgs(tugasIdMock, defaultSize, offset).WillReturnRows(rows)
	s.mock.ExpectQuery(countQuery).WithArgs(tugasIdMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	expData := []debtCollectorEntity.LogTugas{
		{
			ID:          "1",
			Description: "log satu",
			CreatedAt:   "tanggal dibuat",
			UpdatedAt:   "tanggal diupdate",
		},
	}
	expPaging := json.Paging{
		Page:      1,
		TotalData: 3,
	}

	data, paging, err := s.repo.SelectAllLogByTugasId(tugasIdMock, pageMock, sizeMock)
	s.NoError(err)
	s.Equal(expPaging, paging)
	s.Equal(expData, data)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestSelectLateDebiturById_Success() {
	userIdMock := "1"
	dcCityMock := "malang"
	expId := "1"

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expId)
	query := regexp.QuoteMeta(`SELECT DISTINCT ON (u.id) u.id
	FROM cicilan c
	INNER JOIN pinjaman p ON c.pinjaman_id = p.id
	INNER JOIN users u ON p.user_id = u.id
	INNER JOIN detail_users du ON u.id = du.user_id
	LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
	WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
	AND du.city ILIKE '%' || $2 || '%'
	AND ct.user_id IS NULL AND u.id = $3`)

	s.mock.ExpectQuery(query).WithArgs(AnyTime{}, dcCityMock, userIdMock).WillReturnRows(rows)

	lateDebtor, err := s.repo.SelectLateDebiturById(userIdMock, dcCityMock)

	s.NoError(err)
	s.Equal(expId, lateDebtor)
}

func (s *DebtCollectorRepoTestSuite) TestCreateClaimTugas_Success() {
	dcIdMock := "1"
	userIdMock := "2"

	query := regexp.QuoteMeta(`INSERT INTO claim_tugas(user_id,collector_id) VALUES($1,$2);`)
	s.mock.ExpectExec(query).WithArgs(userIdMock, dcIdMock).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.CreateClaimTugas(dcIdMock, userIdMock)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestSelectAllTugas_Success() {
	dcIdMock := "1"
	statusMock := "ongoing"
	pageMock := 1
	sizeMock := 0
	rows := sqlmock.NewRows([]string{"id", "user_id", "status"}).
		AddRow("1", "2", "ongoing")

	query := regexp.QuoteMeta(`SELECT id,user_id,status
	FROM claim_tugas
	WHERE collector_id = $1
	AND ($2 = '' OR status = $2::claim_status)
	ORDER BY created_at ASC LIMIT $3 OFFSET $4`)
	countQuery := regexp.QuoteMeta(`SELECT COUNT(*)
	FROM claim_tugas
	WHERE collector_id = $1
	AND ($2 = '' OR status = $2::claim_status);`)

	offset := (pageMock - 1) * sizeMock
	defaultSize := 10
	s.mock.ExpectQuery(query).WithArgs(dcIdMock, statusMock, defaultSize, offset).WillReturnRows(rows)
	s.mock.ExpectQuery(countQuery).WithArgs(dcIdMock, statusMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	expData := []debtCollectorEntity.Tugas{
		{
			ID:     "1",
			UserId: "2",
			Status: "ongoing",
		},
	}
	expPaging := json.Paging{
		Page:      1,
		TotalData: 3,
	}

	data, paging, err := s.repo.SelectAllTugas(dcIdMock, statusMock, pageMock, sizeMock)
	s.NoError(err)
	s.Equal(expPaging, paging)
	s.Equal(expData, data)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DebtCollectorRepoTestSuite) TestCountOngoingTugas_Success() {
	dcIdMock := "1"
	expTotalTugas := 2

	query := regexp.QuoteMeta(`SELECT COUNT(*) FROM claim_tugas WHERE collector_id = $1 AND status = 'ongoing';`)

	s.mock.ExpectQuery(query).WithArgs(dcIdMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expTotalTugas))

	totalTugas, err := s.repo.CountOngoingTugas(dcIdMock)

	s.NoError(err)
	s.Equal(expTotalTugas, totalTugas)
}
