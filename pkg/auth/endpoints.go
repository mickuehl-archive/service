package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/svc"
)

// CreateJWTAuthorizationEndpoint creates an JWT authorization
func CreateJWTAuthorizationEndpoint(c *gin.Context) {
	var ar AuthorizationRequest

	// this endpoint is secured by a master token i.e. a shared secret between
	// the service and the client, NOT a JWT token !!
	bearer := GetBearerToken(c)
	if bearer != env.GetString("MASTER_KEY", "") {
		svc.StandardNotAuthorizedResponse(c)
		return
	}

	err := c.BindJSON(&ar)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	token, err := CreateJWTToken(ar.Secret, ar.Realm, ar.ClientID, ar.UserID, ar.Scope, ar.Duration)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	now := util.Timestamp()
	a := Authorization{
		ClientID:  ar.ClientID,
		Name:      ar.Realm,
		Token:     token,
		TokenType: ar.ClientType,
		UserID:    ar.UserID,
		Scope:     ar.Scope,
		Expires:   now + (ar.Duration * 86400), // Duration days from now
		AuthType:  "jwt",
		Created:   now,
		Updated:   now,
	}
	err = CreateAuthorization(appengine.NewContext(c.Request), &a)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	resp := AuthorizationResponse{
		Realm:    ar.Realm,
		ClientID: ar.ClientID,
		Token:    token,
	}

	svc.StandardJSONResponse(c, &resp, nil)
}

// ValidateJWTAuthorizationEndpoint verifies that the token is valid and exists in the authorization table
func ValidateJWTAuthorizationEndpoint(c *gin.Context) {
	token := GetBearerToken(c)
	if token == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	auth, err := FindAuthorization(appengine.NewContext(c.Request), token)
	if auth == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	if !auth.IsValid() {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.Status(http.StatusAccepted)
}
