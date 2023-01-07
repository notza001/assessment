package expenses_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/notza001/assessment/expenses"
	"github.com/stretchr/testify/assert"
)

var sqlMock sqlmock.Sqlmock
var sqlDb *sql.DB

var (
	ExpenseJson = `{
		"title": "strawberry smoothie 3",
		"amount": 4444,
		"note": "night market promotion discount 100 bath", 
		"tags": ["beverage"]
	}`
)

func TestCreate(t *testing.T) {

	teardown := setupTest(t)
	defer teardown(t)

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(ExpenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	expenses.Db = sqlDb

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	sqlMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

	// Assertions
	if assert.NoError(t, expenses.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

}

func TestGetAll(t *testing.T) {

	teardown := setupTest(t)
	defer teardown(t)

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	expenses.Db = sqlDb

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(
		1, "title", 100, "note", pq.Array([]string{"tag1", "tag2"}))
	sqlMock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expensess").ExpectQuery().WillReturnRows(rows)

	// Assertions
	if assert.NoError(t, expenses.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func TestGetById(t *testing.T) {

	teardown := setupTest(t)
	defer teardown(t)

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	expenses.Db = sqlDb

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(
		1, "title", 100, "note", pq.Array([]string{"tag1", "tag2"}))
	sqlMock.ExpectPrepare(regexp.QuoteMeta("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")).ExpectQuery().WithArgs("1").WillReturnRows(rows)

	// Assertions
	if assert.NoError(t, expenses.GetById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func setupTest(tb testing.TB) func(tb testing.TB) {
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
