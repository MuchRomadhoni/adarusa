package json

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	// JSONResponse - struct for json response success
	jsonResponse struct {
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	jsonResponseWithPaging struct {
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Paging  *Paging     `json:"paging,omitempty"`
	}

	// JSONResponse - struct for json response error
	jsonErrorResponse struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	Paging struct {
		Page      int `json:"page,omitempty"`
		TotalData int `json:"totalData,omitempty"`
	}

	ValidationField struct {
		FieldName string `json:"field"`
		Message   string `json:"message"`
	}

	jsonBadRequestResponse struct {
		Message          string            `json:"message"`
		ErrorDescription []ValidationField `json:"error_description,omitempty"`
	}

	JwtClaim struct {
		jwt.StandardClaims
		UserId string
		Email string `json:"email"`
		Roles string `json:"role"`
		Status string `json:"status"`
	}
)

func NewResponseSuccess(c *gin.Context, result interface{}, message string) {
	c.JSON(http.StatusOK, jsonResponse{
		Message: message,
		Data:    result,
	})
}

func NewResponseSuccessWithPaging(c *gin.Context, result interface{}, paging Paging, message string) {
	c.JSON(http.StatusOK, jsonResponseWithPaging{
		Message: message,
		Data:    result,
		Paging:  &paging,
	})
}

func NewResponseBadRequestValidator(c *gin.Context, validationField []ValidationField, message string) {
	c.JSON(http.StatusBadRequest, jsonBadRequestResponse{
		Message:          message,
		ErrorDescription: validationField,
	})
}

func NewResponseBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, jsonResponse{
		Message: message,
	})
}

func NewResponseError(c *gin.Context, err string) {
	log.Error().Msg(err)
	c.JSON(http.StatusInternalServerError, jsonErrorResponse{
		Message: "internal server error",
		Error:   err,
	})
}

func NewResponseNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, jsonResponse{
		Message: message,
	})
}

func NewResponseForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, jsonResponse{
		Message: message,
	})
}

func NewResponseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, jsonResponse{
		Message: message,
	})
}

func NewAbortForbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, jsonResponse{
		Message: message,
	})
}

func NewAbortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, jsonResponse{
		Message: message,
	})
}
