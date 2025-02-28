package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `sql:"name" json:"name,omitempty" form:"name"`
	Email    string `sql:"email" json:"email,omitempty" form:"email"`
	Password string `sql:"password" json:"password,omitempty" form:"password"`
}

func (m *DBMapper) Signup(user *User) (*User, error) {
	b, _ := json.Marshal(user)
	fmt.Println(string(b))

	fmt.Println(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	if err := m.DB.Insert(user); err != nil {
		return nil, errors.WithStack(err)
	}
	user.Password = ""
	return user, nil
}
