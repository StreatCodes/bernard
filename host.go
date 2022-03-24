package main

import (
	"encoding/json"
	"net/http"
)

type createHostReq struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Descrpition string `json:"descrpition"`
}

type createHostRes struct {
	ID int64 `json:"id"`
}

func (service *Service) HandleCreateHost(w http.ResponseWriter, r *http.Request) {
	req := createHostReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := service.DB.CreateHost(req.Domain, req.Name, req.Descrpition)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := createHostRes{
		ID: id,
	}

	json.NewEncoder(w).Encode(res)
}
