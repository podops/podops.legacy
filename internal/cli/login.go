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
		return fmt.Errorf(messagedef.MsgArgumentMissing, "EMAIL")
	}

	if c.Args().Len() == 1 {

		// login EMAIL
		email := c.Args().Get(0)

		if !podops.ValidEmail(email) {
			return fmt.Errorf(messagedef.MsgLoginInvalidEmail, email)
		}

		loginRequest := authentication.AuthorizationRequest{
			Realm:  client.Realm(),
			UserID: email,
		}

		status, err := post(client.APIEndpoint()+loginEndpoint, &loginRequest, nil)
		if err != nil {
			return err
		}

		switch status {
		case http.StatusCreated:
			printMsg(messagedef.MsgLoginNewAccount)
			return nil
		case http.StatusNoContent:
			printMsg(messagedef.MsgLoginVerification)
			return nil
		case http.StatusForbidden:
			return fmt.Errorf(messagedef.MsgLoginError)
		default:
			return fmt.Errorf(messagedef.MsgServerError, status)
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
			return err
		}

		switch status {
		case http.StatusOK:
			if err := storeLogin(response.UserID, response.Token); err != nil {
				printMsg(messagedef.MsgErrorUpdatingConfig)
				return nil
			}
			fmt.Println(messagedef.MsgLoginSuccess)
			return nil
		case http.StatusUnauthorized:
			return fmt.Errorf(messagedef.MsgAuthenticationTokenExpired)
		case http.StatusNotFound:
			return fmt.Errorf(messagedef.MsgAuthenticationTokenInvalid)
		default:
			return fmt.Errorf(messagedef.MsgServerError, status)
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
		return err
	}

	if status == http.StatusNoContent {
		clearLogin()
		printMsg(messagedef.MsgLogoutSuccess)
	} else {
		return fmt.Errorf(messagedef.MsgServerError, status)
	}

	return nil
}
