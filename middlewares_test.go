package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StreatCodes/bernard/database"
	_ "github.com/mattn/go-sqlite3"
)

func generateSession(service *Service) (string, error) {
	userId, err := service.DB.CreateUser("user@example.com", "user", "password")
	if err != nil {
		return "", err
	}
	bToken, err := service.DB.CreateSession(userId)
	return fmt.Sprintf("%x", bToken), err
}

func TestUserAuthMiddleWareAllow(t *testing.T) {
	db, err := initDB(true)
	if err != nil {
		t.Errorf("Error intilizing DB: %s", err)
	}

	service := Service{DB: db}

	token, err := generateSession(&service)
	if err != nil {
		t.Errorf("Error generating token: %s", err)
	}

	//Custom test handler to check it's called when it's meant to be with the appropriate context
	testHandlerWasCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ContextUserKey).(database.User)

		if user.Email != "user@example.com" {
			t.Errorf("expected context user email to be 'user@example.com' found %s", user.Email)
		}
		testHandlerWasCalled = true
	})

	//Call middleware with token
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/host", nil)
	r.Header.Add("Authorization", token)

	handlerToTest := service.userAuthMiddleWare(testHandler)
	handlerToTest.ServeHTTP(w, r)
	res := w.Result()

	//The handler should run with the valid token
	if !testHandlerWasCalled {
		t.Errorf("expected test handler to be called")
	}
	if res.StatusCode != 200 {
		t.Errorf("expected status code to be 200 found %d", res.StatusCode)
	}
}

func TestUserAuthMiddleWareDeny(t *testing.T) {
	db, err := initDB(true)
	if err != nil {
		t.Errorf("Error intilizing DB: %s", err)
	}

	service := Service{DB: db}

	_, err = generateSession(&service)
	if err != nil {
		t.Errorf("Error generating token: %s", err)
	}

	//Custom test handler to check it's called when it's meant to be with the appropriate context
	testHandlerWasCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandlerWasCalled = true
	})

	invalidTokens := []string{
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"8fe7fc984ef19134142f6bfb2c1c7236a7db9a37252468156024165372a82b3491660c5a8d9edaea11bd1e37d797e2bd9ca1c6c7ef475e30f9fb58a6",
		"8fe7fc984ef19134142f6bfb2c1c7236a7db9a37252468156024165372a82b3491660c5a8d9edaea11bd1e37d797e2bd9ca1c6c7ef475e30f9fb58a6ef475e30f9fb58a6",
		"8fe7fc984ef19134142f6bfb2c1c7236a7db9a37252468156024165372a82b3491660c5a8d9edaea11bd1e37d797e2bd9ca1c6c7ef475e3xxxxxxxxxxxxxxxxx",
		"8fe7fc984ef19134142f6bfb2c1c7236a7db9a37252468156024165372a82b3491660c5a8d9edaea11bd1e37d797e2bd9ca1c6c7ef47invalid-characters",
		"",
	}

	for _, invalidToken := range invalidTokens {
		//Call middleware with token
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/host", nil)
		r.Header.Add("Authorization", invalidToken)

		handlerToTest := service.userAuthMiddleWare(testHandler)
		handlerToTest.ServeHTTP(w, r)
		res := w.Result()

		//The handler should run with the valid token
		if testHandlerWasCalled {
			t.Errorf("Request handler shouldn't have been called with the invalid token %s", invalidToken)
		}
		if res.StatusCode < 400 {
			t.Errorf("expected status code to be >= 400 found %d", res.StatusCode)
		}
	}

	// Request with no authorization header should return 401
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/host", nil)

	handlerToTest := service.userAuthMiddleWare(testHandler)
	handlerToTest.ServeHTTP(w, r)
	res := w.Result()

	//The handler should run with the valid token
	if testHandlerWasCalled {
		t.Errorf("Request handler shouldn't have been called without an authorization token")
	}
	if res.StatusCode != 401 {
		t.Errorf("expected status code to be 401 found %d", res.StatusCode)
	}

}
