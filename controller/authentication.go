package controller

import (
	"e-commerce_api/database"
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	Buyer Role = iota
	Seller
)

type User struct {
	Username string
	Password string
	Email    string
	Role     Role
}

type Response struct {
	Message string `json:"message"`
}

func GetStringRole(role Role) string {
	if role > 1 || role < 0 {
		return ""
	}
	roles := []string{"Buyer", "Seller"}
	return roles[role]
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: err.Error(),
		})
		return
	}

	if !IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Email is not valid",
		})
		return
	}

	//Check username, email uniqueness
	var tmp string
	query := `SELECT username FROM users WHERE username=$1 OR email=$2`
	err = database.DB.Get(&tmp, query, user.Username, user.Email)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "User with this username or email already exists.",
		})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: err.Error(),
		})
		return
	}

	//Store into database
	query = `INSERT INTO users (ID, username, hashedPassword, email, role) VALUES ($1, $2, $3, $4, $5)`
	tx := database.DB.MustBegin()
	tx.MustExec(query, uuid.NewString(), user.Username, hashedPassword, user.Email, GetStringRole(user.Role))
	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Account successfully registered!",
	})
}
