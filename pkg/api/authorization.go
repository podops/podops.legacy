package api

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/auth"
	"github.com/txsvc/service/pkg/svc"
)

type (
	authorizationRequest struct {
		Secret     string `json:"secret" binding:"required"`
		Realm      string `json:"realm" binding:"required"`
		ClientID   string `json:"client_id" binding:"required"`
		ClientType string `json:"client_type" binding:"required"` // user,app,bot
		UserID     string `json:"user_id" binding:"required"`
		Scope      string `json:"scope" binding:"required"`
		Duration   int64  `json:"duration" binding:"required"`
	}

	authorizationResponse struct {
		Realm    string `json:"realm" binding:"required"`
		ClientID string `json:"client_id" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}
)

// CreateJWTAuthorizationEndpoint creates an JWT authorization
func CreateJWTAuthorizationEndpoint(c *gin.Context) {
	var ar authorizationRequest

	// this endpoint is secured by a master token i.e. a shared secret between the service and the client
	bearer := svc.GetBearerToken(c)
	if bearer != env.GetString("MASTER_KEY", "") {
		svc.StandardNotAuthorizedResponse(c)
		return
	}

	err := c.BindJSON(&ar)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	token, err := auth.CreateJWTToken(ar.Secret, ar.Realm, ar.ClientID, ar.UserID, ar.Scope, ar.Duration)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	now := util.Timestamp()
	a := auth.Authorization{
		ClientID:  ar.ClientID,
		Name:      ar.Realm,
		Token:     token,
		TokenType: ar.ClientType,
		UserID:    ar.UserID,
		Scope:     ar.Scope,
		Expires:   now + (ar.Duration * 86400),
		AuthType:  "jwt",
		Created:   now,
		Updated:   now,
	}
	err = auth.CreateAuthorization(appengine.NewContext(c.Request), &a)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	resp := authorizationResponse{
		Realm:    ar.Realm,
		ClientID: ar.ClientID,
		Token:    token,
	}

	svc.StandardJSONResponse(c, &resp, nil)
}
