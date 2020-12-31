package main

//
// FIXME
// This is only a placeholder/example. It creates a token but does not create the AUTH/ACCOUNT structs in the backend.
//

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/podops/podops"
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

	cl, err := podops.NewClient("")
	if err != nil {
		log.Fatal("Error:" + err.Error())
		os.Exit(1)
	}
	token, err := cl.CreateToken(secret, realm, clientID, userID, scope, duration)
	if err != nil {
		log.Fatal("Error:" + err.Error())
		os.Exit(1)
	}

	fmt.Printf("token='%s'\n\n", token)
}
