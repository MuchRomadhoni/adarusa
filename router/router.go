package router

import (
	"database/sql"

	"fp_pinjaman_online/src/admin/adminDelivery"
	"fp_pinjaman_online/src/admin/adminUsecase"
	adminRepository "fp_pinjaman_online/src/admin/repository"

	"fp_pinjaman_online/src/debitur/debiturDelivery"
	"fp_pinjaman_online/src/debitur/debiturRepository"
	"fp_pinjaman_online/src/debitur/debiturUsecase"

	"fp_pinjaman_online/src/debtCollector/debtCollectorDelivery"
	"fp_pinjaman_online/src/debtCollector/debtCollectorRepository"
	"fp_pinjaman_online/src/debtCollector/debtCollectorUseCase"

	"fp_pinjaman_online/src/users/userDelivery"
	"fp_pinjaman_online/src/users/userRepository"
	"fp_pinjaman_online/src/users/userUseCase"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {

	//debitur
	debiturRepository := debiturRepository.NewDebiturRepository(db)
	debiturUC := debiturUsecase.NewDebiturUsecase(debiturRepository)
	debiturDelivery.NewDebiturDelivery(v1Group, debiturUC)

	adminRepo := adminRepository.NewAdminRepository(db)
	adminUC := adminUsecase.NewAdminUsecase(adminRepo)
	adminDelivery.NewAdminDelivery(v1Group, adminUC)

	userRepository := userRepository.NewUserRepository(db)
	userUC := userUseCase.NewUserUseCase(userRepository)
	userDelivery.NewUserDelivery(v1Group, userUC)

	debtCollectorRepo := debtCollectorRepository.NewDebtCollectorRepository(db)
	debtCollectorUC := debtCollectorUseCase.NewDebtCollectorUseCase(debtCollectorRepo, userRepository)
	debtCollectorDelivery.NewDebtCollectorDelivery(v1Group, debtCollectorUC)
}
