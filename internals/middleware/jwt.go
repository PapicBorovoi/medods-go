package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/PapicBorovoi/medods-go/api"
	"github.com/golang-jwt/jwt/v5"
)

type idKey string
const idKeyVal = idKey("id")
var ErrUnauthorized = fmt.Errorf("unauthorized")

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		})

		if err != nil || !token.Valid {
			api.RequestErrorHandler(w, ErrUnauthorized)
			return
		}

		id := token.Claims.(jwt.MapClaims)["id"].(string)

		next.ServeHTTP(w, r.WithContext(setID(r.Context(), id)))
	})
}

func GetID(ctx context.Context) string {
	return ctx.Value(idKeyVal).(string)
}

func setID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idKeyVal, id)
}