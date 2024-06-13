package userDelivery

import (
	"bytes"
	jsonEndoce "encoding/json"
	"errors"
	"fp_pinjaman_online/model/userDto"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/src/users/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserDeliverySuite struct {
	suite.Suite
	mokUC  *mocks.UserUseCase
	router *gin.Engine
}

func (s *UserDeliverySuite) SetupSuite() {
	s.mokUC = &mocks.UserUseCase{Mock: mock.Mock{}}
	s.router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validation.ValidationPassword)
	}

	v1Group := s.router.Group("/v1")
	NewUserDelivery(v1Group, s.mokUC)
}

func TestUserDeliverySuite(t *testing.T) {
	suite.Run(t, new(UserDeliverySuite))
}

func (s *UserDeliverySuite) TestLogin_Success() {
	reqBody := userDto.LoginRequest{
		Email: "test@example.com",
		Password: "password",
	}
	jsonValue, _ := jsonEndoce.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/users/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	token := "some_token"
	s.mokUC.On("Login", reqBody).Return(token, nil)

	s.router.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var response map[string]interface{}
	err := jsonEndoce.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["message"])
	assert.Equal(s.T(), token, response["data"].(map[string]interface{})["token"])
}

func (s *UserDeliverySuite) TestLogin_Fail() {
	reqBody := userDto.LoginRequest{
		Email: "test@example.com",
		Password: "wrong_password",
	}
	jsonValue, _ := jsonEndoce.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/users/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	s.mokUC.On("Login", reqBody).Return("", errors.New("internal server error"))

	s.router.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	var response map[string]interface{}
	err := jsonEndoce.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "internal server error", response["message"])
}

func (s *UserDeliverySuite) TestCreateUser_Success() {
	reqBody := userDto.CreateRequest{
		Email:    "test@example.com",
		Password: "Password1!",
	}
	jsonValue, _ := jsonEndoce.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/users/debitur/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	s.mokUC.On("CreateUser", reqBody, 2).Return(nil)

	s.router.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var response map[string]interface{}
	err := jsonEndoce.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["message"])
}

func (s *UserDeliverySuite) TestCreateUser_Failure() {
	reqBody := userDto.CreateRequest{
		Email:    "test@example.com",
		Password: "Password1!",
		Name:     "Test User",
	}
	jsonValue, _ := jsonEndoce.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/users/debitur/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	s.mokUC.On("CreateUser", reqBody, 2).Return(errors.New("email already exists"))

	s.router.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	var response map[string]interface{}
	err := jsonEndoce.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "email already exists", response["message"])
}