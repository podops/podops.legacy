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
			fmt.Println("New account created. Check your inbox and confirm the email address.")
			return nil
		case http.StatusNoContent:
			fmt.Println("Login verificaction sent. Check your inbox.")
			return nil
		case http.StatusForbidden:
			fmt.Println("Already logged-in, use 'po logout' first.")
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
		fmt.Println(errordef.MsgArgumentCountMismatch, 2, c.Args().Len())
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
			fmt.Println("Error updating config.")
			return nil
		}
		fmt.Println("Sucessfully authenticated.")
		return nil
	case http.StatusUnauthorized:
		fmt.Println("Token is expired")
		return nil
	case http.StatusNotFound:
		fmt.Println("Invalid token")
		return nil
	default:
		return fmt.Errorf(errordef.MsgStatus, status)
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
		fmt.Println("Logout successful.")
	} else {
		return fmt.Errorf(errordef.MsgStatus, status)
	}

	return nil
}
