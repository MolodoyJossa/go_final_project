package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTKey() []byte {
	return []byte(os.Getenv("TODO_SECRET"))
}

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func passwordHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func NewToken(password string) (string, error) {
	hash := passwordHash(password)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["hash"] = hash
	claims["exp"] = time.Now().Add(8 * time.Hour).Unix()

	return token.SignedString(getJWTKey())
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(signinResponse{Error: "Некорректный запрос"})
		return
	}
	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" || req.Password != expected {
		json.NewEncoder(w).Encode(signinResponse{Error: "Неверный пароль"})
		return
	}
	token, err := NewToken(req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(signinResponse{Error: "Ошибка генерации токена"})
		return
	}
	json.NewEncoder(w).Encode(signinResponse{Token: token})
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return getJWTKey(), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["hash"] != passwordHash(pass) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
