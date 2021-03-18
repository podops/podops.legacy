package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fupas/platform/pkg/platform"
	"github.com/labstack/echo/v4"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/api"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint = "http://localhost:8080"

	realm  = "podops"
	userID = "me@podops.dev"
)

func createAuthRequestJSON(real, user, client, token string) string {
	return fmt.Sprintf(`{"realm":"%s","user_id":"%s","client_id":"%s","token":"%s"}`, realm, user, client, token)
}
func TestLoginSequence(t *testing.T) {

	apiv1.DefaultAPIEndpoint = endpoint

	t.Cleanup(func() {
		account, err := LookupAccount(context.TODO(), realm, userID)
		if err == nil && account != nil {
			a := authorizationKey(account.Realm, account.ClientID)
			platform.DataStore().Delete(context.TODO(), a)

			k := accountKey(realm, userID)
			platform.DataStore().Delete(context.TODO(), k)
		}

	})

	loginStep1(t) // new account, request login
	loginStep2(t) // confirm the new account
	loginStep3(t) // exchange temporary token for a permanent token
}

func loginStep1(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.LoginRequestRoute, strings.NewReader(createAuthRequestJSON(realm, userID, "", "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := LoginEndpoint(c)

	if assert.NoError(t, err) {
		account, err := LookupAccount(context.TODO(), realm, userID)
		if assert.NoError(t, err) && assert.NotNil(t, account) {
			assert.NotEqual(t, int64(0), account.Ext1)
			assert.Equal(t, http.StatusCreated, rec.Result().StatusCode)
		}
	}
}

func loginStep2(t *testing.T) {

	account, err := LookupAccount(context.TODO(), realm, userID)
	if assert.NoError(t, err) && assert.NotNil(t, account) {

		url := fmt.Sprintf("/login/%s", account.Ext1)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		r := e.Router()
		r.Add(http.MethodGet, api.LoginConfirmationRoute, LoginConfirmationEndpoint)

		c := e.NewContext(req, rec)
		r.Find(http.MethodGet, url, c)

		handler := c.Handler()
		err := handler(c)

		if assert.NoError(t, err) {
			account, err = LookupAccount(context.TODO(), realm, userID)
			if assert.NoError(t, err) && assert.NotNil(t, account) {
				assert.NotEqual(t, int64(0), account.Confirmed)
				assert.Equal(t, AccountLoggedOut, account.Status)
				assert.NotEqual(t, int64(0), account.Ext2)
				assert.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
			}
		}
	}
}

func loginStep3(t *testing.T) {

	account, err := LookupAccount(context.TODO(), realm, userID)
	if assert.NoError(t, err) && assert.NotNil(t, account) {

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(createAuthRequestJSON(realm, userID, account.ClientID, account.Ext2)))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := GetAuthorizationEndpoint(c)

		if assert.NoError(t, err) {
			account, err := LookupAccount(context.TODO(), realm, userID)
			if assert.NoError(t, err) && assert.NotNil(t, account) {
				assert.Equal(t, AccountActive, account.Status)
				assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
			}
		}
	}
}
