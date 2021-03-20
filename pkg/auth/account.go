package auth

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
)

const (
	// DatastoreAccounts collection ACCOUNTS
	DatastoreAccounts string = "ACCOUNTS"

	// AccountActive indicates a confirmed account with a valid login
	AccountActive = 1
	// AccountLoggedOut indicates a confirmed account without a valid login
	AccountLoggedOut = 0
	// AccountDeactivated indicates an account that has been deactivated due to
	// e.g. account deletion or UserID swap
	AccountDeactivated = -1
	// AccountBlocked signals an issue with the account that needs intervention
	AccountBlocked = -2
	// AccountUnconfirmed well guess what?
	AccountUnconfirmed = -3
)

type (

	// Account represents an account for a user or client (e.g. API, bot)
	Account struct {
		Realm    string `json:"realm"`     // KEY
		UserID   string `json:"user_id"`   // KEY external id for the entity e.g. email for a user
		ClientID string `json:"client_id"` // a unique id within [realm,user_id]
		// status and other metadata
		Status int `json:"status"` // default == AccountUnconfirmed
		// login auditing
		LastLogin  int64  `json:"-"`
		LoginCount int    `json:"-"`
		LoginFrom  string `json:"-"`
		// internal
		Ext1      string `json:"-"` // universal field, used as needed. e.g to confirm the account and then to request the real token
		Ext2      string `json:"-"`
		Expires   int64  `json:"-"` // 0 == never
		Confirmed int64  `json:"-"`
		Created   int64  `json:"-"`
		Updated   int64  `json:"-"`
	}
)

// LookupAccount retrieves an account within a given realm
func LookupAccount(ctx context.Context, realm, userID string) (*Account, error) {
	var account Account
	k := accountKey(realm, userID)

	err := platform.DataStore().Get(ctx, k, &account)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// FindAccountByToken retrieves an account bases on either the temporary token or the auth token
func FindAccountByToken(ctx context.Context, token string) (*Account, error) {
	var accounts []*Account
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreAccounts).Filter("Ext1 =", token), &accounts); err != nil {
		return nil, err
	}
	if accounts == nil {
		return nil, nil
	}
	return accounts[0], nil
}

// CreateAccount creates an new account within a given realm
func CreateAccount(ctx context.Context, realm, userID string) (*Account, error) {
	now := util.Timestamp()
	token, _ := util.ShortUUID()
	uid, _ := util.ShortUUID() // FIXME verify that uid is unique

	account := &Account{
		Realm:     realm,
		UserID:    userID,
		ClientID:  uid,
		Status:    AccountUnconfirmed,
		Ext1:      token,
		Expires:   util.IncT(util.Timestamp(), DefaultAuthenticationExpiration),
		Confirmed: 0,
		Created:   now,
		Updated:   now,
	}

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func UpdateAccount(ctx context.Context, account *Account) error {
	k := accountKey(account.Realm, account.UserID)
	account.Updated = util.Timestamp()

	if _, err := platform.DataStore().Put(ctx, k, account); err != nil {
		return err
	}
	return nil
}

func accountKey(realm, user string) *datastore.Key {
	return datastore.NameKey(DatastoreAccounts, namedKey(realm, user), nil)
}