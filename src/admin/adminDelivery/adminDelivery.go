package adminDelivery

import (
	"fp_pinjaman_online/model/dto/adminDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/pkg/middleware"
	adminInterface "fp_pinjaman_online/src/admin"
	"strconv"

	"github.com/gin-gonic/gin"
)

type adminDelivery struct {
	adminUc adminInterface.AdminUsecase
}

func NewAdminDelivery(v1Group *gin.RouterGroup, adminUc adminInterface.AdminUsecase) {
	handler := adminDelivery{
		adminUc: adminUc,
	}
	adminGroup := v1Group.Group("/admin")
	adminGroup.Use(middleware.JWTAuthWithRoles("admin"))
	{
		adminGroup.PUT("/:id/verify", handler.VerifyAndUpdateUser)
		adminGroup.POST("/verify-pinjaman", handler.VerifyAndCreateCicilan)
		adminGroup.POST("/verify-tugas", handler.VerifyAndSendBalanceDC)
		adminGroup.POST("/withdrawal", handler.VerifyWithdrawalDC)
	}
}

func (a *adminDelivery) VerifyAndUpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	var req adminDto.RequestUpdateStatusUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	req.ID = id

	res, err := a.adminUc.VerifyAndUpdateUser(req)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, res, "success")
}

func (a *adminDelivery) VerifyAndCreateCicilan(ctx *gin.Context) {
	var req adminDto.RequestVerifyLoan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	res, err := a.adminUc.VerifyAndCreateCicilan(req)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, res, "success")
}

func (a *adminDelivery) VerifyAndSendBalanceDC(ctx *gin.Context) {
	var req adminDto.RequestUpdateClaimTugas
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}
	res, err := a.adminUc.VerifyAndSendBalanceDC(req)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, res, "success")

}

func (a *adminDelivery) VerifyWithdrawalDC(ctx *gin.Context) {
	var req adminDto.RequestWithdrawal
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	res, err := a.adminUc.VerifyWithdrawalDC(req)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}
	json.NewResponseSuccess(ctx, res, "success")
}
