package database

import "crypto/rand"

const TokenLength = 64

func generateToken() ([]byte, error) {
	token := make([]byte, TokenLength)

	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
