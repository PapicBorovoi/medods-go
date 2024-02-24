package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PapicBorovoi/medods-go/api"
	"github.com/PapicBorovoi/medods-go/internals/db"
	"github.com/PapicBorovoi/medods-go/internals/middleware"
	"github.com/PapicBorovoi/medods-go/internals/tools"
)

var ErrInvalidMethod = fmt.Errorf("invalid method")

func Handle(mux *http.ServeMux) {
	mux.Handle("/jwt",  http.HandlerFunc(createToken))
	mux.Handle("/protected", middleware.JWT(http.HandlerFunc(refresh)))
}

func createToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		api.RequestErrorHandler(w, ErrInvalidMethod)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		api.RequestErrorHandler(w, fmt.Errorf("ID is required"))
		return
	}
	

	token, refreshToken , err := tools.CreateTokens(id)

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	db.Delete(id);

	err = db.Create(db.Auth{ID: id, RefreshToken: refreshToken});

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	response := api.GetJWTPairResponse{
		AccessToken: token,
		RefreshToken: refreshToken,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}
}

func refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		api.RequestErrorHandler(w, ErrInvalidMethod)
		return
	}

	tokenString := r.Header.Get("Authorization")[len("Bearer "):]
	id := middleware.GetID(r.Context())

	result, err := db.Read(id)

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	if result.RefreshToken != tokenString {
		api.RequestErrorHandler(w, fmt.Errorf("invalid refresh token"))
		return
	} 

	db.Delete(id)

	token, refreshToken , err := tools.CreateTokens(id)

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	err = db.Create(db.Auth{ID: id, RefreshToken: refreshToken})

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	response := api.GetJWTPairResponse{
		AccessToken: token,
		RefreshToken: refreshToken,
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}
}