package controllers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"simple-api/database"
	"simple-api/models"
	"time"
)

const (
	SecretKey  = "secret"
	CookieName = "jwt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, _ := bcrypt.GenerateFromPassword(user.Password, 10)
	user.Password = password
	if err = database.DB.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var receivedUser, dbUser models.User

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&receivedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.DB.Where("email = ?", receivedUser.Email).First(&dbUser)

	if dbUser.Id == 0 {
		http.Error(w, "wrong user/password", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword(dbUser.Password, receivedUser.Password); err != nil {
		http.Error(w, "wrong user/password", http.StatusNotFound)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    dbUser.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Expires:  time.Now().Add(-time.Second),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
