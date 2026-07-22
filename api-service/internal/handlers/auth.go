package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"api-service/internal/config"
	"api-service/internal/models"
	"api-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID         string `json:"id"`
	TelegramID string `json:"telegramId"`
	Role       string `json:"role"`
}

type TokenClaims struct {
	UserID     string `json:"userId"`
	TelegramID string `json:"telegramId"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Login == "" || req.Password == "" {
		SendError(w, http.StatusBadRequest, "Login and password are required")
		return
	}

	if len(req.Password) < 6 {
		SendError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	repo := repository.NewRepository[models.User]()
	filters := map[string]interface{}{
		"telegramId": req.Login,
	}

	existingUsers, _, err := repo.GetAll(filters, "", 1, 0, []string{})
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to query user: "+err.Error())
		return
	}

	if len(existingUsers) > 0 {
		SendError(w, http.StatusConflict, "User already exists")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	newUser := models.User{
		TelegramId:   req.Login,
		PasswordHash: string(passwordHash),
		Role:         models.RoleUser,
	}

	err = repo.Create(&newUser)
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	cfg := config.LoadConfig()
	claims := TokenClaims{
		UserID:     newUser.ID,
		TelegramID: newUser.TelegramId,
		Role:       string(newUser.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	SendJSON(w, http.StatusCreated, RegisterResponse{
		Token: tokenString,
		User: UserInfo{
			ID:         newUser.ID,
			TelegramID: newUser.TelegramId,
			Role:       string(newUser.Role),
		},
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Login == "" || req.Password == "" {
		SendError(w, http.StatusBadRequest, "Login and password are required")
		return
	}

	repo := repository.NewRepository[models.User]()
	filters := map[string]interface{}{
		"telegramId": req.Login,
	}

	users, _, err := repo.GetAll(filters, "", 1, 0, []string{})
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to query user: "+err.Error())
		return
	}

	if len(users) == 0 {
		SendError(w, http.StatusUnauthorized, "Invalid login or password")
		return
	}

	user := users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		SendError(w, http.StatusUnauthorized, "Invalid login or password")
		return
	}

	cfg := config.LoadConfig()
	claims := TokenClaims{
		UserID:     user.ID,
		TelegramID: user.TelegramId,
		Role:       string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	SendJSON(w, http.StatusOK, LoginResponse{
		Token: tokenString,
		User: UserInfo{
			ID:         user.ID,
			TelegramID: user.TelegramId,
			Role:       string(user.Role),
		},
	})
}

// AdminMiddleware проверяет JWT токен или X-Telegram-Id и роль администратора
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		telegramId := r.Header.Get("X-Telegram-Id")

		// 1. Пробуем JWT
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
			cfg := config.LoadConfig()
			token, err := jwt.ParseWithClaims(bearerToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*TokenClaims); ok && claims.Role == string(models.RoleAdmin) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// 2. Фолбек на X-Telegram-Id (для обратной совместимости и локальной разработки)
		if telegramId != "" {
			repo := repository.NewRepository[models.User]()
			users, _, err := repo.GetAll(map[string]interface{}{"telegramId": telegramId}, "", 1, 0, nil)
			if err == nil && len(users) > 0 {
				if users[0].Role == models.RoleAdmin {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		SendError(w, http.StatusUnauthorized, "Missing or invalid authentication")
	})
}