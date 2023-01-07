//go:build integration
// +build integration

package expenses_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/notza001/assessment/expenses"
	"github.com/stretchr/testify/assert"
)

var serverPort = 80
var dbUrl = "postgres://root:root@db/go-test-db?sslmode=disable"

func TestIntegrationCreate(t *testing.T) {

	ExpenseJson := strings.NewReader(`{
		"title": "strawberry smoothie 12",
		"amount": 4,
		"note": "night market promotion discount 100 bath", 
		"tags": ["beverage"]
	}`)

	expenses.InitTable(dbUrl)
	teardown := setupIntegrationTest(t)
	defer teardown(t)

	// Setup
	e := echo.New()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/", serverPort), ExpenseJson)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, expenses.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

}

func TestIntegrationGetAll(t *testing.T) {

	expenses.InitTable(dbUrl)
	teardown := setupIntegrationTest(t)
	defer teardown(t)

	// Setup
	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/", serverPort), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, expenses.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func setupIntegrationTest(tb testing.TB) func(tb testing.TB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		tb.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlDb = db
	sqlMock = mock

	return func(tb testing.TB) {
		db.Close()
	}
}
