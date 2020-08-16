package handlers

import (
	"encoding/json"
	"net/http"
)

type Resp struct {
	Info string `json:"info"`
}

var Err error

func Header(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Vary", "Origin")
}

func Response(w http.ResponseWriter, r *http.Request, status int, info string) {
	var response Resp
	
	response.Info = info

	Header(w)
	
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
