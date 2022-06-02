package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type staff struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func isLoggedIn(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func sendLoginJwtCookie(c echo.Context) error {
	var (
		api = c.Get("api").(*api)

		passwordHash string
	)

	u := new(staff)
	if err := c.Bind(u); err != nil {
		return err
	}

	if err := api.db.QueryRow(context.Background(), api.q["get-password-hash"], u.Username).Scan(&passwordHash); err != nil {
		return c.String(http.StatusNotFound, "username not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password)); err != nil {
		return c.String(http.StatusForbidden, "login failed")
	}

	expiration := time.Now().Add(24 * 7 * time.Hour)

	claims := &jwtClaims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(api.jwtKey)
	if err != nil {
		return err
	}

	cookie := cookieDefaults()
	cookie.Value = tokenString
	cookie.Expires = expiration

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "login successful")
}

func sendLogoutCookie(c echo.Context) error {
	cookie := cookieDefaults()
	cookie.MaxAge = -1

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "logout successful")
}

func cookieDefaults() *http.Cookie {
	cookie := new(http.Cookie)

	cookie.Name = "_ge_auth"
	cookie.Path =  "/"
	cookie.SameSite = 3
	cookie.HttpOnly = true
	cookie.Secure = false

	return cookie
}