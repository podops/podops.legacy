package auth

import (
	"context"
	"testing"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
	"github.com/stretchr/testify/assert"
)

const (
	accountTestRealm = "account_test"
	accountTestUser  = "account_test_user"
)

func TestLookupFailure(t *testing.T) {
	account, err := LookupAccount(context.TODO(), accountTestRealm, accountTestUser)
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
	account, err := LookupAccount(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)
		assert.Equal(t, accountTestRealm, account.Realm)
		assert.Equal(t, accountTestUser, account.UserID)
	}
}

func TestUpdateAccount(t *testing.T) {
	t.Cleanup(func() {
		k := accountKey(accountTestRealm, accountTestUser)
		platform.DataStore().Delete(context.TODO(), k)
	})

	account1, err := LookupAccount(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account1)

		now := util.Timestamp()
		account1.Confirmed = now
		err = UpdateAccount(context.TODO(), account1)

		if assert.NoError(t, err) {
			account2, err := LookupAccount(context.TODO(), account1.Realm, account1.UserID)
			if assert.NoError(t, err) {
				assert.Equal(t, now, account2.Confirmed)
			}
		}
	}
}
