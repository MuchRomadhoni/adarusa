package debtCollectorRepository

import (
	"database/sql"
	"fmt"
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	"fp_pinjaman_online/src/debtCollector"
	"strings"
	"time"
)

type debtCollectorRepository struct {
	db *sql.DB
}

func NewDebtCollectorRepository(db *sql.DB) debtCollector.DebtCollectorRepository {
	return &debtCollectorRepository{db}
}

// TODO
// this func should be in users repository
func (repo *debtCollectorRepository) SelectDebtCollectorById(id string) (debtCollectorEntity.DebtCollector, error) {
	var debtCollector debtCollectorEntity.DebtCollector
	query := "SELECT u.id,du.fullname,du.city FROM users u JOIN detail_users du ON du.user_id = u.id WHERE u.id = $1"
	err := repo.db.QueryRow(query, id).Scan(&debtCollector.ID, &debtCollector.FullName, &debtCollector.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return debtCollectorEntity.DebtCollector{}, fmt.Errorf("users with id: %v not found", id)
		}
		return debtCollectorEntity.DebtCollector{}, err
	}
	return debtCollector, nil
}

func (repo *debtCollectorRepository) SelectTugasById(tugasId string) (debtCollectorEntity.Tugas, error) {
	var tugas debtCollectorEntity.Tugas
	query := "SELECT id,user_id,collector_id,status FROM claim_tugas WHERE id = $1"
	err := repo.db.QueryRow(query, tugasId).Scan(&tugas.ID, &tugas.UserId, &tugas.CollectorId, &tugas.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return debtCollectorEntity.Tugas{}, fmt.Errorf("tugas with id: %v not found", tugasId)
		}
		return debtCollectorEntity.Tugas{}, err
	}
	return tugas, nil
}

func (repo *debtCollectorRepository) SelectAllLateDebitur(dcCity string, page, size int) ([]debtCollectorEntity.LateDebtor, json.Paging, error) {
	// late if cicilan = unpaid more than 2 months
	var rows *sql.Rows
	var err error
	var offset int
	var newPaging json.Paging

	if page == 0 || size == 0 {
		page = 1
		size = 10
	}

	query := `SELECT u.id,du.fullname,du.address,SUM(c.jumlah_bayar) AS unpaid
	FROM cicilan c
	INNER JOIN pinjaman p ON c.pinjaman_id = p.id
	INNER JOIN users u ON p.user_id = u.id
	INNER JOIN detail_users du ON u.id = du.user_id
	LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
	WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
	AND du.city ILIKE '%' || $2 || '%'
	AND ct.user_id IS NULL
	GROUP BY u.id, du.fullname, du.address`

	countQuery := `SELECT COUNT(*) FROM (SELECT DISTINCT ON (u.id) u.id
	FROM cicilan c
	INNER JOIN pinjaman p ON c.pinjaman_id = p.id
	INNER JOIN users u ON p.user_id = u.id
	INNER JOIN detail_users du ON u.id = du.user_id
	LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
	WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
	AND du.city ILIKE '%' || $2 || '%'
	AND ct.user_id IS NULL);`

	offset = (page - 1) * size
	query += " LIMIT $3 OFFSET $4"
	lateMonthLimit := time.Now().AddDate(0, -2, 0)
	rows, err = repo.db.Query(query, lateMonthLimit, dcCity, size, offset)
	if err != nil {
		return nil, json.Paging{}, err
	}
	defer rows.Close()

	err = repo.db.QueryRow(countQuery, lateMonthLimit, dcCity).Scan(&newPaging.TotalData)
	if err != nil {
		return nil, json.Paging{}, err
	}

	listLateDebtors := scanLateDebitur(rows)
	newPaging.Page = page
	return listLateDebtors, newPaging, nil
}

