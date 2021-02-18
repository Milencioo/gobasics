package models

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	tokenSecret = []byte(os.Getenv("token_secret"))
)
//User struct
type User struct {
	ID              uuid.UUID `json : "id"`
	CreatedAt       time.Time `json : "_"`
	UpdatedAt       time.Time `json : "_"`
	Email           string    `json : "email" `
	PasswordHash    string    `json: "_" `
	Password        string    `json: "password"`
	PasswordConfirm string    `json: "password_confirm"`
}
//Register registration
func (u *User) Register(conn *pgx.Conn) error {
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("password err")
	}
	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("password need to be match")
	}
	if len(u.Email) < 5 {
		return fmt.Errorf("email err")
	}
	u.Email = strings.ToLower(u.Email)
	row := conn.QueryRow(context.Background(), "SELECT id from user_account WHERE email = $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	if err != pgx.ErrNoRows {
		fmt.Println("found user")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("user exist")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Error("err")
	}
	u.PasswordHash = string(pwdHash)
	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, email, password_hash) 
	VALUES ($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash) 

	return err

}

//GetAuthToken auth
func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour = 25).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err


}
//IsAuthenticated authentication 
func (u *User) IsAuthenticated( conn *pgx.Conn ) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hash from 
	user_account WHERE email=$1", u.Email)
	err := row.Scan(&u.ID, &u.PaswordHash)
	if err == pgx.ErrNoRows {
		fmt.Println("err")
		return fmt.Errorf("error")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("err")
	}
	return nil

}
