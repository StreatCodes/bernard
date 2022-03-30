package main

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/StreatCodes/bernard/database"
)

type ContextKey string

const ContextUserKey ContextKey = "user"
const ContextHostKey ContextKey = "host"

func (service *Service) userAuthMiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		if encodedToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(encodedToken) != database.TokenLength*2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := hex.DecodeString(encodedToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := service.DB.SelectUserByToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (service *Service) hostAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		if len(encodedToken) != database.TokenLength*2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := hex.DecodeString(encodedToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		host, err := service.DB.SelectHostByKey(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextHostKey, host)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
