package debtCollectorDelivery

import (
	"fp_pinjaman_online/model/dto/debtCollectorDto"
	"fp_pinjaman_online/model/dto/json"

	"fp_pinjaman_online/pkg/middleware"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/src/debtCollector"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type debtCollectorDelivery struct {
	debtCollUC debtCollector.DebtCollectorUseCase
}

func NewDebtCollectorDelivery(v1Group *gin.RouterGroup, debtCollUC debtCollector.DebtCollectorUseCase) {
	handler := debtCollectorDelivery{
		debtCollUC: debtCollUC,
	}
	dcGroup := v1Group.Group("/debt-collector")
	dcGroup.Use(middleware.JWTAuthWithRoles("dc"), middleware.VerifiedOnly())
	{
		dcGroup.GET("/late-debitur/:id", handler.GetLateDebtor)      // get late debitur info from tugas
		dcGroup.GET("/late-debitur", handler.GetAllLateDebtor)       // get all debitur nunggak
		dcGroup.POST("/tugas/create", handler.AddTugas)              // claim tugas
		dcGroup.GET("/tugas", handler.GetAllTugas)                   // get all tugas atau user yang pernah di tagih
		dcGroup.GET("/tugas/:id/log-tugas", handler.GetAllLogTugas)  // get all log
		dcGroup.POST("/log-tugas/create", handler.AddLogTugas)       // membuat log tugas baru
		dcGroup.GET("/log-tugas/:id", handler.GetLogTugas)           // get log detail
		dcGroup.PUT("/log-tugas/:id", handler.EditLogTugas)          // edit log
		dcGroup.DELETE("/log-tugas/:id", handler.DeleteLogTugas)     // hapus log
		dcGroup.GET("/balance", handler.GetBalance)                  // menampilkan saldo atau gaji
		dcGroup.POST("/balance/withdraw", handler.CreateWithdrawReq) // withdraw saldo
	}
}

// TODO
// add upload file for job proof ?
func (d *debtCollectorDelivery) AddLogTugas(ctx *gin.Context) {
	var payload debtCollectorDto.NewLogTugasPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(ctx, "invalid payload")
		return
	}

	err := d.debtCollUC.CreateLogTugas(payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (d *debtCollectorDelivery) GetLogTugas(ctx *gin.Context) {
	var param debtCollectorDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	log, err := d.debtCollUC.GetLogTugasById(param.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, log, "success")
}

func (d *debtCollectorDelivery) EditLogTugas(ctx *gin.Context) {
	var param debtCollectorDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	var payload debtCollectorDto.UpdateLogTugasPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(ctx, "invalid payload")
		return
	}

	err := d.debtCollUC.EditLogTugasById(param.ID, payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (d *debtCollectorDelivery) DeleteLogTugas(ctx *gin.Context) {
	var param debtCollectorDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	err := d.debtCollUC.DeleteLogTugasById(param.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (d *debtCollectorDelivery) GetAllLogTugas(ctx *gin.Context) {
	var param debtCollectorDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	var queryParams debtCollectorDto.Query
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}
	page, _ := strconv.Atoi(queryParams.Page)
	size, _ := strconv.Atoi(queryParams.Size)

	logsList, paging, err := d.debtCollUC.GetAllLogTugas(param.ID, page, size)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	if len(logsList) == 0 {
		json.NewResponseSuccess(ctx, nil, "data not found")
		return
	}

	json.NewResponseSuccessWithPaging(ctx, logsList, paging, "")
}

func (d *debtCollectorDelivery) GetAllLateDebtor(ctx *gin.Context) {
	var queryParams debtCollectorDto.Query
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}
	page, _ := strconv.Atoi(queryParams.Page)
	size, _ := strconv.Atoi(queryParams.Size)

	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	lateDebtorsList, paging, err := d.debtCollUC.GetAllLateDebtorByCity(dcId.(string), page, size)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	if len(lateDebtorsList) == 0 {
		json.NewResponseSuccess(ctx, nil, "data not found")
		return
	}

	json.NewResponseSuccessWithPaging(ctx, lateDebtorsList, paging, "")
}

func (d *debtCollectorDelivery) AddTugas(ctx *gin.Context) {
	var payload debtCollectorDto.NewTugasPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(ctx, "invalid payload")
		return
	}

	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	err := d.debtCollUC.ClaimTugas(dcId.(string), payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		if strings.Contains(err.Error(), "maximum") {
			json.NewResponseBadRequest(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (d *debtCollectorDelivery) GetAllTugas(ctx *gin.Context) {
	var queryParams debtCollectorDto.Query
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	page, _ := strconv.Atoi(queryParams.Page)
	size, _ := strconv.Atoi(queryParams.Size)

	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	listTugas, paging, err := d.debtCollUC.GetAllTugas(dcId.(string), queryParams.Status, page, size)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	if len(listTugas) == 0 {
		json.NewResponseSuccess(ctx, nil, "data not found")
		return
	}

	json.NewResponseSuccessWithPaging(ctx, listTugas, paging, "")
}

func (d *debtCollectorDelivery) GetBalance(ctx *gin.Context) {
	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	balance, err := d.debtCollUC.GetBalanceByUserId(dcId.(string))
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, map[string]float64{"balance": balance}, "")
}

func (d *debtCollectorDelivery) CreateWithdrawReq(ctx *gin.Context) {
	var payload debtCollectorDto.WithdrawalReqPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
		json.NewResponseBadRequest(ctx, "invalid payload")
		return
	}

	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	err := d.debtCollUC.CreateWithdrawRequest(dcId.(string), payload.Amount)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (d *debtCollectorDelivery) GetLateDebtor(ctx *gin.Context) {
	var param debtCollectorDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	dcId, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseError(ctx, "failed to get user id")
		return
	}

	data, err := d.debtCollUC.GetDebtorData(param.ID, dcId.(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, data, "success")
}
