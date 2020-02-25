package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type File struct {
	ID       int64  `json:"id"`
	FileName string `sql:"file_name" json:"file_name"`
	Content  []byte `sql:"content" json:"content"`
	ContentType string    `sql:"content_type" json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	NoteID      int64     `sql:"note_id" json:"note_id"`
}

func (m *DBMapper) SaveFile(noteID int64, fileName string, buffer *bytes.Buffer) error {
	file := &File{}

	err := m.DB.Model(file).Where("note_id = ?", noteID).First()
	if err == nil {
		now := time.Now().UTC()
		file.UpdatedAt = now
		file.FileName = fileName
		file.Content = buffer.Bytes()
		file.ContentType = http.DetectContentType(buffer.Bytes())
		if uErr := m.DB.Update(file); uErr != nil {
			return errors.Wrapf(uErr, "couldn't update file")
		}Î
	} else {
		now := time.Now().UTC()
		file.CreatedAt = now
		file.UpdatedAt = now
		file.FileName = fileName
		file.Content = buffer.Bytes()
		file.ContentType = http.DetectContentType(buffer.Bytes())
		file.NoteID = noteID

		if insertErr := m.DB.Insert(file); insertErr != nil {
			return errors.Wrapf(insertErr, "couldn't insert file")
		}
	}

	return nil
}
