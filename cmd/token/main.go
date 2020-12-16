package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/txsvc/service/pkg/auth"
)

var (
	secret   string
	clientID string
	userID   string
	scope    string
	realm    string
	duration int64
)

func main() {
	// parse the command line for all the data that goes into the token
	flag.StringVar(&secret, "secret", "", "Secret used to sign the token")
	flag.StringVar(&clientID, "client", "", "The client ID the token belongs to")
	flag.StringVar(&userID, "user", "", "The user ID the token belongs to")
	flag.StringVar(&scope, "scope", "", "The scope of the request token")
	flag.StringVar(&realm, "realm", "", "The realm the token is valid for")
	flag.Int64Var(&duration, "duration", 30, "Validity of the token in days")
	flag.Parse()

	token, err := auth.CreateJWTToken(secret, realm, clientID, userID, scope, duration)
	if err != nil {
		log.Fatal("Error:" + err.Error())
		os.Exit(1)
	}

	fmt.Printf("token='%s'\n\n", token)
}
