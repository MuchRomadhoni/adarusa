package userRepository

import (
	"database/sql"
	"fmt"
	"fp_pinjaman_online/model/dcFormDto"
	"fp_pinjaman_online/model/debiturFormDto"
	"fp_pinjaman_online/model/entity/usersEntity"
	"fp_pinjaman_online/model/userDto"
	"fp_pinjaman_online/src/users"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) users.UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) CreateUser(req userDto.CreateRequest, roleId int) error {
	query := "INSERT INTO users(email, password, role_id) VALUES ($1,$2, $3)"
	_, err := repo.db.Exec(query, req.Email, req.Password, roleId)
	return err
}

func (repo *userRepository) Login(req userDto.LoginRequest) (string, error) {
	var hashedPassword string

	query := "SELECT password FROM users WHERE email=$1"
	err := repo.db.QueryRow(query, req.Email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", nil
	}

	return hashedPassword, nil
}

func (repo *userRepository) UserExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := repo.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (repo *userRepository) GetUserByEmail(email string) (userDto.User, error) {
	var user userDto.User
	query := `
        SELECT u.id, u.email, u.password, r.roles_name, u.status
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.email=$1`
	err := repo.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Roles, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with email %s not found", email)
		}
		return user, err
	}
	return user, nil
}

func (repo *userRepository) CreateDetailDebitur(req debiturFormDto.Debitur) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	// Check if nik is unique
	var exists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM detail_users WHERE nik=$1 AND user_id != $2)`, req.DetailUser.Nik, req.DetailUser.UserID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		tx.Rollback()
		return fmt.Errorf("nik already exists")
	}

	// Check if detailuser already exists
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM detail_users WHERE user_id=$1)`, req.DetailUser.UserID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		_, err = tx.Exec(`UPDATE detail_users SET nik=$1, fullname=$2, phone_number=$3, address=$4, city=$5 WHERE user_id=$6`, req.DetailUser.Nik, req.DetailUser.Fullname, req.DetailUser.PhoneNumber, req.DetailUser.Address, req.DetailUser.City, req.DetailUser.UserID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec(`INSERT INTO detail_users (user_id, nik, fullname, phone_number, address, city) VALUES ($1, $2, $3, $4, $5, $6)`, req.DetailUser.UserID, req.DetailUser.Nik, req.DetailUser.Fullname, req.DetailUser.PhoneNumber, req.DetailUser.Address, req.DetailUser.City)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = upsertUserJobDetail(tx, req.UserJobs)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = upsertEmergencyContact(tx, req.EmergencyContact)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *userRepository) CreateDetailDc(req dcFormDto.DetailDC) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	// Check if nik is unique
	var exists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM detail_users WHERE nik=$1 AND user_id != $2)`, req.Nik, req.UserID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		tx.Rollback()
		return fmt.Errorf("nik already exists")
	}

	// Check if detailuser already exists
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM detail_users WHERE user_id=$1)`, req.UserID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		_, err = tx.Exec(`UPDATE detail_users SET nik=$1, fullname=$2, phone_number=$3, address=$4, city=$5 WHERE user_id=$6`, req.Nik, req.Fullname, req.PhoneNumber, req.Address, req.City, req.UserID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec(`INSERT INTO detail_users (user_id, nik, fullname, phone_number, address, city) VALUES ($1, $2, $3, $4, $5, $6)`, req.UserID, req.Nik, req.Fullname, req.PhoneNumber, req.Address, req.City)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func upsertUserJobDetail(tx *sql.Tx, jobDetail debiturFormDto.UserJobs) error {
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users_job_detail WHERE user_id=$1)`, jobDetail.UserID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = tx.Exec(`
        UPDATE users_job_detail SET job_name=$1, gaji=$2, office_name=$3, office_contact=$4, address=$5 WHERE user_id=$6
        `, jobDetail.JobName, jobDetail.Salary, jobDetail.OfficeName, jobDetail.OfficeContact, jobDetail.OfficeAddress, jobDetail.UserID)
	} else {
		_, err = tx.Exec(`
        INSERT INTO users_job_detail (user_id, job_name, gaji, office_name, office_contact, address)
        VALUES ($1, $2, $3, $4, $5, $6)
        `, jobDetail.UserID, jobDetail.JobName, jobDetail.Salary, jobDetail.OfficeName, jobDetail.OfficeContact, jobDetail.OfficeAddress)
	}
	return err
}

func upsertEmergencyContact(tx *sql.Tx, contact debiturFormDto.EmergencyContact) error {
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM kontak_darurat WHERE user_id=$1)`, contact.UserID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = tx.Exec(`
        UPDATE kontak_darurat SET name=$1, phone_number=$2 WHERE user_id=$3
        `, contact.Name, contact.PhoneNumber, contact.UserID)
	} else {
		_, err = tx.Exec(`
        INSERT INTO kontak_darurat (user_id, name, phone_number)
        VALUES ($1, $2, $3)
        `, contact.UserID, contact.Name, contact.PhoneNumber)
	}
	return err
}

