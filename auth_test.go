package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestHandleLogin(t *testing.T) {
	db, err := initDB(true)
	if err != nil {
		t.Errorf("Error intilizing DB: %s", err)
	}

	service := Service{DB: db}

	service.DB.CreateUser("user@example.com", "user", "password")
	service.DB.CreateUser("admin@example.com", "admin", "Secur3")

	reqs := []loginReq{
		{Email: "user@example.com", PlainPassword: "password"},
		{Email: "admin@example.com", PlainPassword: "wrongPassword"},
		{Email: "not-exist@example.com", PlainPassword: "nothing"},
	}

	expectedStatus := []int{
		http.StatusOK,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
	}

	for index, req := range reqs {
		buf := bytes.NewBuffer([]byte{})
		json.NewEncoder(buf).Encode(req)

		r := httptest.NewRequest(http.MethodPost, "/login", buf)
		w := httptest.NewRecorder()

		service.HandleLogin(w, r)
		res := w.Result()
		defer res.Body.Close()

		//TODO check body result?
		// data, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	t.Errorf("expected error to be nil got %v", err)
		// }

		// fmt.Printf("%s", data)

		if res.StatusCode != expectedStatus[index] {
			t.Errorf("expected status %d found %d", expectedStatus[index], res.StatusCode)
		}
	}
}
