package api

import (
	"errors"
	"fmt"
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
