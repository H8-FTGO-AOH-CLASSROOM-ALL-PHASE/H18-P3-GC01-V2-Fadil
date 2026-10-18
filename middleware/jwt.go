package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken signs a JWT valid for 24 hours.
func GenerateToken(secret string, userID int64, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// SecretLookup returns the JWT secret that belongs to the given user.
type SecretLookup func(ctx context.Context, userID int64) (string, error)

// JWT validates the "Authorization: Bearer <token>" header using the
// per-user secret of the user named in the token claims.
func JWT(lookup SecretLookup) echo.MiddlewareFunc {
	unauthorized := echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get(echo.HeaderAuthorization)
			tokenString, found := strings.CutPrefix(header, "Bearer ")
			if !found || tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed bearer token")
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				// The claim is only used to find the secret; the signature is still verified with it.
				c, ok := t.Claims.(*Claims)
				if !ok || c.UserID == 0 {
					return nil, errors.New("missing user id")
				}
				secret, err := lookup(context.Background(), c.UserID)
				if err != nil || secret == "" {
					return nil, errors.New("unknown user")
				}
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err != nil || !token.Valid {
				return unauthorized
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
