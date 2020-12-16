package auth

import (
	"time"
)

// CreateJWTToken creates a token that can be used for JWT authentication / authorization
func CreateJWTToken(secret, realm, clientID, userID, scope string, duration int64) (string, error) {

	a, err := GetSecureJWTMiddleware(realm, secret)
	if err != nil {
		return "", err
	}

	a.Timeout = (time.Duration)(duration*24) * time.Hour

	// the claim
	client := Client{
		ClientID: clientID,
		UserID:   userID,
		Scope:    scope,
	}

	token, _, err := a.TokenGenerator(&client)
	if err != nil {
		return "", err
	}

	return token, nil
}
