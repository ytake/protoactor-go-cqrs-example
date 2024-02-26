package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-cqrs-example/internal/action"
	"github.com/ytake/protoactor-go-cqrs-example/internal/config"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/internal/logger"
	"github.com/ytake/protoactor-go-cqrs-example/internal/route"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(logger.MiddlewareFactory(
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil)).
			With("env", "production")))
	db, err := mysql.NewConn(config.MysqlConfig().FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	defer db.Close()
	actor, err := route.NewRestAPIActorSystem(db)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Simple Domain Sample!")
	})
	ah, err := action.NewUserRegistration(actor)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.POST("/user/registration", ah.Handle)
	e.GET("/users", action.NewUserList(mysql.NewUserFindStore(db)).Handle)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":1323"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
