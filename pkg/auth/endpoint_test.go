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

func verifyAccountAndAuth(t *testing.T) bool {
	account, err := LookupAccount(context.TODO(), realm, userID)
	if err == nil && account != nil {
		auth, err := LookupAuthorization(context.TODO(), account.Realm, account.ClientID)
		if err == nil && auth != nil {
			return true
		}
	}
	return false
}

// Scenario 1: new account, login, account confirmation, token swap
func TestLoginScenario1(t *testing.T) {
	fmt.Println("scenario 1")

	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)
	cleaner()

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	loginStep2(t, account.Ext1, http.StatusNoContent) // confirm the new account, send auth token

	account = getAccount(t)
	loginStep3(t, account.Ext2, http.StatusOK) // exchange auth token for a permanent token

	assert.True(t, verifyAccountAndAuth(t))
}

// Scenario 2: new account, login, duplicate login request
func TestLoginScenario2(t *testing.T) {
	fmt.Println("scenario 2")

	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account
	account1 := getAccount(t)

	loginStep1(t, http.StatusCreated) // existing account, request login again, create the account
	account2 := getAccount(t)

	// requires a new token
	assert.NotEqual(t, account1.Ext1, account2.Ext1)

	loginStep2(t, account2.Ext1, http.StatusNoContent) // confirm the new account, send auth token

	account3 := getAccount(t)
	loginStep3(t, account3.Ext2, http.StatusOK) // exchange auth token for a permanent token

	assert.True(t, verifyAccountAndAuth(t))
}

// Scenario 3: new account, login, duplicate account confirmation
func TestLoginScenario3(t *testing.T) {
	fmt.Println("scenario 3")

	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	token := account.Ext1

	loginStep2(t, token, http.StatusNoContent) // confirm the new account, send auth token
	loginStep2(t, token, http.StatusNotFound)  // confirm again

	account = getAccount(t)
	loginStep3(t, account.Ext2, http.StatusOK) // exchange auth token for a permanent token

	assert.True(t, verifyAccountAndAuth(t))
}

// Scenario 4: new account, login, account confirmation, duplicate token swap
func TestLoginScenario4(t *testing.T) {
	fmt.Println("scenario 4")

	apiv1.DefaultAPIEndpoint = endpoint
	//t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	loginStep2(t, account.Ext1, http.StatusNoContent) // confirm the new account, send auth token

	account = getAccount(t)
	token := account.Ext2

	loginStep3(t, token, http.StatusOK) // exchange auth token for a permanent token
	loginStep3(t, token, http.StatusUnauthorized)

	assert.True(t, verifyAccountAndAuth(t))
}

func loginStep1(t *testing.T, status int) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.LoginRequestRoute, strings.NewReader(createAuthRequestJSON(realm, userID, "", "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := LoginEndpoint(c)

	if assert.NoError(t, err) {
		account := getAccount(t)
		assert.NotEqual(t, int64(0), account.Ext1)
		assert.Equal(t, status, rec.Result().StatusCode)
	}
}

func loginStep2(t *testing.T, token string, status int) {

	url := fmt.Sprintf("/login/%s", token)
	fmt.Println("confirm: " + url)

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
		assert.Equal(t, status, rec.Result().StatusCode)
	}
}

func loginStep3(t *testing.T, token string, status int) {

	account := getAccount(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(createAuthRequestJSON(realm, userID, account.ClientID, token)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := GetAuthorizationEndpoint(c)

	if assert.NoError(t, err) {
		account := getAccount(t)
		assert.Equal(t, AccountActive, account.Status)
		assert.Equal(t, status, rec.Result().StatusCode)
	}
}
