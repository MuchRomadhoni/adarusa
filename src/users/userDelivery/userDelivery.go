package userDelivery

import (
	"fmt"
	"fp_pinjaman_online/config/cloudinary"
	"fp_pinjaman_online/model/dcFormDto"
	"fp_pinjaman_online/model/debiturFormDto"
	"fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/userDto"
	"fp_pinjaman_online/pkg/middleware"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/src/users"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type userDelivery struct {
	userUC users.UserUseCase
}

func NewUserDelivery(v1Group *gin.RouterGroup, userUc users.UserUseCase) {
	handler := userDelivery{
		userUC: userUc,
	}
	userGroup := v1Group.Group("/users")
	userGroup.POST("/login", handler.login)
	userGroup.POST("/:role/create", handler.createUser)
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.POST("/upload/form", handler.uploadFiles)
		userGroup.POST("/rekening", handler.updateAccountNumber)
	}

	userGroup.POST("/debitur/form", middleware.JWTAuthWithRoles("debitur"), handler.createDetailDebitur)
	userGroup.POST("/dc/form", middleware.JWTAuthWithRoles("dc"), handler.createDetailDC)

	userGroup.Use(middleware.JWTAuthWithRoles("admin"))
	{
		userGroup.GET("/alldata/:roles", handler.getDataByRole)
		userGroup.GET("/data/:id", handler.GetUserDataById)
	}
}

func (c *userDelivery) createUser(ctx *gin.Context) {
	role := ctx.Param("role")

	var roleId int
	switch role {
	case "debitur":
		roleId = 2 // Debitur role ID
	case "dc":
		roleId = 3 // Debt collector role ID
	default:
		json.NewResponseBadRequest(ctx, "Invalid role")
		return
	}

	var req userDto.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	err := c.userUC.CreateUser(req, roleId)
	if err != nil {
		json.NewResponseBadRequest(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (c *userDelivery) login(ctx *gin.Context) {
	var req userDto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationErro := validation.GetValidationError(err)
		if len(validationErro) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationErro, "bad request")
			return
		}
	}

	token, err := c.userUC.Login(req)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}
	json.NewResponseSuccess(ctx, map[string]interface{}{"token": token}, "success")
}

func (c *userDelivery) createDetailDebitur(ctx *gin.Context) {
	userIdStr, exists := ctx.Get("userId")
	if !exists {
		json.NewAbortUnauthorized(ctx, "unauthorized")
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(ctx, "invalid userId format")
		return
	}

	var debt debiturFormDto.Debitur
	if err := ctx.ShouldBindJSON(&debt); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, "bad request")
			return
		}
	}

	// set userID from token JWT, no need input on body json
	debt.DetailUser.UserID = userId
	debt.UserJobs.UserID = userId
	debt.EmergencyContact.UserID = userId

	// Pengecekan jika userID dalam body berbeda dengan userID dari token
	if debt.DetailUser.UserID != userId {
		json.NewAbortForbidden(ctx, "forbidden: cannot modify another user's data")
		return
	}

	err = c.userUC.CreateDetailDebitur(debt)
	if err != nil {
		if err.Error() == "nik already exists" {
			json.NewResponseBadRequest(ctx, "nik already exists")
		} else {
			json.NewResponseError(ctx, err.Error())
		}
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (c *userDelivery) createDetailDC(ctx *gin.Context) {
	userIdStr, exists := ctx.Get("userId")
	if !exists {
		json.NewAbortUnauthorized(ctx, "unauthorized")
		return
	}
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(ctx, "invalid userId format")
		return
	}

	var dc dcFormDto.DetailDC
	if err := ctx.ShouldBindJSON(&dc); err != nil {
		validationErro := validation.GetValidationError(err)
		if len(validationErro) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationErro, "bad request")
			return
		}
	}

	// set userID from token JWT, no need input on body json
	dc.UserID = userId
	if dc.UserID != userId {
		json.NewAbortForbidden(ctx, "forbidden: cannot modify another user's data")
		return
	}

	err = c.userUC.CreateDetailDc(dc)
	if err != nil {
		if err.Error() == "nik already exists" {
			json.NewResponseBadRequest(ctx, "nik already exists")
		} else {
			json.NewResponseError(ctx, err.Error())
		}
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func (c *userDelivery) uploadFiles(ctx *gin.Context) {
	// get userId and roles from context
	roles, _ := ctx.Get("roleName")
	userIdStr, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseUnauthorized(ctx, "unauthorized")
		return
	}

	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(ctx, "invalid userID")
		return
	}

	fullname, err := c.userUC.GetFullname(userId)
	fmt.Println("fullname:", fullname)
	if err != nil {
		json.NewResponseError(ctx, "unable to fect user detail")
		return
	}

	// handle file upload
	fileKTP, err := ctx.FormFile("foto_ktp")
	if err != nil {
		json.NewResponseBadRequest(ctx, "no foto_ktp file is received")
		return
	}
	fileSelfie, err := ctx.FormFile("foto_selfie")
	if err != nil {
		json.NewResponseBadRequest(ctx, "no foto_selfie is received")
		return
	}

	ktpURL, err := uploadFileToCloudinary(ctx, fileKTP, roles.(string), fullname, "ktp")
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}
	selfieURL, err := uploadFileToCloudinary(ctx, fileSelfie, roles.(string), fullname, "selfie")
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	// update photo path in db
	err = c.userUC.UpdatePhotoPaths(userId, ktpURL, selfieURL)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}

