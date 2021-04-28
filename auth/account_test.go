package auth

import (
	"context"
	"testing"

	"github.com/fupas/platform/pkg/platform"
	"github.com/stretchr/testify/assert"
	"github.com/txsvc/spa/pkg/timestamp"
)

const (
	accountTestRealm = "account_test"
	accountTestUser  = "account_test_user"
)

func cleanup() {
	account, _ := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if account != nil {
		k := accountKey(accountTestRealm, account.ClientID)
		platform.DataStore().Delete(context.TODO(), k)
	}
}

func TestFindAccountByUserIDFail(t *testing.T) {
	cleanup()

	account, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.Nil(t, account)
	}
}

func TestCreateAccount(t *testing.T) {
	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)
		assert.Equal(t, int64(0), account.Confirmed)
		assert.Equal(t, AccountUnconfirmed, account.Status)
	}
}

func TestLookupAccount(t *testing.T) {
	account1, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account1)
	}

	account2, err := LookupAccount(context.TODO(), accountTestRealm, account1.ClientID)
	if assert.NoError(t, err) {
		assert.NotNil(t, account2)
		assert.Equal(t, account1.Realm, account2.Realm)
		assert.Equal(t, account1.UserID, account2.UserID)
	}
}

func TestUpdateAccount(t *testing.T) {
	t.Cleanup(cleanup)

	account1, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account1)

		now := timestamp.Now()
		account1.Confirmed = now
		err = UpdateAccount(context.TODO(), account1)

		if assert.NoError(t, err) {
			account2, err := LookupAccount(context.TODO(), account1.Realm, account1.ClientID)
			if assert.NoError(t, err) {
				assert.Equal(t, now, account2.Confirmed)
			}
		}
	}
}
