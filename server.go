package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notza001/assessment/expenses"
)

func sayHello(c echo.Context) error {
	return c.JSON(http.StatusCreated, "Hello")
}

func main() {

	e := echo.New()
	// url := os.Getenv("DATABASE_URL")
	url := os.Getenv("DATABASE_URL")
	expenses.InitTable(url)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth != "November 10, 2009" {
				return c.NoContent(http.StatusUnauthorized)
			}
			return next(c)
		}
	})

	e.POST("/expenses", expenses.Create)
	e.GET("/expenses", expenses.GetAll)
	e.PUT("/expenses/:id", expenses.Update)
	e.GET("/expenses/:id", expenses.GetById)

	go func() {
		err := e.Start(os.Getenv("PORT"))
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	fmt.Println("application started.")
	fmt.Printf("port number: %s\n", os.Getenv("PORT"))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
