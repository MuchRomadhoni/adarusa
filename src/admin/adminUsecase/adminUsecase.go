package adminUsecase

import (
	"fmt"
	"fp_pinjaman_online/model/dto/adminDto"
	"fp_pinjaman_online/pkg/validation"
	adminInterface "fp_pinjaman_online/src/admin"
	"log"
	"math"
	"time"
)

type adminUsecase struct {
	repo adminInterface.AdminRepository
}

func NewAdminUsecase(repo adminInterface.AdminRepository) adminInterface.AdminUsecase {
	return &adminUsecase{repo: repo}
}

func (uc *adminUsecase) VerifyAndUpdateUser(req adminDto.RequestUpdateStatusUser) (adminDto.AdminResponse, error) {

	user, err := uc.repo.RetrieveUserStatusById(req.ID)
	if err != nil {
		log.Printf("Error retrieving user with ID %d: %v", req.ID, err)
		return adminDto.AdminResponse{}, err
	}

	if user == nil {
		log.Printf("No user found with ID %d", req.ID)
		return adminDto.AdminResponse{}, fmt.Errorf("no user found with ID %d", req.ID)
	}

	if req.Status == "verified" {
		if user.RoleID == 2 && !validation.ValidateUserComplete(*user) {
			log.Printf("Verification failed for user ID %d: incomplete user information", req.ID)
			return adminDto.AdminResponse{}, fmt.Errorf("verification failed: missing some data information for user ID %d", req.ID)
		} else if user.RoleID == 3 {
			log.Printf("Debt collector with ID %d verified with minimal data", req.ID)
		}
	} else {
		log.Printf("No verification needed for user ID %d with status %s", req.ID, req.Status)
	}

	if user.Status == req.Status {
		return adminDto.AdminResponse{}, fmt.Errorf("user with ID %d already verified", req.ID)
	}

	if user.Status != req.Status {
		err := uc.repo.UpdateUserStatus(req.ID, req.Status)
		if err != nil {
			log.Printf("Failed to update status for user ID %d: %v", req.ID, err)
			return adminDto.AdminResponse{}, err
		}

		if req.Status == "verified" {
			now := time.Now()
			user.VerifiedAt = &now
			err = uc.repo.InsertLimitId(req.LimitID, req.ID)
			if err != nil {
				return adminDto.AdminResponse{}, fmt.Errorf("failed to insert limit ID %d: %v", req.LimitID, err)
			}
		}
	}
	return adminDto.AdminResponse{
		ID:     user.UserID,
		Email:  user.Email,
		Status: "verified",
	}, nil
}

func (uc *adminUsecase) VerifyAndCreateCicilan(req adminDto.RequestVerifyLoan) (adminDto.LoanResponse, error) {
	loan, err := uc.repo.RetrievePinjamanById(req.LoanID)
	if err != nil {
		log.Printf("Error retrieving loan: %v", err)
		return adminDto.LoanResponse{}, fmt.Errorf("error retrieving loan: %w", err)
	}

	if loan.StatusPengajuan == req.StatusPengajuan {
		log.Printf("No action needed: loan %d is already %s", loan.ID, req.StatusPengajuan)
		return adminDto.LoanResponse{
			LoanID:  loan.ID,
			UserID:  loan.UserID,
			Message: fmt.Sprintf("Loan is already %s", req.StatusPengajuan),
		}, nil
	}

	user, err := uc.repo.RetrieveUserLimitByAdmin(req.UserID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return adminDto.LoanResponse{}, fmt.Errorf("error retrieving user: %w", err)
	}

	if user.MaxPinjaman.Float64 < loan.JumlahPinjaman {
		err = uc.repo.UpdateLoanStatus(loan.ID, "rejected") // Auto-reject due to credit limit
		if err != nil {
			log.Printf("Error auto-rejecting loan: %v", err)
			return adminDto.LoanResponse{}, fmt.Errorf("error auto-rejecting loan: %w", err)
		}
		log.Printf("Loan %d auto-rejected due to credit limit", loan.ID)
		return adminDto.LoanResponse{
			LoanID:  loan.ID,
			UserID:  loan.UserID,
			Message: "Loan rejected due to exceeding credit limit",
		}, nil
	}

	err = uc.repo.UpdateLoanStatus(loan.ID, req.StatusPengajuan)
	if err != nil {
		log.Printf("Error updating loan status: %v", err)
		return adminDto.LoanResponse{}, fmt.Errorf("error updating loan status: %w", err)
	}

	if req.StatusPengajuan == "approved" {
		// Create payment schedule as loan is approved
		monthlyPayment := CalculateMonthlyPayment(loan.JumlahPinjaman, loan.BungaPerBulan*12, loan.Tenor)
		dueDate := time.Now()
		for i := 1; i <= loan.Tenor; i++ {
			dueDate = dueDate.AddDate(0, 1, 0)
			if err := uc.repo.InsertCicilan(loan.ID, dueDate, monthlyPayment, "unpaid"); err != nil {
				return adminDto.LoanResponse{}, fmt.Errorf("error creating payment schedule: %w", err)
			}
		}
		log.Printf("Payment schedule created for loan %d", loan.ID)
		return adminDto.LoanResponse{
			LoanID:  loan.ID,
			UserID:  loan.UserID,
			Message: "Payment schedule created successfully",
		}, nil
	}

	return adminDto.LoanResponse{
		LoanID:  loan.ID,
		UserID:  loan.UserID,
		Message: "Loan status updated successfully",
	}, nil
}

