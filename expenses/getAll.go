package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetAll(c echo.Context) error {

	db := open(DbUrl)
	defer db.Close()

	query, err := db.Query(`SELECT * FROM expenses`)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, "item not found")
	}

	var response = []Expense{}
	for query.Next() {
		var expense Expense
		err := query.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		response = append(response, expense)
	}

	return c.JSON(http.StatusOK, response)
}
