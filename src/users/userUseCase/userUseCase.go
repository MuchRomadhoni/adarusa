package userUseCase

import (
	"errors"
	"fp_pinjaman_online/model/dcFormDto"
	"fp_pinjaman_online/model/debiturFormDto"
	"fp_pinjaman_online/model/entity/usersEntity"
	"fp_pinjaman_online/model/userDto"
	"fp_pinjaman_online/pkg/middleware"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/src/users"
)

type userUC struct {
	userRepo users.UserRepository
}

func NewUserUseCase(userRepo users.UserRepository) users.UserUseCase {
	return &userUC{userRepo}
}

func (useCase *userUC) CreateUser(req userDto.CreateRequest, roleId int) error {
	exists, err := useCase.userRepo.UserExists(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exist")
	}

	req.Password, err = validation.HashedPassword(req.Password)
	if err != nil {
		return err
	}

	return useCase.userRepo.CreateUser(req, roleId)
}

func (useCase *userUC) Login(req userDto.LoginRequest) (string, error) {
	hashedPassword, err := useCase.userRepo.Login(req)
	if err != nil {
		return "", nil
	}
	if hashedPassword == "" {
		return "", errors.New("invalid email or password")
	}

	// Get user by their email
	user, err := useCase.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", err
	}

	token, err := middleware.GenerateTokenJwt(user.Id, user.Email, user.Roles, user.Status, 1)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (useCase *userUC) GetUserByEmail(email string) (userDto.User, error) {
	return useCase.userRepo.GetUserByEmail(email)
}

func (useCase *userUC) CreateDetailDebitur(debitur debiturFormDto.Debitur) error {
	return useCase.userRepo.CreateDetailDebitur(debitur)
}

func (useCase *userUC) CreateDetailDc(dc dcFormDto.DetailDC) error {
	return useCase.userRepo.CreateDetailDc(dc)
}

func (dbt *userUC) GetFullname(userId int) (string, error) {
	return dbt.userRepo.GetFullname(userId)
}

func (useCase *userUC) GetUserDataById(userId string) (interface{}, error) {
	roleId, err := useCase.userRepo.GetRolesById(userId)
	if err != nil {
		return nil, err
	}

	detailData, err := useCase.userRepo.GetUserDetailByUserId(userId)
	if err != nil {
		return nil, err
	}

	// return if user = debt collector
	if roleId == "3" {
		return detailData, nil
	}

	jobData, err := useCase.userRepo.GetUserJobDetailByUserId(userId)
	if err != nil {
		return nil, err
	}

	mrgcContact, err := useCase.userRepo.GetEmergencyContactByUserId(userId)
	if err != nil {
		return nil, err
	}

	return usersEntity.DetailedUserData{
		PersonalData:     detailData,
		EmploymentData:   jobData,
		EmergencyContact: mrgcContact,
	}, nil
}

func (useCase *userUC) UpdatePhotoPaths(userId int, fotoKTP, fotoSelfie string) error {
	return useCase.userRepo.UpdatePhotoPaths(userId, fotoKTP, fotoSelfie)
}

func (useCase *userUC) GetDataByRole(role, status string, page, size int) ([]debiturFormDto.DetailDebitur, int, error) {
    offset := (page - 1) * size
    debitur, totalData, err := useCase.userRepo.GetDataByRole(role, status, size, offset)
    if err != nil {
        return nil, 0, err
    }
    return debitur, totalData, nil
}

func (useUC *userUC) UpdateBankAccount(userId int, accountNumber, bankName string) error {
	exists, err := useUC.userRepo.IsBankAccExist(userId, accountNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("account number already exist, add another account number")
	}
	
	return useUC.userRepo.UpdateBankAccount(userId, accountNumber, bankName)
}

func (useUC *userUC) IsBankAccExist(userId int, accountNumber string) (bool, error) {
	return useUC.userRepo.IsBankAccExist(userId, accountNumber)
}