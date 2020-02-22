package api

import (
	"fmt"
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
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty",`
	User    string `json:"user,omitempty"`
}

func (m *DBMapper) Signin(cred *User) (*LoginResponse, error) {
	fmt.Println("inside signin >>>>>>>>>>>>>>>>>>>>>>>>>>")

	fmt.Println(cred.Email)
	loginResponse := &LoginResponse{}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	user := &User{}
	err := m.DB.Model(user).
		Where("email = ?", cred.Email).
		First()

	if err != nil {
		if err == pg.ErrNoRows {
			loginResponse.Status = false
			loginResponse.Message = "Email not found. Please register if you have not registered yet."
			return loginResponse, nil
		}
		return nil, err
	}

	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)
	fmt.Printf("%s=%s", user.Password, cred.Password)

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		loginResponse.Status = false
		loginResponse.Message = "Invalid login credentials. Please try again"
		//var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
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
			loginResponse.Status = false
			loginResponse.Message = "cannot create jwt token"
			return loginResponse, nil
		}

		loginResponse.Status = true
		loginResponse.Message = "logged in"
		loginResponse.Token = tokenString
		loginResponse.User = user.Name
		return loginResponse, nil
	}

}
