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

func cleaner() {
	account, err := LookupAccount(context.TODO(), realm, userID)
	if err == nil && account != nil {
		a := authorizationKey(account.Realm, account.ClientID)
		platform.DataStore().Delete(context.TODO(), a)

		k := accountKey(realm, userID)
		platform.DataStore().Delete(context.TODO(), k)
	}
}

func getAccount(t *testing.T) *Account {
	account, err := LookupAccount(context.TODO(), realm, userID)
	if assert.NoError(t, err) {
		if assert.NotNil(t, account) {
			return account
		}
	}
	t.FailNow()
	return nil
}

// Scenario 1:
// - no account
// - request login, create the account
// - confirm the account and send auth token
// - exchange auth token for permanent token
// - delete account & authorization
func TestLoginScenario1(t *testing.T) {
	fmt.Println("Scenario 1")

	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t) // new account, request login, create the account
	loginStep2(t) // confirm the new account, send auth token
	loginStep3(t) // exchange auth token for a permanent token

	fmt.Println("Scenario 1. done")
}

// Scenario 2:
// - no account
// - request login, create the account
// - request login again, expect to reuse existing account
// - continue as Scenario 1
func TestLoginScenario2(t *testing.T) {
	fmt.Println("Scenario 2")

	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t) // new account, request login, create the account
	account1 := getAccount(t)

	loginStep1(t) // existing account, request login again, create the account
	account2 := getAccount(t)

	// requires a new token
	assert.NotEqual(t, account1.Ext1, account2.Ext1)

	loginStep2(t) // confirm the new account, send auth token
	loginStep3(t) // exchange auth token for a permanent token

	fmt.Println("Scenario 2. done")
}

func loginStep1(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.LoginRequestRoute, strings.NewReader(createAuthRequestJSON(realm, userID, "", "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := LoginEndpoint(c)

	if assert.NoError(t, err) {
		account := getAccount(t)
		assert.NotEqual(t, int64(0), account.Ext1)
		assert.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	}
}

func loginStep2(t *testing.T) {

	account := getAccount(t)

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
		account := getAccount(t)
		assert.NotEqual(t, int64(0), account.Confirmed)
		assert.Equal(t, AccountLoggedOut, account.Status)
		assert.NotEqual(t, int64(0), account.Ext2)
		assert.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
	}
}

func loginStep3(t *testing.T) {

	account := getAccount(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(createAuthRequestJSON(realm, userID, account.ClientID, account.Ext2)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := GetAuthorizationEndpoint(c)

	if assert.NoError(t, err) {
		account := getAccount(t)
		assert.Equal(t, AccountActive, account.Status)
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	}
}
