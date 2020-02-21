package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jeroldleslie/my-notes-backend/internal/api"
	"github.com/jeroldleslie/my-notes-backend/internal/log_conf"
	"github.com/jeroldleslie/my-notes-backend/internal/utils"
	"github.com/jeroldleslie/my-notes-backend/internal/utils/db_logger"
	"github.com/jeroldleslie/my-notes-backend/internal/utils/stage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// APIHandler handles API routes
type APIHandler struct {
	API *api.API
}

const (
	version = "v1"
	port    = "8000"
)

func main() {
	if !stage.IsProd() {
		if err := log_conf.AllLevelFiles(".", "my-notes-api", logrus.TraceLevel); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}

	db, err := utils.ConnectToPostgres("")
	if err != nil {
		fmt.Errorf("%+v", err)
		panic(err)
	}
	db_logger.AddDbLogger(db, true)

	mapper := api.DBMapper{
		DB: db,
	}
	mapper.InitCache()
	defer db.Close()
	api := &APIHandler{
		API: &api.API{
			DBMapper: &api.DBMapper{
				DB: db,
			},
		},
	}

	serve(api)
}

func serve(a *APIHandler) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/v1/ping", ping)
	e.Logger.Fatal(e.Start(":" + port))
}

func ping(c echo.Context) error {
	response := utils.Response{
		StatusCode: http.StatusOK,
		Data:       http.StatusText(http.StatusOK),
	}
	return response.Send(&c)
}

func errorToHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Errorf("%+v", err)
	switch err {
	case api.ErrNoteNotFound:
		return http.StatusNotFound
	case api.ErrBadInput:
		return http.StatusBadRequest
	case api.ErrNoData:
		return http.StatusNotFound // Request from QAs
		//return http.StatusNoContent
	}
	return http.StatusInternalServerError
}
