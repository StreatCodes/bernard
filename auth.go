package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type loginReq struct {
	Email         string `json:"email"`
	PlainPassword string `json:"password"`
}

type loginRes struct {
	Token string `json:"token"`
}

func (service *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	req := loginReq{}
	json.NewDecoder(r.Body).Decode(&req)

	user, err := service.DB.SelectUserByCreds(req.Email, req.PlainPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := service.DB.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	res := loginRes{Token: fmt.Sprintf("%x", token)}

	json.NewEncoder(w).Encode(res)
}
