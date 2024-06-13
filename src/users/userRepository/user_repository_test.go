package userRepository

import (
	"database/sql"
	"fp_pinjaman_online/model/dcFormDto"
	"fp_pinjaman_online/model/debiturFormDto"
	"fp_pinjaman_online/model/entity/usersEntity"
	"fp_pinjaman_online/model/userDto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Setup a mock database connection
func setupMockDB() (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock, func() {
		db.Close()
	}
}

// Test CreateUser
func TestCreateUser_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)
	req := userDto.CreateRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(req.Email, req.Password, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateUser(req, 1)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

// Test Login
func TestLogin_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)
	req := userDto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mock.ExpectQuery(`SELECT password FROM users WHERE email=\$1`).
		WithArgs(req.Email).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("hashedpassword"))

	hashedPassword, err := repo.Login(req)
	assert.NoError(t, err)
	assert.Equal(t, "hashedpassword", hashedPassword)
	mock.ExpectationsWereMet()
}

// Test UserExists
func TestUserExists_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email=\\$1\\)").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.UserExists("test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
	mock.ExpectationsWereMet()
}

// Test GetUserByEmail
func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT u.id, u.email, u.password, r.roles_name, u.status FROM users u JOIN roles r ON u.role_id = r.id WHERE u.email=\$1`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "roles_name", "status"}).
			AddRow("1", "test@example.com", "hashedpassword", "user", "status"))

	user, err := repo.GetUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "1", user.Id)  // if Id is a string
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "user", user.Roles)
	assert.Equal(t, "status", user.Status)
	mock.ExpectationsWereMet()
}

func TestGetUserByEmail_Failed(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT u.id, u.email, u.password, r.roles_name, u.status FROM users u JOIN roles r ON u.role_id = r.id WHERE u.email=\$1`).
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByEmail("notfound@example.com")
	assert.Error(t, err)
	assert.Equal(t, "user with email notfound@example.com not found", err.Error())
	assert.Equal(t, "", user.Id)
	assert.Equal(t, "", user.Email)
	assert.Equal(t, "", user.Password)
	assert.Equal(t, "", user.Roles)
	assert.Equal(t, "", user.Status)
	mock.ExpectationsWereMet()
}

// Test UpdateBankAccount
func TestUpdateBankAccount(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	mock.ExpectExec("INSERT INTO rekening").
		WithArgs(1, "1234567890", "Test Bank").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateBankAccount(1, "1234567890", "Test Bank")
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

// Test IsBankAccExist
func TestIsBankAccExist(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM rekening WHERE user_id=\\$1 AND account_number=\\$2\\)").
		WithArgs(1, "1234567890").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.IsBankAccExist(1, "1234567890")
	assert.NoError(t, err)
	assert.True(t, exists)
	mock.ExpectationsWereMet()
}

