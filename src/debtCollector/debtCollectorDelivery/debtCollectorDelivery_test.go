package debtCollectorDelivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fp_pinjaman_online/model/dto/debtCollectorDto"
	dtoJson "fp_pinjaman_online/model/dto/json"
	"fp_pinjaman_online/model/entity/debtCollectorEntity"
	mymock "fp_pinjaman_online/src/debtCollector/debtCollectorUseCase/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define the test suite
type DebtCollectorDeliverySuite struct {
	suite.Suite
	mockUC *mymock.DebtCollectorUseCaseMock
	router *gin.Engine
}

// SetupSuite is run once before the suite's tests are run
func (s *DebtCollectorDeliverySuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest is run before each test in the suite
func (s *DebtCollectorDeliverySuite) SetupTest() {
	s.mockUC = &mymock.DebtCollectorUseCaseMock{Mock: mock.Mock{}}
	s.router = gin.Default()
	v1Group := s.router.Group("/v1")
	NewDebtCollectorDelivery(v1Group, s.mockUC)
}

func TestDebtCollectorDeliverySuite(t *testing.T) {
	suite.Run(t, new(DebtCollectorDeliverySuite))
}

func (s *DebtCollectorDeliverySuite) TestAddLogTugas_Success() {
	payload := debtCollectorDto.NewLogTugasPayload{TugasId: "1", Description: "Test Task"}
	s.mockUC.Mock.On("CreateLogTugas", payload).Return(nil)

	body, err := json.Marshal(payload)
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPost, "/v1/debt-collector/log-tugas/create", bytes.NewBuffer(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "success")
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestAddLogTugas_Fail() {
	completePayload := debtCollectorDto.NewLogTugasPayload{TugasId: "1", Description: "Test Task"}
	invalidPayload := debtCollectorDto.NewLogTugasPayload{TugasId: "salah", Description: "Test Task"}
	s.Run("error bad request validator", func() {
		body, err := json.Marshal(invalidPayload)
		s.NoError(err)

		req, err := http.NewRequest(http.MethodPost, "/v1/debt-collector/log-tugas/create", bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
		s.Contains(w.Body.String(), "bad request")
		s.mockUC.Mock.AssertExpectations(s.T())
	})

	s.Run("error not found", func() {
		expError := errors.New("log not found")
		s.mockUC.Mock.On("CreateLogTugas", completePayload).Return(expError)

		body, err := json.Marshal(completePayload)
		s.NoError(err)

		req, err := http.NewRequest(http.MethodPost, "/v1/debt-collector/log-tugas/create", bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusNotFound, w.Code)
		s.Contains(w.Body.String(), "not found")
		s.mockUC.Mock.AssertExpectations(s.T())
	})
}

func (s *DebtCollectorDeliverySuite) TestEditLogTugas_Success() {
	param := debtCollectorDto.Param{ID: "1"}
	payload := debtCollectorDto.UpdateLogTugasPayload{Description: "Test Task"}
	s.mockUC.Mock.On("EditLogTugasById", param.ID, payload).Return(nil)

	body, err := json.Marshal(payload)
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPut, "/v1/debt-collector/log-tugas/1", bytes.NewBuffer(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "success")
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestEditLogTugas_BindJSONError() {
	payload := map[string]interface{}{
		"description": 123,
	}

	body, err := json.Marshal(payload)
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPut, "/v1/debt-collector/log-tugas/1", bytes.NewBuffer(body))
	s.NoError(err)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "invalid payload")
}

func (s *DebtCollectorDeliverySuite) TestEditLogTugas_BindParamError() {
	req, err := http.NewRequest(http.MethodPut, "/v1/debt-collector/log-tugas/s", nil)
	s.NoError(err)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "bad request")
}

func (s *DebtCollectorDeliverySuite) TestEditLogTugas_NotFoundError() {
	param := debtCollectorDto.Param{ID: "1"}
	payload := debtCollectorDto.UpdateLogTugasPayload{Description: "Test Task"}
	expErr := errors.New("data not found")
	s.mockUC.Mock.On("EditLogTugasById", param.ID, payload).Return(expErr)

	body, err := json.Marshal(payload)
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPut, "/v1/debt-collector/log-tugas/1", bytes.NewBuffer(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "not found")
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestDeleteLogTugas_Success() {
	param := debtCollectorDto.Param{ID: "1"}
	s.mockUC.Mock.On("DeleteLogTugasById", param.ID).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/debt-collector/log-tugas/1", nil)
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "success")
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestDeleteLogTugas_BindParamError() {
	req, err := http.NewRequest(http.MethodDelete, "/v1/debt-collector/log-tugas/s", nil)
	s.NoError(err)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "bad request")
}

func (s *DebtCollectorDeliverySuite) TestDeleteLogTugas_NotFoundError() {
	param := debtCollectorDto.Param{ID: "1"}
	expErr := errors.New("data not found")
	s.mockUC.Mock.On("DeleteLogTugasById", param.ID).Return(expErr)

	req, err := http.NewRequest(http.MethodDelete, "/v1/debt-collector/log-tugas/1", nil)
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "not found")
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestGetAllLogTugas_Success() {
	param := debtCollectorDto.Param{ID: "1"}
	expPaging := dtoJson.Paging{
		Page:      1,
		TotalData: 10,
	}
	expLogList := []debtCollectorEntity.LogTugas{
		{
			ID:          "1",
			TugasId:     "2",
			Description: "tes desc",
		},
	}
	s.mockUC.Mock.On("GetAllLogTugas", param.ID, 1, 10).Return(expLogList, expPaging, nil)
	req, err := http.NewRequest(http.MethodGet, "/v1/debt-collector/tugas/1/log-tugas?page=1&size=10", nil)
	s.NoError(err, "error request")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.mockUC.Mock.AssertExpectations(s.T())
}

func (s *DebtCollectorDeliverySuite) TestGetAllLogTugas_BindUriError() {
	req, err := http.NewRequest(http.MethodGet, "/v1/debt-collector/tugas/invalid/log-tugas", nil)
	s.NoError(err)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "bad request")
}

func (s *DebtCollectorDeliverySuite) TestGetAllLogTugas_BindParamError() {
	req, err := http.NewRequest(http.MethodGet, "/v1/debt-collector/tugas/1/log-tugas?page=s&size=salah", nil)
	s.NoError(err)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "bad request")
}

func (s *DebtCollectorDeliverySuite) TestGetAllLateDebtor_Success() {
	// query := debtCollectorDto.Query{}
	expPaging := dtoJson.Paging{
		Page:      1,
		TotalData: 10,
	}
	expLateDebtorList := []debtCollectorEntity.LateDebtor{
		{
			ID:           "1",
			FullName:     "nama tes",
			Address:      "jl. tes",
			UnpaidAmount: 1000000,
		},
	}
	s.mockUC.Mock.On("GetAllLateDebtorByCity", "1", 0, 0).Return(expLateDebtorList, expPaging, nil)
	req, err := http.NewRequest(http.MethodGet, "/v1/debt-collector/late-debitur", nil)
	s.NoError(err, "error request")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.mockUC.Mock.AssertExpectations(s.T())
}
