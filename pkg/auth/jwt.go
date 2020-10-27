package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/service/pkg/jwt"
	"github.com/txsvc/service/pkg/svc"
)

type (

	// Client represents the claim of the client calling the API
	Client struct {
		ClientID string `json:"client_id"`
		UserID   string `json:"user_id"`
		Scope    string `json:"scope"`
	}
)

const (
	identityKey = "client_id"
)

// GetSecureJWTMiddleware instantiates a JWT middleware and all the necessary handlers
func GetSecureJWTMiddleware(realm, secretKey string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm: realm,
		Key:   []byte(secretKey),
		//Timeout:         timeout,
		//MaxRefresh:      maxRefresh,
		IdentityKey:     identityKey,
		PayloadFunc:     PayloadMappingHandler,
		IdentityHandler: IdentityHandler,
		Authenticator:   nil, // none provided as we do not have a 'login' function for API clients
		Authorizator:    ScopeAuthorizationHandler,
		//Unauthorized:    Unauthorized,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

// PayloadMappingHandler extracts the client_id, user_id and scope of the request
func PayloadMappingHandler(data interface{}) jwt.MapClaims {
	if v, ok := data.(*Client); ok {
		return jwt.MapClaims{
			"client_id": v.ClientID,
			"user_id":   v.UserID,
			"scope":     v.Scope,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler returns the Client structure
func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	if claims[identityKey] == "" {
		// FIXME: see Issue #170, check if identityKey exists in claims
		return nil
	}
	return &Client{
		ClientID: claims[identityKey].(string),
		UserID:   claims["user_id"].(string),
		Scope:    claims["scope"].(string),
	}
}

// ScopeAuthorizationHandler checks for required scopes
func ScopeAuthorizationHandler(data interface{}, c *gin.Context) bool {
	// FIXME: this is a very simple and naive implementation !
	if v, ok := data.(*Client); ok {
		rr := svc.GetRequiredScopes(c)
		return strings.Contains(v.Scope, rr)
	}
	return false
}
