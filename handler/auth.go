package handler

import (
	"errors"
	"net/http"
	"strings"

	"H18-P3-GC01-V2-Fadil/middleware"
	"H18-P3-GC01-V2-Fadil/repository"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users     *repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(users *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{users: users, jwtSecret: jwtSecret}
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func bindCredentials(c echo.Context) (*credentials, error) {
	var req credentials
	if err := c.Bind(&req); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") || req.Password == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "valid email and password are required")
	}
	return &req, nil
}

// Register - POST /register
func (h *AuthHandler) Register(c echo.Context) error {
	req, err := bindCredentials(c)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &repository.User{Email: req.Email, Password: string(hash)}
	if err := h.users.Create(c.Request().Context(), user); err != nil {
		if errors.Is(err, repository.ErrEmailDuplicate) {
			return echo.NewHTTPError(http.StatusConflict, "email already registered")
		}
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "user registered successfully"})
}

// Login - POST /login
func (h *AuthHandler) Login(c echo.Context) error {
	req, err := bindCredentials(c)
	if err != nil {
		return err
	}

	invalid := echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	user, err := h.users.FindByEmail(c.Request().Context(), req.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return invalid
	}
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return invalid
	}

	token, err := middleware.GenerateToken(h.jwtSecret, user.ID, user.Email)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token, "token_type": "Bearer"})
}
