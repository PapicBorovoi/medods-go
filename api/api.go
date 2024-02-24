package api

import (
	"net/http"
	"encoding/json"
)

type GetJWTPairParams struct {
	Id string
}

type GetJWTPairResponse struct {
	AccessToken string
	RefreshToken string
}

type Error struct {
	Code int
	Message string
}

func writeError(w http.ResponseWriter, msg string, code int) {
	resp := Error{
		Code: code,
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var (
	RequestErrorHandler = func (w http.ResponseWriter, err error)  {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func (w http.ResponseWriter) {
		writeError(w, "Internal error", http.StatusInternalServerError)
	}
)