func (repo *debtCollectorRepository) InsertLogTugas(newLogPayload debtCollectorDto.NewLogTugasPayload) error {
	query := "INSERT INTO log_tugas(tugas_id,description) VALUES ($1, $2)"
	_, err := repo.db.Exec(query, newLogPayload.TugasId, newLogPayload.Description)
	if err != nil {
		return err
	}
	return nil
}

func (repo *debtCollectorRepository) UpdateLogTugasById(storedLog debtCollectorEntity.LogTugas, updateLogPayload debtCollectorDto.UpdateLogTugasPayload) error {
	if strings.TrimSpace(updateLogPayload.Description) != "" {
		storedLog.Description = updateLogPayload.Description
	}
	query := "UPDATE log_tugas SET description = $1,updated_at = $2 WHERE id = $3"
	_, err := repo.db.Exec(query, storedLog.Description, time.Now(), storedLog.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *debtCollectorRepository) SoftDeleteLogTugasById(logTugasId string) error {
	query := "UPDATE log_tugas SET deleted_at = $1 WHERE id = $2"
	_, err := repo.db.Exec(query, time.Now(), logTugasId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *debtCollectorRepository) SelectLogTugasById(logTugasId string) (debtCollectorEntity.LogTugas, error) {
	var logTugas debtCollectorEntity.LogTugas
	query := "SELECT id,tugas_id,description,created_at,updated_at FROM log_tugas WHERE id = $1 AND deleted_at IS NULL"
	err := repo.db.QueryRow(query, logTugasId).Scan(&logTugas.ID, &logTugas.TugasId, &logTugas.Description, &logTugas.CreatedAt, &logTugas.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return debtCollectorEntity.LogTugas{}, fmt.Errorf("log tugas with id: %v not found", logTugasId)
		}
		return debtCollectorEntity.LogTugas{}, err
	}
	return logTugas, nil
}

func (repo *debtCollectorRepository) SelectAllLogByTugasId(tugasId string, page, size int) ([]debtCollectorEntity.LogTugas, json.Paging, error) {
	var rows *sql.Rows
	var err error
	var offset int
	var newPaging json.Paging

	if page == 0 || size == 0 {
		page = 1
		size = 10
	}

	query := `SELECT id,description,created_at,updated_at
	FROM log_tugas
	WHERE tugas_id = $1 AND deleted_at IS NULL
	ORDER BY created_at ASC`

	countQuery := `SELECT COUNT(*)
	FROM log_tugas
	WHERE tugas_id = $1 AND deleted_at IS NULL`

	offset = (page - 1) * size
	query += " LIMIT $2 OFFSET $3"
	rows, err = repo.db.Query(query, tugasId, size, offset)
	if err != nil {
		return nil, json.Paging{}, err
	}
	defer rows.Close()

	err = repo.db.QueryRow(countQuery, tugasId).Scan(&newPaging.TotalData)
	if err != nil {
		return nil, json.Paging{}, err
	}

	logList := scanTugasLogs(rows)
	newPaging.Page = page
	return logList, newPaging, nil
}

func (repo *debtCollectorRepository) SelectLateDebiturById(userId, dcCity string) (string, error) {
	var id string
	query := `SELECT DISTINCT ON (u.id) u.id
	FROM cicilan c
	INNER JOIN pinjaman p ON c.pinjaman_id = p.id
	INNER JOIN users u ON p.user_id = u.id
	INNER JOIN detail_users du ON u.id = du.user_id
	LEFT JOIN claim_tugas ct ON u.id = ct.user_id AND ct.status = 'ongoing'
	WHERE c.tanggal_jatuh_tempo < $1 AND c.status = 'unpaid'
	AND du.city ILIKE '%' || $2 || '%'
	AND ct.user_id IS NULL AND u.id = $3`

	lateMonthLimit := time.Now().AddDate(0, -2, 0)
	err := repo.db.QueryRow(query, lateMonthLimit, dcCity, userId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("id invalid or not found")
		}
		return "", err
	}
	return id, nil
}