func (uc *adminUsecase) VerifyAndSendBalanceDC(req adminDto.RequestUpdateClaimTugas) (adminDto.ClaimTugasResponse, error) {
	claimTugas, err := uc.repo.RetrieveTugasById(req.ID)
	if err != nil {
		return adminDto.ClaimTugasResponse{}, fmt.Errorf("error retrieving claim task: %v", err)
	}

	if claimTugas.StatusTugas == "done" {
		return adminDto.ClaimTugasResponse{}, fmt.Errorf("claim task is already marked as done")
	}

	err = uc.repo.UpdateClaimTugas(req.ID, "done")
	if err != nil {
		return adminDto.ClaimTugasResponse{}, fmt.Errorf("failed to update claim task status: %v", err)
	}

	const rewardAmount = 500000
	err = uc.repo.UpdateBalance(claimTugas.CollectorID, rewardAmount)
	if err != nil {
		return adminDto.ClaimTugasResponse{}, fmt.Errorf("failed to update balance: %v", err)
	}

	response := adminDto.ClaimTugasResponse{
		ID:      req.ID,
		Status:  "done",
		Message: "Claim task updated to done and balance dc increased by 500,000.",
	}
	return response, nil
}

func (uc *adminUsecase) VerifyWithdrawalDC(req adminDto.RequestWithdrawal) (adminDto.WithdrawalResponse, error) {
	withdrawal, err := uc.repo.RetrieveWithdrawalById(req.ID)
	if err != nil {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("error retrieving withdrawal: %v", err)
	}

	if withdrawal.UserID != req.UserID {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("user ID does not match the withdrawal record")
	}

	if withdrawal.Status == "rejected" {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("withdrawal is rejected and cannot be processed")
	}

	if withdrawal.Status != "pending" {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("withdrawal is not in a state that can be processed, current status: %s", withdrawal.Status)
	}

	err = uc.repo.UpdateWithdrawalStatus(withdrawal.ID, "paid")
	if err != nil {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("failed to update withdrawal status to paid: %v", err)
	}

	// Update the user's balance after successfully updating the withdrawal status.
	err = uc.repo.UpdateBalance(withdrawal.UserID, -withdrawal.Amount)
	if err != nil {
		return adminDto.WithdrawalResponse{}, fmt.Errorf("failed to update user balance: %v", err)
	}

	response := adminDto.WithdrawalResponse{
		ID:     withdrawal.ID,
		UserID: withdrawal.UserID,
		Amount: withdrawal.Amount,
		Status: "paid",
	}

	return response, nil

}

func CalculateMonthlyPayment(principal float64, annualInterestRate float64, tenor int) float64 {
	monthlyInterestRate := (annualInterestRate / 100) / 12

	numberOfPayments := float64(tenor)

	if monthlyInterestRate == 0 {
		return principal / numberOfPayments
	}
	monthlyPayment := principal * (monthlyInterestRate * math.Pow(1+monthlyInterestRate, numberOfPayments)) / (math.Pow(1+monthlyInterestRate, numberOfPayments) - 1)

	return math.Ceil(monthlyPayment)
}
