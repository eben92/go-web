package server

import "net/http"

type ApiError struct {
	Error string `json:"error"`
}

type ApiResp struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error