func uploadFileToCloudinary(ctx *gin.Context, file *multipart.FileHeader, role, fullName, fileType string) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	uploadParams := uploader.UploadParams{
		Folder: fmt.Sprintf("uploads/%s/%s/%s", role, fullName, fileType),
	}

	uploadResp, err := cloudinary.Cloudinary.Upload.Upload(ctx, fileContent, uploadParams)
	if err != nil {
		return "", err
	}

	return uploadResp.SecureURL, nil
}

func (c *userDelivery) getDataByRole(ctx *gin.Context) {
	role := ctx.Param("roles")
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	status := ctx.DefaultQuery("status", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		json.NewResponseBadRequest(ctx, "bad request: invalid page parameter")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		json.NewResponseBadRequest(ctx, "bad request: invalid size parameter")
		return
	}

	debitur, totalData, err := c.userUC.GetDataByRole(role, status, page, size)
	if err != nil {
		json.NewResponseError(ctx, err.Error())
		return
	}
	if len(debitur) == 0 {
		json.NewResponseSuccess(ctx, "", "success")
		return
	}

	paging := json.Paging{Page: page, TotalData: totalData}

	json.NewResponseSuccessWithPaging(ctx, debitur, paging, "success")
}

func (c *userDelivery) GetUserDataById(ctx *gin.Context) {
	var param userDto.Params
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request")
			return
		}
	}

	resp, err := c.userUC.GetUserDataById(param.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error())
			return
		}
		json.NewResponseError(ctx, err.Error())
		return
	}

	json.NewResponseSuccess(ctx, resp, "success")
}

func (c *userDelivery) updateAccountNumber(ctx *gin.Context) {
	userIdStr, exists := ctx.Get("userId")
	if !exists {
		json.NewResponseUnauthorized(ctx, "unauthorized")
		return
	}

	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		json.NewResponseError(ctx, "invalid userID")
		return
	}

	var request userDto.CreateBankAccount
	if err := ctx.ShouldBindJSON(&request); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequestValidator(ctx, validationError, "bad request body json")
			return
		}
	}

	request.UserID = userId
	if request.UserID != userId {
		json.NewAbortForbidden(ctx, "forbidden")
		return
	}

	err = c.userUC.UpdateBankAccount(userId, request.AccountNumber, request.BankName)
	if err != nil {
		if err.Error() == "account number already exist, add another account number" {
			json.NewResponseBadRequest(ctx, "account number already exist, add another account number")
			return
		}
		json.NewResponseError(ctx, "bad request")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success")
}
