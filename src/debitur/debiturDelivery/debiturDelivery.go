package debiturDelivery

import (
	"fp_pinjaman_online/model/dto/debiturDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/pkg/middleware"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/src/debitur"
	"strconv"

	"github.com/gin-gonic/gin"
)

type debiturDelivery struct {
	debiturUC debitur.DebiturUsecase
}

func NewDebiturDelivery(v1Group *gin.RouterGroup, debiturUC debitur.DebiturUsecase) {
	handler := debiturDelivery{
		debiturUC: debiturUC,
	}
	usersGroup := v1Group.Group("/users")
	debiturGroup := usersGroup.Group("/debitur")
	debiturGroup.Use(middleware.JWTAuthWithRoles("debitur"), middleware.VerifiedOnly())
	{
		debiturGroup.POST("/create/pinjaman", handler.PengajuanPinjaman)
		debiturGroup.GET("/pinjaman", handler.GetPengajuanPinjaman)
		debiturGroup.GET("/cicilan/:id", handler.GetCicilan)
		debiturGroup.POST("/cicilan/pay", handler.CicilanPay)
		debiturGroup.GET("/cicilan/verify/:id", handler.CicilanVerify)
	}
}

func (u *debiturDelivery) PengajuanPinjaman(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		json.NewResponseUnauthorized(c, "unauthorized")
		return
	}

	userID, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(c, "invalid userID")
		return
	}

	var req debiturDto.DebiturRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(c, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(c, "tipe data salah")
		return
	}
	if req.Tenor < 1 {
		json.NewResponseError(c, "tenor must be greater than 0")
		return
	}
	if req.Tenor > 12 {
		json.NewResponseError(c, "tenor must be less than or equal to 12")
		return
	}
	err = u.debiturUC.PengajuanPinjaman(userID, req.JumlahPinjaman, req.Tenor, req.Description)
	if err != nil {
		json.NewResponseError(c, err.Error())
		return
	}
	json.NewResponseSuccess(c, nil, "success")
}

func (u *debiturDelivery) GetPengajuanPinjaman(c *gin.Context) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		json.NewResponseUnauthorized(c, "unauthorized")
		return
	}

	id, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(c, "invalid userID")
		return
	}

	data, err := u.debiturUC.GetPengajuanPinjaman(id)
	if err != nil {
		json.NewResponseError(c, err.Error())
		return
	}

	json.NewResponseSuccess(c, data, "success")
}

func (u *debiturDelivery) GetCicilan(c *gin.Context) {
	id := c.Param("id") //pinjaman ID
	status := c.Query("status")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	data, paging, err := u.debiturUC.GetCicilan(page, size, offset, id, status)
	if err != nil {
		json.NewResponseError(c, err.Error())
		return
	}
	json.NewResponseSuccessWithPaging(c, data, paging, "success")
}

func (u *debiturDelivery) CicilanPay(c *gin.Context) {
	var req debiturDto.CicilanPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(c, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(c, "tipe data salah")
		return
	}
	data, err := u.debiturUC.CicilanPayment(req.PinjamanId, req.TotalBayar)
	if err != nil {
		json.NewResponseError(c, err.Error())
		return
	}
	json.NewResponseSuccess(c, data, "success")
}

func (u *debiturDelivery) CicilanVerify(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := u.debiturUC.CicilanVerify(id)
	if err != nil {
		json.NewResponseError(c, err.Error())
		return
	}
	json.NewResponseSuccess(c, nil, "success")
}
