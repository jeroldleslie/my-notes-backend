package api

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	tableName struct{} `sql:"user" json:"-"`
	Email     string   `sql:"email" json:"email,omitempty" form:"email"`
	Password  string   `sql:"password" json:"password,omitempty" form:"password"`
}

type LoginResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty",`
	User    string `json:"user,omitempty"`
	UserID  int64  `json:"user_id,omitempty"`
}

func (m *DBMapper) Signin(cred *User) (*LoginResponse, error) {
	loginResponse := &LoginResponse{}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	user := &User{}
	err := m.DB.Model(user).
		Where("email = ?", cred.Email).
		First()

	if err != nil {
		if err == pg.ErrNoRows {
			loginResponse.Status = "FAILURE"
			loginResponse.Message = "Email not found"
			return loginResponse, nil
		}
		return nil, err
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		loginResponse.Status = "FAILURE"
		loginResponse.Message = "Incorrect password"
		return loginResponse, nil
	} else {
		type NotesClaims struct {
			UserID int64  `json:"id"`
			Name   string `json:"name"`
			Email  string `json:"email"`
			jwt.StandardClaims
		}

		claims := NotesClaims{
			user.ID,
			user.Name,
			user.Email,
			jwt.StandardClaims{
				ExpiresAt: expiresAt,
				Issuer:    "notes",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signingKey := []byte("notesthebestalltime")

		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			loginResponse.Status = "FAILURE"
			loginResponse.Message = "cannot create jwt token"
			return loginResponse, nil
		}

		loginResponse.Status = "SUCCESS"
		loginResponse.Message = "logged in"
		loginResponse.Token = tokenString
		loginResponse.User = user.Name
		loginResponse.UserID = user.ID
		return loginResponse, nil
	}

}