func (repo *userRepository) UpdatePhotoPaths(userId int, fotoKTP, fotoSelfie string) error {
	_, err := repo.db.Exec(`
    UPDATE detail_users SET foto_ktp=$1, foto_selfie=$2 WHERE user_id=$3
    `, fotoKTP, fotoSelfie, userId)
	return err
}

func (repo *userRepository) GetFullname(userId int) (string, error) {
	var fullname string
	err := repo.db.QueryRow("SELECT fullname FROM detail_users WHERE user_id=$1", userId).Scan(&fullname)
	fmt.Println("fullname:", fullname)
	if err != nil {
		return "", nil
	}

	return fullname, nil
}

func (repo *userRepository) GetDataByRole(role, status string, limit, offset int) ([]debiturFormDto.DetailDebitur, int, error) {
	var debitur []debiturFormDto.DetailDebitur
	var totalData int

	// Base queries with conditional status filter
	query := `
        SELECT u.id, du.nik, du.fullname, du.phone_number, du.address, du.city, du.foto_ktp, du.foto_selfie, du.limit_id
        FROM users u
        JOIN detail_users du ON u.id = du.user_id
        JOIN roles r ON u.role_id = r.id
        JOIN limit_pinjaman lp ON du.limit_id = lp.id
        WHERE r.roles_name = $1
        AND ($2 = '' OR u.status = $2::user_status)
        LIMIT $3 OFFSET $4`

	countQuery := `
        SELECT count(*)
        FROM users u
        JOIN detail_users du ON u.id = du.user_id
        JOIN roles r ON u.role_id = r.id
        JOIN limit_pinjaman lp ON du.limit_id = lp.id
        WHERE r.roles_name = $1
        AND ($2 = '' OR u.status = $2::user_status)`

	// Arguments for the queries
	args := []interface{}{role, status, limit, offset}

	// Execute the main query
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process the result set
	for rows.Next() {
		var dbt debiturFormDto.DetailDebitur
		err := rows.Scan(&dbt.UserID, &dbt.Nik, &dbt.Fullname, &dbt.PhoneNumber, &dbt.Address, &dbt.City, &dbt.FotoKtp, &dbt.FotoSelfie, &dbt.LimitID)
		if err != nil {
			return nil, 0, err
		}
		debitur = append(debitur, dbt)
	}

	// Execute the count query
	countArgs := args[:2] // Use only role and status for the count query
	err = repo.db.QueryRow(countQuery, countArgs...).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	return debitur, totalData, nil
}

func (repo *userRepository) GetRolesById(userId string) (string, error) {
	var role string
	query := "SELECT role_id FROM users WHERE id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("data not found")
		}
		return "", err
	}
	return role, nil
}

func (repo *userRepository) GetUserDetailByUserId(userId string) (usersEntity.DetailUser, error) {
	var detail usersEntity.DetailUser
	query := "SELECT nik,fullname,phone_number,address,city,foto_ktp,foto_selfie FROM detail_users WHERE user_id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&detail.NIK, &detail.FullName, &detail.PhoneNumber, &detail.Address, &detail.City, &detail.FotoKtp, &detail.FotoSelfie)
	if err != nil {
		if err == sql.ErrNoRows {
			return usersEntity.DetailUser{}, fmt.Errorf("data not found")
		}
		return usersEntity.DetailUser{}, err
	}
	return detail, nil
}

func (repo *userRepository) GetUserJobDetailByUserId(userId string) (usersEntity.UserJobDetail, error) {
	var job usersEntity.UserJobDetail
	query := "SELECT job_name,gaji,office_name,office_contact,address FROM users_job_detail WHERE user_id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&job.JobName, &job.Salary, &job.OfficeName, &job.OfficeContact, &job.OfficeAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return usersEntity.UserJobDetail{}, fmt.Errorf("data not found")
		}
		return usersEntity.UserJobDetail{}, err
	}
	return job, nil
}

func (repo *userRepository) GetEmergencyContactByUserId(userId string) (usersEntity.EmergencyContact, error) {
	var mrgcContact usersEntity.EmergencyContact
	query := "SELECT name,phone_number FROM kontak_darurat WHERE user_id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&mrgcContact.ContactName, &mrgcContact.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return usersEntity.EmergencyContact{}, fmt.Errorf("data not found")
		}
		return usersEntity.EmergencyContact{}, err
	}
	return mrgcContact, nil
}

func (repo *userRepository) UpdateBankAccount(userId int, accountNumber, bankName string) error {
	_, err := repo.db.Exec(`INSERT INTO rekening (user_id, account_number, bank_name) VALUES ($1, $2, $3)`, userId, accountNumber, bankName)
	return err
}

func (repo *userRepository) IsBankAccExist(userId int, accountNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rekening WHERE user_id=$1 AND account_number=$2)`
	err := repo.db.QueryRow(query, userId, accountNumber).Scan(&exists)
	return exists, err
}
