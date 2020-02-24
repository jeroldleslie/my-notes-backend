package api

import (
	"bytes"
	"errors"
	"fmt"
	//"mime/multipart"
)

//TODO reorganize this errors
var ErrNoteNotFound = errors.New("NOTE_NOT_FOUND")
var ErrBadInput = errors.New("ERR_BAD_INPUT")
var ErrNoDataAvailable = errors.New("NO_DATA_AVAILABLE")
var ErrNotFound = errors.New("ERR_NOT_FOUND")
var ErrNoData = fmt.Errorf("ERR_NO_DATA")

type API struct {
	DBMapper *DBMapper
	/* Redis    *redis.Client */
}

func (api *API) Signup(user *User) (*User, error) {
	return api.DBMapper.Signup(user)
}

func (api *API) Signin(cred *User) (*LoginResponse, error) {
	return api.DBMapper.Signin(cred)
}

func (api *API) CreateNote(note *Note) (*Note, error) {
	return api.DBMapper.SaveNote(note)
}

func (api *API) UpdateNote(note *Note) (*Note, error) {
	return api.DBMapper.UpdateNote(note)
}

func (api *API) GetUserNotes(userID int) (*[]Note, error) {
	return api.DBMapper.GetUserNotes(userID)
}

func (api *API) DeleteNotes(id int) error {
	return api.DBMapper.DeleteNotes(id)
}

func (api *API) Upload(noteID int64, fileName string, buffer *bytes.Buffer) error {
	return api.DBMapper.SaveFile(noteID, fileName, buffer)
}
