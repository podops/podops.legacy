package cli

import (
	"fmt"
	"net/http"

	"github.com/txsvc/platform/v2/authentication"

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
	email := c.Args().First()

	if email != "" {
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
			fmt.Println(messagedef.MsgLoginNewAccount)
			return nil
		case http.StatusNoContent:
			fmt.Println(messagedef.MsgLoginVerification)
			return nil
		case http.StatusForbidden:
			fmt.Println(messagedef.MsgLoginError)
			return nil
		default:
			return fmt.Errorf(messagedef.MsgStatus, status)
		}
	} else {
		fmt.Println(messagedef.MsgArgumentMissing, "EMAIL")
	}

	return nil
}

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {

	if c.Args().Len() != 2 {
		printMsg(messagedef.MsgArgumentCountMismatch, 2, c.Args().Len())
		return nil
	}

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
			fmt.Println(messagedef.MsgErrorUpdatingConfig)
			return nil
		}
		fmt.Println(messagedef.MsgAuthenticationSuccess)
		return nil
	case http.StatusUnauthorized:
		fmt.Println(messagedef.MsgAuthenticationTokenExpired)
		return nil
	case http.StatusNotFound:
		fmt.Println(messagedef.MsgAuthenticationTokenInvalid)
		return nil
	default:
		return fmt.Errorf(messagedef.MsgStatus, status)
	}

}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {

	m := loadNetrc().FindMachine(machineEntry)
	if m == nil {
		return fmt.Errorf(messagedef.MsgClientError)
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
		fmt.Println(messagedef.MsgLogoutSuccess)
	} else {
		return fmt.Errorf(messagedef.MsgStatus, status)
	}

	return nil
}
