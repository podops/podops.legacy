package cli

import (
	"fmt"
	"net/http"

	"github.com/podops/podops/auth"
	"github.com/podops/podops/internal/errordef"

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
		loginRequest := auth.AuthorizationRequest{
			Realm:  client.Realm(),
			UserID: email,
		}

		status, err := post(client.APIEndpoint()+loginEndpoint, &loginRequest, nil)
		if err != nil {
			return err
		}

		switch status {
		case http.StatusCreated:
			fmt.Println(errordef.MsgCLINewAccount)
			return nil
		case http.StatusNoContent:
			fmt.Println(errordef.MsgCLILoginVerification)
			return nil
		case http.StatusForbidden:
			fmt.Println(errordef.MsgCLILoginError)
			return nil
		default:
			return fmt.Errorf(errordef.MsgStatus, status)
		}
	} else {
		fmt.Println(errordef.MsgMissingArgument, "EMAIL")
	}

	return nil
}

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {

	if c.Args().Len() != 2 {
		printMsg(errordef.MsgArgumentCountMismatch, 2, c.Args().Len())
		return nil
	}

	authRequest := auth.AuthorizationRequest{
		Realm:  client.Realm(),
		UserID: c.Args().Get(0),
		Token:  c.Args().Get(1),
	}
	response := auth.AuthorizationRequest{}

	status, err := post(client.APIEndpoint()+authEndpoint, &authRequest, &response)
	if err != nil {
		return err
	}

	switch status {
	case http.StatusOK:
		if err := storeLogin(response.UserID, response.Token); err != nil {
			fmt.Println(errordef.MsgCLIErrorUpdateConfig)
			return nil
		}
		fmt.Println(errordef.MsgCLIAuthSuccess)
		return nil
	case http.StatusUnauthorized:
		fmt.Println(errordef.MsgCLITokenExpired)
		return nil
	case http.StatusNotFound:
		fmt.Println(errordef.MsgCLITokenInvalid)
		return nil
	default:
		return fmt.Errorf(errordef.MsgCLIStatus, status)
	}

}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {

	m := loadNetrc().FindMachine(machineEntry)
	if m == nil {
		return fmt.Errorf(errordef.MsgCLIError)
	}
	request := auth.AuthorizationRequest{
		Realm:  client.Realm(),
		UserID: m.Login,
	}

	status, err := post(client.APIEndpoint()+logoutEndpoint, &request, nil)
	if err != nil {
		return err
	}

	if status == http.StatusNoContent {
		clearLogin()
		fmt.Println(errordef.MsgCLILogout)
	} else {
		return fmt.Errorf(errordef.MsgCLIStatus, status)
	}

	return nil
}
