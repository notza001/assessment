package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetById(c echo.Context) error {

	db := open(DbUrl)
	defer db.Close()
	id := c.Param("id")

	Expense, err := getExpense(db, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "item not found")
	}

	return c.JSON(http.StatusOK, Expense)
}

func getExpense(db *sql.DB, id string) (Expense, error) {
	var Expense = Expense{}

	sqlstm, err := db.Prepare(`SELECT id, title, amount, note, tags  FROM expenses WHERE id = $1`)
	if err != nil {
		return Expense, err
	}

	query := sqlstm.QueryRow(id)

	err = query.Scan(&Expense.Id, &Expense.Title, &Expense.Amount, &Expense.Note, pq.Array(&Expense.Tags))

	if err != nil {
		return Expense, err
	}

	return Expense, nil
}
