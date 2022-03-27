package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

func (service *Service) HandleGetHosts(w http.ResponseWriter, r *http.Request) {
	var err error
	after := r.URL.Query().Get("after")
	limit := r.URL.Query().Get("limit")

	iAfter := int64(0)
	if after != "" {
		iAfter, err = strconv.ParseInt(after, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	iLimit := int64(20)
	if limit != "" {
		iLimit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	hosts, err := service.DB.SelectHosts(iAfter, iLimit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(hosts)
}

func (service *Service) HandleGetHost(w http.ResponseWriter, r *http.Request) {
	hostID := chi.URLParam(r, "id")
	iHostID, err := strconv.ParseInt(hostID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	host, err := service.DB.SelectHostById(iHostID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(host)
}

func (service *Service) HandleGetHostLogs(w http.ResponseWriter, r *http.Request) {
	hostID := chi.URLParam(r, "id")
	iHostID, err := strconv.ParseInt(hostID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO implement start / end
	host, err := service.DB.SelectHostLogsBetween(iHostID, time.Now().Add(-7*24*time.Hour), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(host)

}
