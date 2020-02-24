package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-pg/pg/orm"
	"github.com/jeroldleslie/my-notes-backend/internal/api"
	"github.com/jeroldleslie/my-notes-backend/internal/utils"

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
	/* if !stage.IsProd() {
		if err := log_conf.AllLevelFiles(".", "my-notes-api", logrus.TraceLevel); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	} */

	db, err := utils.ConnectToPostgres("")
	if err != nil {
		fmt.Errorf("%+v", err)
		panic(err)
	}
	//db_logger.AddDbLogger(db, true)

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/api/ping", ping)
	e.POST("/api/signup", a.handleSignup)
	e.POST("/api/signin", a.handleSignin)
	g := e.Group("/api/notes")

	g.POST("", a.handleCreateNote)
	g.GET("/user_notes/:user_id", a.handleGetNote)
	g.GET("/search", a.handleSearch)
	g.PUT("/:id", a.handleUpdateNote)
	g.DELETE("/:id", a.handleDeleteNote)
	g.POST("/file", a.handleUpload)
	g.GET("/file/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		file := &api.File{}
		a.API.DBMapper.DB.Model(file).Where("id = ?", id).Select()

		c.Response().Header().Set("Content-Type", file.ContentType)
		c.Response().Header().Set("Content-Length", strconv.Itoa(len(file.Content)))

		return c.Blob(http.StatusOK, file.ContentType, file.Content)

	})

	e.Logger.Fatal(e.Start(":" + port))
}

func ping(c echo.Context) error {
	data := "{ \"message\": \"hello angular I am from API. Hello leslie too\"}"
	response := utils.Response{
		StatusCode: http.StatusOK,
		Data:       data, //http.StatusText(http.StatusOK),
	}
	return response.Send(&c)
}

func (a *APIHandler) handleSignup(c echo.Context) error {
	//data := "{ \"message\": \"hello angular I am from API. Hello leslie too\"}"

	user := api.User{}
	var userResult *api.User
	err := c.Bind(&user)
	if err == nil {
		fmt.Println(">>>>>>")
		fmt.Println(user.Password)
		userResult, err = a.API.Signup(&user)
	}

	user.Password = ""
	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(err),
		Data:       userResult,
		Error:      err,
	}
	return response.Send(&c)
}

func (a *APIHandler) handleSignin(c echo.Context) error {
	//data := "{ \"message\": \"hello angular I am from API. Hello leslie too\"}"

	cred := api.User{}
	err := c.Bind(&cred)
	if err != nil {
		response := utils.Response{
			StatusCode: errorToHTTPStatusCode(err),
			Data:       nil,
			Error:      err,
		}
		return response.Send(&c)
	}

	b, _ := json.Marshal(cred)
	fmt.Println(string(b))
	fmt.Println(string(b))
	fmt.Println(string(b))
	fmt.Println(string(b))
	fmt.Println(string(b))
	fmt.Println(string(b))

	loginResponse, errf := a.API.Signin(&cred)

	if errf != nil {
		response := utils.Response{
			StatusCode: errorToHTTPStatusCode(errf),
			Data:       nil,
			Error:      errf,
		}
		return response.Send(&c)
	}
	response := utils.Response{
		StatusCode: http.StatusOK,
		Data:       loginResponse,
		Error:      nil,
	}
	return response.Send(&c)
}

func (a *APIHandler) handleCreateNote(c echo.Context) error {

	note := api.Note{}
	err := c.Bind(&note)
	if err != nil {
		response := utils.Response{
			StatusCode: errorToHTTPStatusCode(err),
			Data:       nil,
			Error:      err,
		}
		return response.Send(&c)
	}
	savedNote, errf := a.API.CreateNote(&note)

	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(errf),
		Data:       savedNote,
		Error:      errf,
	}

	return response.Send(&c)
}

func (a *APIHandler) handleGetNote(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	notes, err := a.API.GetUserNotes(userID)

	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(err),
		Data:       notes,
		Error:      err,
	}
	return response.Send(&c)
}

func (a *APIHandler) handleUpdateNote(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	note := api.Note{}
	err := c.Bind(&note)
	if err != nil {
		response := utils.Response{
			StatusCode: errorToHTTPStatusCode(err),
			Data:       nil,
			Error:      err,
		}
		return response.Send(&c)
	}
	note.ID = int64(id)
	savedNote, errf := a.API.UpdateNote(&note)

	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(errf),
		Data:       savedNote,
		Error:      errf,
	}

	return response.Send(&c)
}
func (a *APIHandler) handleDeleteNote(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	deleteResponse := "{\"status\":\"SUCCESS\"}"
	err := a.API.DeleteNotes(id)
	if err != nil {
		deleteResponse = "{\"status\":\"FAILURE\"}"
	}
	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(err),
		Data:       deleteResponse,
		Error:      err,
	}
	return response.Send(&c)
}

func (a *APIHandler) handleUpload(c echo.Context) error {
	//noteID, _ := strconv.Atoi(c.Param("note_id"))

	//-----------
	// Read file
	//-----------

	// Source
	noteID, _ := strconv.Atoi(c.FormValue("note_id"))
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	bodyBuf := &bytes.Buffer{}

	//iocopy
	_, err = io.Copy(bodyBuf, src)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(bodyBuf.Bytes())
	fmt.Println(contentType)
	fmt.Println(file.Size)
	fmt.Println(len(bodyBuf.Bytes()))

	err = a.API.Upload(int64(noteID), file.Filename, bodyBuf)

	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(err),
		Data:       noteID,
		Error:      err,
	}
	return response.Send(&c)
}

func (a *APIHandler) handleSearch(c echo.Context) error {
	searchText := c.QueryParam("text")
	priority := c.QueryParam("priority")
	userID := c.QueryParam("user_id")

	notes := &[]api.Note{}

	mod := a.API.DBMapper.DB.Model(notes).
		Column("note.*").
		ColumnExpr("f.id as image_id").
		Join("LEFT JOIN files as f").JoinOn("note.id=f.note_id").
		Where("note.user_id = ?", userID)

	if strings.TrimSpace(searchText) != "" {
		mod = mod.
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.Where("note.content ILIKE ?", "%"+searchText+"%").
					WhereOr("note.title ILIKE ?", "%"+searchText+"%")
				return q, nil
			})
	}

	if strings.TrimSpace(priority) != "" && priority != "NONE" {
		mod = mod.Where("note.priority = ?", priority)
	}
	err := mod.Order("note.created_at DESC").Select()
	response := utils.Response{
		StatusCode: errorToHTTPStatusCode(err),
		Data:       notes,
		Error:      err,
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
