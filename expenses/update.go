package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func Update(c echo.Context) error {
	db := open(DbUrl)

	_, err_validate := getExpense(db, c.Param("id"))

	if err_validate != nil {
		return c.JSON(http.StatusNotFound, Err{Message: err_validate.Error()})
	}

	var expense = Expense{}
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	defer db.Close()

	query := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5`
	_, err = db.Exec(query, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags), expense.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, expense)

}
