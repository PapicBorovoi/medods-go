package tools

import (
	"net/http"
	"os"
	"time"

	"github.com/PapicBorovoi/medods-go/internals/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func CreateTokens(id string) (string, string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	refreshTokenString, refreshErr := refreshToken.SignedString(
		[]byte(os.Getenv("JWT_REFRESH_SECRET")),
	)

	if err != nil {
		return "", "", err
	} else if refreshErr != nil {
		return "", "", refreshErr
	}

	return tokenString, refreshTokenString, nil
}

func VerifyToken(tokenString string, r *http.Request) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	id := middleware.GetID(r.Context())

	return id, nil
}