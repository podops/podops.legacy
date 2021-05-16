package cli

import (
	"fmt"
	"net/http"

	"github.com/txsvc/platform/v2/pkg/authentication"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/messagedef"

	"github.com/urfave/cli/v2"
)

const (
	loginEndpoint  = "/_a/login"
	logoutEndpoint = "/logout"
	authEndpoint   = "/_a/auth"
)

// FIXME replace all the messages with consts

// LoginCommand logs into the service
func LoginCommand(c *cli.Context) error {

	if c.Args().Len() == 0 {
		printError(c, fmt.Errorf(messagedef.MsgArgumentMissing, "EMAIL"))
		return nil
	}

	if c.Args().Len() == 1 {

		// login EMAIL
		email := c.Args().Get(0)

		if !podops.ValidEmail(email) {
			printError(c, fmt.Errorf(messagedef.MsgLoginInvalidEmail, email))
			return nil
		}

		loginRequest := authentication.AuthorizationRequest{
			Realm:  client.Realm(),
			UserID: email,
		}

		status, err := post(client.APIEndpoint()+loginEndpoint, &loginRequest, nil)
		if err != nil {
			printError(c, err)
			return nil
		}

		switch status {
		case http.StatusCreated:
			printMsg(messagedef.MsgLoginNewAccount)
			return nil
		case http.StatusNoContent:
			printMsg(messagedef.MsgLoginVerification)
			return nil
		case http.StatusForbidden:
			printError(c, fmt.Errorf(messagedef.MsgLoginError))
			return nil
		default:
			printError(c, fmt.Errorf(messagedef.MsgServerError, status))
			return nil
		}
	} else if c.Args().Len() == 2 {

		// login EMAIL TOKEN

		authRequest := authentication.AuthorizationRequest{
			Realm:  client.Realm(),
			UserID: c.Args().Get(0),
			Token:  c.Args().Get(1),
		}
		response := authentication.AuthorizationRequest{}

		status, err := post(client.APIEndpoint()+authEndpoint, &authRequest, &response)
		if err != nil {
			printError(c, err)
			return nil
		}

		switch status {
		case http.StatusOK:
			if err := storeLogin(response.UserID, response.Token); err != nil {
				printError(c, fmt.Errorf(messagedef.MsgErrorUpdatingConfig))
				return nil
			}
			fmt.Println(messagedef.MsgLoginSuccess)
			return nil
		case http.StatusUnauthorized:
			printError(c, fmt.Errorf(messagedef.MsgAuthenticationTokenExpired))
			return nil
		case http.StatusNotFound:
			printError(c, fmt.Errorf(messagedef.MsgAuthenticationTokenInvalid))
			return nil
		default:
			printError(c, fmt.Errorf(messagedef.MsgServerError, status))
			return nil
		}
	}

	printMsg(messagedef.MsgTooManyArguments)
	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {

	m := loadNetrc().FindMachine(machineEntry)
	if m == nil {
		printMsg(messagedef.MsgNotLoggedIn)
		return nil
	}
	request := authentication.AuthorizationRequest{
		Realm:  client.Realm(),
		UserID: m.Login,
	}

	status, err := post(client.APIEndpoint()+logoutEndpoint, &request, nil)
	if err != nil {
		printError(c, err)
		return nil
	}

	if status == http.StatusNoContent {
		clearLogin()
		printMsg(messagedef.MsgLogoutSuccess)
	} else {
		printError(c, fmt.Errorf(messagedef.MsgServerError, status))
	}

	return nil
}
