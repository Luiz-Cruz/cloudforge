package environment_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestApp(handler *environment.Handler) *fiber.App {
	app := fiber.New()
	environment.RegisterRoutes(app, handler)
	return app
}

func TestHandler_CreateEnvironment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	expectedState := environment.EnvironmentState{
		TransactionID: "tx-abc",
		Name:          "prod-cluster",
		Type:          "kubernetes",
		Status:        "PENDING",
	}

	mockSvc.EXPECT().StartProvisioning(gomock.Any(), "prod-cluster", "kubernetes").
		Return(expectedState, nil).Times(1)

	payload := map[string]string{"name": "prod-cluster", "type": "kubernetes"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/environments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestHandler_CreateEnvironment_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	req := httptest.NewRequest(http.MethodPost, "/environments", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandler_CreateEnvironment_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	mockSvc.EXPECT().StartProvisioning(gomock.Any(), "prod-cluster", "kubernetes").
		Return(environment.EnvironmentState{}, errors.New("provisioning failed")).Times(1)

	payload := map[string]string{"name": "prod-cluster", "type": "kubernetes"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/environments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestHandler_GetEnvironment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	expectedState := &environment.EnvironmentState{
		TransactionID: "tx-123",
		Name:          "prod-cluster",
		Status:        "AVAILABLE",
	}

	mockSvc.EXPECT().GetEnvironment(gomock.Any(), "tx-123").Return(expectedState, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/environments/tx-123", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandler_GetEnvironment_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	mockSvc.EXPECT().GetEnvironment(gomock.Any(), "tx-unknown").Return(nil, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/environments/tx-unknown", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_GetEnvironment_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	mockSvc.EXPECT().GetEnvironment(gomock.Any(), "tx-123").Return(nil, errors.New("db down")).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/environments/tx-123", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestHandler_ListEnvironments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	mockSvc.EXPECT().ListEnvironments(gomock.Any()).
		Return([]*environment.EnvironmentState{{TransactionID: "tx-1"}}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/environments", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandler_ListEnvironments_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler := environment.NewHandler(mockSvc)
	app := setupTestApp(handler)

	mockSvc.EXPECT().ListEnvironments(gomock.Any()).
		Return(nil, errors.New("scan error")).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/environments", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
