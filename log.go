package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/StreatCodes/bernard/database"
)

type newLogReq struct {
	Level   int64  `json:"level"`
	Service string `json:"service"`
	Content string `json:"content"`
}

func (service *Service) HandleNewLog(w http.ResponseWriter, r *http.Request) {
	host := r.Context().Value(ContextHostKey).(database.Host)

	req := newLogReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log := database.Log{
		Time:    time.Now(),
		HostID:  host.ID,
		Level:   req.Level,
		Service: req.Service,
		Content: req.Content,
	}
	err = service.DB.InsertLog(log)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