func TestCreateDetailDebitur_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	req := debiturFormDto.Debitur{
		DetailUser: debiturFormDto.DetailDebitur{
			UserID:     1,
			Nik:        "1234567890",
			Fullname:   "John Doe",
			PhoneNumber: "123456789",
			Address:    "123 Main St",
			City:       "Metropolis",
		},
		UserJobs: debiturFormDto.UserJobs{
			UserID:        1,
			JobName:       "Engineer",
			Salary:        5000,
			OfficeName:    "Tech Co.",
			OfficeContact: "office_contact",
			OfficeAddress: "office_address",
		},
		EmergencyContact: debiturFormDto.EmergencyContact{
			UserID:      1,
			Name:        "Jane Doe",
			PhoneNumber: "0987654321",
		},
	}

	mock.ExpectBegin()

	// Check if NIK is unique
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM detail_users WHERE nik=\$1 AND user_id != \$2\)`).
		WithArgs(req.DetailUser.Nik, req.DetailUser.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Check if detail user already exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM detail_users WHERE user_id=\$1\)`).
		WithArgs(req.DetailUser.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Insert detail user
	mock.ExpectExec(`INSERT INTO detail_users \(user_id, nik, fullname, phone_number, address, city\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(req.DetailUser.UserID, req.DetailUser.Nik, req.DetailUser.Fullname, req.DetailUser.PhoneNumber, req.DetailUser.Address, req.DetailUser.City).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Upsert user jobs
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users_job_detail WHERE user_id=\$1\)`).
		WithArgs(req.UserJobs.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(`INSERT INTO users_job_detail \(user_id, job_name, gaji, office_name, office_contact, address\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(req.UserJobs.UserID, req.UserJobs.JobName, req.UserJobs.Salary, req.UserJobs.OfficeName, req.UserJobs.OfficeContact, req.UserJobs.OfficeAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Upsert emergency contact
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM kontak_darurat WHERE user_id=\$1\)`).
		WithArgs(req.EmergencyContact.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(`INSERT INTO kontak_darurat \(user_id, name, phone_number\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(req.EmergencyContact.UserID, req.EmergencyContact.Name, req.EmergencyContact.PhoneNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.CreateDetailDebitur(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateDetailDc_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	req := dcFormDto.DetailDC{
		UserID:      1,
		Nik:         "1234567890",
		Fullname:    "John Doe",
		PhoneNumber: "123456789",
		Address:     "123 Main St",
		City:        "Metropolis",
	}

	mock.ExpectBegin()

	// Check if NIK is unique
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM detail_users WHERE nik=\$1 AND user_id != \$2\)`).
		WithArgs(req.Nik, req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Check if detail user already exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM detail_users WHERE user_id=\$1\)`).
		WithArgs(req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Insert detail user
	mock.ExpectExec(`INSERT INTO detail_users \(user_id, nik, fullname, phone_number, address, city\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(req.UserID, req.Nik, req.Fullname, req.PhoneNumber, req.Address, req.City).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.CreateDetailDc(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePhotoPaths_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := 1
	fotoKTP := "new_foto_ktp_path"
	fotoSelfie := "new_foto_selfie_path"

	mock.ExpectExec(`UPDATE detail_users SET foto_ktp=\$1, foto_selfie=\$2 WHERE user_id=\$3`).
		WithArgs(fotoKTP, fotoSelfie, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePhotoPaths(userID, fotoKTP, fotoSelfie)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullname_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := 1
	expectedFullname := "John Doe"

	mock.ExpectQuery(`SELECT fullname FROM detail_users WHERE user_id=\$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"fullname"}).AddRow(expectedFullname))

	fullname, err := repo.GetFullname(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedFullname, fullname)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDataByRole_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	role := "some_role"
	status := "some_status"
	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT u.id, du.nik, du.fullname, du.phone_number, du.address, du.city, du.foto_ktp, du.foto_selfie, du.limit_id FROM users u JOIN detail_users du ON u.id = du.user_id JOIN roles r ON u.role_id = r.id JOIN limit_pinjaman lp ON du.limit_id = lp.id WHERE r.roles_name = \$1 AND \(\$2 = '' OR u.status = \$2::user_status\) LIMIT \$3 OFFSET \$4`).
		WithArgs(role, status, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nik", "fullname", "phone_number", "address", "city", "foto_ktp", "foto_selfie", "limit_id"}).
			AddRow(1, "1234567890", "John Doe", "123456789", "123 Main St", "Metropolis", "foto_ktp_path", "foto_selfie_path", 1))

	mock.ExpectQuery(`SELECT count\(\*\) FROM users u JOIN detail_users du ON u.id = du.user_id JOIN roles r ON u.role_id = r.id JOIN limit_pinjaman lp ON du.limit_id = lp.id WHERE r.roles_name = \$1 AND \(\$2 = '' OR u.status = \$2::user_status\)`).
		WithArgs(role, status).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	data, total, err := repo.GetDataByRole(role, status, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, 1, total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRolesById_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := "1"
	role := "some_role"

	mock.ExpectQuery(`SELECT role_id FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(role))

	resultRole, err := repo.GetRolesById(userID)
	assert.NoError(t, err)
	assert.Equal(t, role, resultRole)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserDetailByUserId_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := "1"
	expectedDetail := usersEntity.DetailUser{
		NIK:         "1234567890",
		FullName:    "John Doe",
		PhoneNumber: "123456789",
		Address:     "123 Main St",
		City:        "Metropolis",
		FotoKtp:     "foto_ktp_path",
		FotoSelfie:  "foto_selfie_path",
	}

	mock.ExpectQuery(`SELECT nik,fullname,phone_number,address,city,foto_ktp,foto_selfie FROM detail_users WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"nik", "fullname", "phone_number", "address", "city", "foto_ktp", "foto_selfie"}).
			AddRow(expectedDetail.NIK, expectedDetail.FullName, expectedDetail.PhoneNumber, expectedDetail.Address, expectedDetail.City, expectedDetail.FotoKtp, expectedDetail.FotoSelfie))

	resultDetail, err := repo.GetUserDetailByUserId(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedDetail, resultDetail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserJobDetailByUserId_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := "1"
	expectedJob := usersEntity.UserJobDetail{
		JobName:       "Engineer",
		Salary:        5000,
		OfficeName:    "Tech Co.",
		OfficeContact: "office_contact",
		OfficeAddress: "office_address",
	}

	mock.ExpectQuery(`SELECT job_name,gaji,office_name,office_contact,address FROM users_job_detail WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"job_name", "gaji", "office_name", "office_contact", "address"}).
			AddRow(expectedJob.JobName, expectedJob.Salary, expectedJob.OfficeName, expectedJob.OfficeContact, expectedJob.OfficeAddress))

	resultJob, err := repo.GetUserJobDetailByUserId(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedJob, resultJob)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEmergencyContactByUserId_Success(t *testing.T) {
	db, mock, teardown := setupMockDB()
	defer teardown()

	repo := NewUserRepository(db)

	userID := "1"
	expectedContact := usersEntity.EmergencyContact{
		ContactName: "Jane Doe",
		PhoneNumber: "0987654321",
	}

	mock.ExpectQuery(`SELECT name,phone_number FROM kontak_darurat WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "phone_number"}).
			AddRow(expectedContact.ContactName, expectedContact.PhoneNumber))

	resultContact, err := repo.GetEmergencyContactByUserId(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedContact, resultContact)
	assert.NoError(t, mock.ExpectationsWereMet())
}