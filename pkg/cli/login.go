package cli

import (
	"fmt"
	"net/http"

	"github.com/podops/podops/pkg/auth"
	"github.com/urfave/cli/v2"
)

const (
	loginEndpoint  = "/_a/login"
	logoutEndpoint = "/_a/logout"
	authEndpoint   = "/_a/auth"
)

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
			fmt.Println("Login error, logout first.") // FIXME better text
			return nil
		default:
			return fmt.Errorf("api error %d", status)
		}
	} else {
		fmt.Println("error: missing email")
	}

	return nil
}

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {

	if c.Args().Len() != 2 {
		fmt.Println("error: missing parameters")
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
		if err := updateNetrc(response.UserID, response.ClientID, response.Token); err != nil {
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
		return fmt.Errorf("api error %d", status)
	}

}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {

	m := loadNetrc().FindMachine(machineEntry)
	if m == nil {
		return fmt.Errorf("cli error")
	}
	request := auth.AuthorizationRequest{
		Realm:    client.Realm(),
		ClientID: m.Account,
		UserID:   m.Login,
	}

	status, err := post(client.APIEndpoint()+logoutEndpoint, &request, nil)
	if err != nil {
		return err
	}

	if status == http.StatusNoContent {
		fmt.Println("Logout successful.")
	} else {
		return fmt.Errorf("api error %d", status)
	}

	return nil
}