func (repo *debtCollectorRepository) CreateClaimTugas(dcId, userId string) error {
	query := "INSERT INTO claim_tugas(user_id,collector_id) VALUES($1,$2);"
	_, err := repo.db.Exec(query, userId, dcId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *debtCollectorRepository) SelectAllTugas(dcId, status string, page, size int) ([]debtCollectorEntity.Tugas, json.Paging, error) {
	var rows *sql.Rows
	var err error
	var offset int
	var newPaging json.Paging

	if page == 0 || size == 0 {
		page = 1
		size = 10
	}

	query := `SELECT id,user_id,status
	FROM claim_tugas
	WHERE collector_id = $1
	AND ($2 = '' OR status = $2::claim_status)
	ORDER BY created_at ASC`

	countQuery := `SELECT COUNT(*)
	FROM claim_tugas
	WHERE collector_id = $1
	AND ($2 = '' OR status = $2::claim_status);`

	offset = (page - 1) * size
	query += " LIMIT $3 OFFSET $4"
	rows, err = repo.db.Query(query, dcId, status, size, offset)
	if err != nil {
		return nil, json.Paging{}, err
	}
	defer rows.Close()

	err = repo.db.QueryRow(countQuery, dcId, status).Scan(&newPaging.TotalData)
	if err != nil {
		return nil, json.Paging{}, err
	}

	tasks := scanTugas(rows)
	newPaging.Page = page
	return tasks, newPaging, nil
}

func (repo *debtCollectorRepository) CountOngoingTugas(dcId string) (int, error) {
	var totalTugas int
	query := "SELECT COUNT(*) FROM claim_tugas WHERE collector_id = $1 AND status = 'ongoing';"
	err := repo.db.QueryRow(query, dcId).Scan(&totalTugas)
	if err != nil {
		return 0, err
	}
	return totalTugas, nil
}

func (repo *debtCollectorRepository) SelectBalanceByUserId(userId string) (float64, error) {
	var balance float64
	query := "SELECT amount FROM balance WHERE user_id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}

func (repo *debtCollectorRepository) CreateWithdrawRequest(userId string, amount float64) error {
	query := "INSERT INTO withdrawal(user_id,amount) VALUES($1, $2);"
	_, err := repo.db.Exec(query, userId, amount)
	if err != nil {
		return err
	}
	return nil
}

func (repo *debtCollectorRepository) SelectDebtorFromTugas(dcId, userId string) (string, error) {
	var debiturId string
	query := "SELECT user_id FROM claim_tugas WHERE user_id = $1 AND collector_id = $2 AND status = 'ongoing'"
	err := repo.db.QueryRow(query, userId, dcId).Scan(&debiturId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("debtor not found")
		}
		return "", err
	}
	return debiturId, nil
}

func scanTugasLogs(rows *sql.Rows) []debtCollectorEntity.LogTugas {
	var logs []debtCollectorEntity.LogTugas
	var err error
	for rows.Next() {
		log := debtCollectorEntity.LogTugas{}
		err = rows.Scan(&log.ID, &log.Description, &log.CreatedAt, &log.UpdatedAt)
		if err != nil {
			panic(err)
		}
		logs = append(logs, log)
	}

	return logs
}

func scanLateDebitur(rows *sql.Rows) []debtCollectorEntity.LateDebtor {
	var debtors []debtCollectorEntity.LateDebtor
	var err error
	for rows.Next() {
		debtor := debtCollectorEntity.LateDebtor{}
		err = rows.Scan(&debtor.ID, &debtor.FullName, &debtor.Address, &debtor.UnpaidAmount)
		if err != nil {
			panic(err)
		}
		debtors = append(debtors, debtor)
	}

	return debtors
}

func scanTugas(rows *sql.Rows) []debtCollectorEntity.Tugas {
	var tasks []debtCollectorEntity.Tugas
	var err error
	for rows.Next() {
		task := debtCollectorEntity.Tugas{}
		err = rows.Scan(&task.ID, &task.UserId, &task.Status)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}
