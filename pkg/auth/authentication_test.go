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
	realm    = "podops"
	userID   = "me@podops.dev"
)

// Scenario 1: new account, login, account confirmation, token swap
func TestLoginScenario1(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)
	cleaner()

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	loginStep2(t, account.Ext1, http.StatusNoContent, true) // confirm the new account, send auth token

	account = getAccount(t)
	loginStep3(t, realm, userID, account.ClientID, account.Ext2, AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)

	auth, _ := LookupAuthorization(context.TODO(), account.Realm, account.ClientID)
	logoutStep(t, realm, userID, account.ClientID, auth.Token, http.StatusNoContent, true)
}

// Scenario 2: new account, login, duplicate login request
func TestLoginScenario2(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account
	account1 := getAccount(t)

	loginStep1(t, http.StatusCreated) // existing account, request login again, create the account
	account2 := getAccount(t)

	// requires a new token
	assert.NotEqual(t, account1.Ext1, account2.Ext1)

	loginStep2(t, account2.Ext1, http.StatusNoContent, true) // confirm the new account, send auth token

	account3 := getAccount(t)
	loginStep3(t, realm, userID, account3.ClientID, account3.Ext2, AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)
}

// Scenario 3: new account, login, duplicate account confirmation
func TestLoginScenario3(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	token := account.Ext1

	loginStep2(t, token, http.StatusNoContent, true)    // confirm the new account, send auth token
	loginStep2(t, token, http.StatusUnauthorized, true) // confirm again

	account = getAccount(t)
	loginStep3(t, realm, userID, account.ClientID, account.Ext2, AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)
}

// Scenario 4: new account, login, account confirmation, duplicate token swap
func TestLoginScenario4(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	loginStep2(t, account.Ext1, http.StatusNoContent, true) // confirm the new account, send auth token

	account = getAccount(t)
	token := account.Ext2

	loginStep3(t, realm, userID, account.ClientID, token, AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	loginStep3(t, realm, userID, account.ClientID, token, AccountActive, http.StatusUnauthorized, true)

	verifyAccountAndAuth(t)
}

// Scenario 5: new account, login, invalid confirmation
func TestLoginScenario5(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	loginStep2(t, "this_is_not_valid", http.StatusUnauthorized, false)
}

// Scenario 6: new account, login, account confirmation, various invalid token swaps
func TestLoginScenario6(t *testing.T) {
	apiv1.DefaultAPIEndpoint = endpoint
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	account := getAccount(t)
	loginStep2(t, account.Ext1, http.StatusNoContent, true) // confirm the new account, send auth token

	account = getAccount(t)
	loginStep3(t, "", "", "", "", AccountLoggedOut, http.StatusBadRequest, false)
	loginStep3(t, "wrong_realm", "wrong_user", account.ClientID, account.Ext2, AccountLoggedOut, http.StatusNotFound, false)
	loginStep3(t, realm, userID, account.ClientID, "wrong_auth_token", AccountLoggedOut, http.StatusUnauthorized, false)
}

// FIXME test account confirmation timeout

// FIXME test auth swap timeout

func loginStep1(t *testing.T, status int) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.LoginRequestRoute, strings.NewReader(createAuthRequestJSON(realm, userID, "", "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := LoginRequestEndpoint(c)

	if assert.NoError(t, err) {
		account := getAccount(t)
		assert.NotEqual(t, int64(0), account.Ext1)
		assert.Equal(t, status, rec.Result().StatusCode)
	}
}

func loginStep2(t *testing.T, token string, status int, validate bool) {

	url := fmt.Sprintf("/login/%s", token)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	r := e.Router()
	r.Add(http.MethodGet, api.LoginConfirmationRoute, LoginConfirmationEndpoint)

	c := e.NewContext(req, rec)
	r.Find(http.MethodGet, url, c)
	err := LoginConfirmationEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			account := getAccount(t)
			assert.NotEqual(t, int64(0), account.Confirmed)
			assert.Equal(t, AccountLoggedOut, account.Status)
			assert.NotEqual(t, int64(0), account.Ext2)
		}
	}
}

func loginStep3(t *testing.T, testRealm, testUser, testClient, testToken string, state, status int, validate bool) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.GetAuthorizationRoute, strings.NewReader(createAuthRequestJSON(testRealm, testUser, "", testToken)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := GetAuthorizationEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			account := getAccount(t)
			assert.Equal(t, state, account.Status)
			assert.Equal(t, "", account.Ext1)
			assert.Equal(t, "", account.Ext2)
		}
	}
}

func logoutStep(t *testing.T, testRealm, testUser, testClient, testToken string, status int, validate bool) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, api.LogoutRequestRoute, strings.NewReader(createAuthRequestJSON(testRealm, testUser, testClient, "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+testToken)
	c := e.NewContext(req, rec)

	err := LogoutRequestEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			account := getAccount(t)
			assert.Equal(t, AccountLoggedOut, account.Status)

		}
	}
}

func createAuthRequestJSON(real, user, client, token string) string {
	return fmt.Sprintf(`{"realm":"%s","user_id":"%s","client_id":"%s","token":"%s"}`, realm, user, client, token)
}

func cleaner() {
	account, err := FindAccountByUserID(context.TODO(), realm, userID)
	if err == nil && account != nil {
		a := authorizationKey(account.Realm, account.ClientID)
		platform.DataStore().Delete(context.TODO(), a)

		k := accountKey(realm, account.ClientID)
		platform.DataStore().Delete(context.TODO(), k)
	}
}

func getAccount(t *testing.T) *Account {
	account, err := FindAccountByUserID(context.TODO(), realm, userID)
	if assert.NoError(t, err) {
		if assert.NotNil(t, account) {
			return account
		}
	}
	t.FailNow()
	return nil
}

func verifyAccountAndAuth(t *testing.T) {
	account, err := FindAccountByUserID(context.TODO(), realm, userID)
	if err == nil && account != nil {
		auth, err := LookupAuthorization(context.TODO(), account.Realm, account.ClientID)
		if err == nil && auth != nil {
			assert.Equal(t, account.ClientID, auth.ClientID)
		}
	}
}
