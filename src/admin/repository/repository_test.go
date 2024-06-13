package adminRepository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveUserStatusById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAdminRepository(db) // adjust this to your actual repository constructor

	columns := []string{
		"user_id", "email", "status", "account_number", "bank_name",
		"emergency_contact_name", "emergency_contact_phone",
		"job_name", "gaji", "office_name",
		"NIK", "fullname", "personal_phone_number", "personal_address",
		"city", "foto_ktp", "foto_selfie",
		"created_at", "updated_at", "verified_at", "deleted_at",
	}

	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(
		1, "example@example.com", "active", "123456789", "BCA",
		"John Doe", "081234567890",
		"Software Engineer", 10000000, "XYZ Corp",
		"123456789012", "Jane Doe", "081234567891", "123 Fake St",
		"Jakarta", []byte("ktp_image"), []byte("selfie_image"),
		time.Now(), time.Now(), time.Now(), nil,
	))

	userCompleteInfo, err := repo.RetrieveUserStatusById(1)

	assert.NoError(t, err)
	assert.NotNil(t, userCompleteInfo)
	assert.Equal(t, 1, userCompleteInfo.UserID)
	assert.Equal(t, "example@example.com", userCompleteInfo.Email)
	assert.Equal(t, "active", userCompleteInfo.Status)

	// Make sure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAdminRepository(db) // Sesuaikan dengan nama constructor yang benar

	query := `
	UPDATE users
		SET status = \$2,
    		verified_at = NOW(),
    		updated_at = NOW()
		WHERE id = \$1;
    `

	mock.ExpectExec(query).WithArgs(1, "verified").WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUserStatus(1, "verified")
	assert.NoError(t, err)

	// Memastikan semua ekspektasi terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
