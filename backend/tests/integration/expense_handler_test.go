package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/luisfucros/expense-tracker-app/internal/config"
	"github.com/luisfucros/expense-tracker-app/internal/handler"
	dbrepo "github.com/luisfucros/expense-tracker-app/internal/repository/db"
	"github.com/luisfucros/expense-tracker-app/internal/router"
	"github.com/luisfucros/expense-tracker-app/internal/service"
	"github.com/luisfucros/expense-tracker-app/tests/integration/testhelpers"
)

func newIntegrationSetup(t *testing.T) (http.Handler, string) {
	t.Helper()

	db, err := testhelpers.SetupTestDB()
	if err != nil {
		t.Skipf("skipping integration test, DB unavailable: %v", err)
	}

	cfg := &config.Config{
		Env:            "test",
		JWTSecret:      "integration-test-secret",
		JWTExpiryHours: 1,
	}

	userRepo := dbrepo.NewUserRepository(db)
	expenseRepo := dbrepo.NewExpenseRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg)
	expenseSvc := service.NewExpenseService(expenseRepo)

	uniqueEmail := fmt.Sprintf("setup+%d@example.com", time.Now().UnixNano())
	authResp, err := authSvc.Register(t.Context(), registerInput(uniqueEmail))
	require.NoError(t, err)

	h := handler.NewHandler(authSvc, expenseSvc, noopLogger())
	r := router.New(cfg, h, noopLogger())

	return r, authResp.Token
}

func TestCreateExpense_Unauthorized(t *testing.T) {
	db, err := testhelpers.SetupTestDB()
	if err != nil {
		t.Skipf("skipping integration test, DB unavailable: %v", err)
	}

	cfg := &config.Config{
		Env:            "test",
		JWTSecret:      "integration-test-secret",
		JWTExpiryHours: 1,
	}

	userRepo := dbrepo.NewUserRepository(db)
	expenseRepo := dbrepo.NewExpenseRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg)
	expenseSvc := service.NewExpenseService(expenseRepo)
	h := handler.NewHandler(authSvc, expenseSvc, noopLogger())
	r := router.New(cfg, h, noopLogger())

	body := `{"title":"Coffee","amount":4.5,"category":"Leisure","date":"2024-01-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateExpense_Success(t *testing.T) {
	r, token := newIntegrationSetup(t)

	body := `{"title":"Coffee","amount":4.5,"category":"Leisure","date":"2024-01-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Coffee", data["title"])
	assert.Equal(t, "Leisure", data["category"])
}

func TestListExpenses_Success(t *testing.T) {
	r, token := newIntegrationSetup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	_, hasExpenses := data["expenses"]
	assert.True(t, hasExpenses)
}

func TestDeleteExpense_Success(t *testing.T) {
	r, token := newIntegrationSetup(t)

	// Create an expense first
	body := `{"title":"To Delete","amount":10.0,"category":"Others","date":"2024-01-20"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var createResp map[string]interface{}
	err := json.NewDecoder(createW.Body).Decode(&createResp)
	require.NoError(t, err)
	data := createResp["data"].(map[string]interface{})
	expenseID := data["id"].(float64)

	// Delete it
	deleteReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/expenses/%.0f", expenseID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusNoContent, deleteW.Code)
}
