package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/txsvc/platform/v2/authentication"
	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/id"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

// This utility creates/updates an API user
func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("Not enough parameters.")
	}

	realm := args[0]
	userID := args[1]
	expires := 0
	if len(args) >= 3 {
		ex, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("Invalid expiration.")
		}
		if ex >= 0 {
			expires = ex
		}
	}

	ctx := context.Background()
	now := timestamp.Now()

	// create or update the account
	acc, err := account.FindAccountByUserID(ctx, realm, userID)
	if err != nil {
		log.Fatal(err)
	}

	if acc == nil {
		uid, _ := id.ShortUUID()
		acc = &account.Account{
			Realm:     realm,
			UserID:    userID,
			ClientID:  uid,
			Status:    account.AccountActive,
			Confirmed: now,
			Created:   now,
			Updated:   now,
		}
	}
	if expires != 0 {
		acc.Expires = now + (int64(expires) * 86400)
	} else {
		acc.Expires = 0
	}
	acc.Updated = now

	if err := account.UpdateAccount(ctx, acc); err != nil {
		log.Fatal(err)
	}

	// create or update the authorization

	ath, err := authentication.LookupAuthorization(ctx, acc.Realm, acc.ClientID)
	if err != nil {
		log.Fatal(err)
	}
	if ath == nil {
		req := authentication.AuthorizationRequest{
			Realm:    realm,
			UserID:   userID,
			ClientID: acc.ClientID,
			Scope:    authentication.ScopeAPIAdmin,
		}
		ath = authentication.NewAuthorization(&req, expires)
	}
	ath.Token = authentication.CreateSimpleToken()
	ath.TokenType = "api"
	if expires != 0 {
		ath.Expires = now + (int64(expires) * 86400)
	} else {
		ath.Expires = 0
	}
	ath.Updated = now

	err = authentication.UpdateAuthorization(ctx, ath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New authorization created. Token='%s'\n", ath.Token)
}
