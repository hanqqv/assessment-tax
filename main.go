package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hanqqv/assessment-tax/postgres"
	"github.com/hanqqv/assessment-tax/tax"
)

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	handler := tax.New(p)
	admin := e.Group("/admin")

	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil
	}))

	e.POST("/tax/calculations", handler.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", handler.CalculateTaxCSVHandler)
	admin.POST("/deductions/personal", handler.SettingPersonalDeductionHandler)

	port := os.Getenv("PORT")

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
