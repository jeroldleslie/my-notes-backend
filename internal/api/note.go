package api

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

type Note struct {
	ID          int64     `json:"id"`
	Title       string    `sql:"title" json:"title" form:"title"`
	Content     string    `sql:"content" json:"content" form:"content"`
	Priority    string    `sql:"priority" json:"priority" form:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RemindFrom  time.Time `json:"remind_from"`
	RemindUntil time.Time `json:"remind_until"`
	Image       string    `sql:"image" json:"image" form:"image"`
	UserID      int64     `sql:"user_id" json:"user_id"`
	ImageID     int64     `sql:"-" json:"image_id"`
}

func (c *Note) MarshalJSON() ([]byte, error) {
	type Alias Note
	noTime := time.Time{}
	o := &struct {
		*Alias
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		RemindFrom    string `json:"remind_from"`
		RemindUntil   string `json:"remind_until"`
		EffectiveFrom string `json:"effective_from,omitempty"`
	}{
		Alias:       (*Alias)(c),
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
		RemindFrom:  c.RemindFrom.Format(time.RFC3339),
		RemindUntil: c.RemindUntil.Format(time.RFC3339),
	}

	if c.CreatedAt == noTime {
		o.CreatedAt = ""
	}
	if c.UpdatedAt == noTime {
		o.UpdatedAt = ""
	}
	if c.RemindFrom == noTime {
		o.RemindFrom = ""
	}
	if c.RemindUntil == noTime {
		o.RemindUntil = ""
	}

	return json.Marshal(o)
}

func (m *DBMapper) SaveNote(note *Note) (*Note, error) {

	now := time.Now().UTC()
	note.CreatedAt = now
	note.UpdatedAt = now
	if insertErr := m.DB.Insert(note); insertErr != nil {
		return nil, errors.Wrapf(insertErr, "couldn't insert note %+v", note)
	}

	return note, nil

}

func (m *DBMapper) GetUserNotes(userID int) (*[]Note, error) {
	/* select *, f.id as image_id from notes n left join files f on f.note_id=n.id; */
	notes := &[]Note{}
	err := m.DB.Model(notes).
		Column("note.*").
		ColumnExpr("f.id as image_id").
		Join("LEFT JOIN files as f").JoinOn("note.id=f.note_id").
		Where("user_id = ?", userID).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return notes, nil
		}
		return nil, err
	}
	return notes, nil
}

func (m *DBMapper) DeleteNotes(id int) error {

	fileMod := m.DB.Model((*File)(nil))
	fileMod = fileMod.Where("note_id = ?", id)
	_, ferr := fileMod.Delete()
	if ferr != nil {
		return ferr
	}

	mod := m.DB.Model((*Note)(nil))
	mod = mod.Where("id = ?", id)
	_, err := mod.Delete()

	if err != nil {
		return err
	}
	return nil
}
