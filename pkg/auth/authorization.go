package auth

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/platform/pkg/platform"
)

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Name      string `json:"name"`                         // name of the domain, realm, tennant etc
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // user,app,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUSerID in Slack
		Scope     string `json:"scope"`                         // a comma separated list of scopes, see below
		Expires   int64  `json:"expires"`                       // 0 = never
		// internal
		AuthType string `json:"-"` // currently: jwt, slack
		Created  int64  `json:"-"`
		Updated  int64  `json:"-"`
	}
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"

	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant salack
	AuthTypeSlack = "slack"
)

// GetToken returns the oauth token of the workspace integration
func GetToken(ctx context.Context, clientID, authType string) (string, error) {
	// ENV always overrides anything else ...
	token := env.GetString(strings.ToUpper(fmt.Sprintf("%s_AUTH_TOKEN", authType)), "")
	if token != "" {
		return token, nil
	}

	// check the in-memory cache
	key := namedKey(clientID, authType)
	token, _ = platform.GetKV(ctx, key)
	if token != "" {
		return token, nil
	}

	auth, err := GetAuthorization(ctx, clientID, authType)
	if err != nil {
		return "", err
	}

	// add the token to the cache
	platform.SetKV(ctx, key, auth.Token, 1800)

	return auth.Token, nil
}

// GetAuthorization looks for an authorization
func GetAuthorization(ctx context.Context, clientID, authType string) (*Authorization, error) {
	var auth Authorization
	k := authorizationKey(clientID, authType)

	if err := platform.DataStore().Get(ctx, k, &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}

// CreateAuthorization creates all data needed for the OAuth fu
func CreateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.ClientID, auth.AuthType)

	// remove the entry from the cache if it is already there ...
	platform.InvalidateKV(ctx, namedKey(auth.ClientID, auth.AuthType))

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := platform.DataStore().Put(ctx, k, auth)
	return err
}

// authorizationKey creates a datastore key for a workspace authorization based on the team_id.
func authorizationKey(clientID, authType string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, namedKey(clientID, authType), nil)
}

func namedKey(clientID, authType string) string {
	return authType + "." + clientID
}
