package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func Create(c echo.Context) error { // Path: expenses/create.go
	var expense Expense
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := open(DbUrl)
	defer db.Close()

	query := `INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).Scan(&expense.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, expense)

}
