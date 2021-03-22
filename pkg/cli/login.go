package cli

import (
	"context"
	"fmt"
	"net/http"

	"github.com/podops/podops"
	"github.com/podops/podops/pkg/auth"
	"github.com/urfave/cli/v2"
)

const (
	loginEndpoint = "/_a/login"
	authEndpoint  = "/_a/auth"
)

// LoginCommand logs into the service
func LoginCommand(c *cli.Context) error {
	email := c.Args().First()

	if email != "" {
		cl, err := podops.NewClient(context.TODO(), "")
		if err != nil {
			return err
		}

		loginRequest := auth.AuthorizationRequest{
			Realm:  cl.Realm(),
			UserID: email,
		}

		status, err := post(cl.APIEndpoint()+loginEndpoint, &loginRequest, nil)
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

	email := c.Args().Get(0)
	token := c.Args().Get(1)

	cl, err := podops.NewClient(context.TODO(), "")
	if err != nil {
		return err
	}

	authRequest := auth.AuthorizationRequest{
		Realm:  cl.Realm(),
		UserID: email,
		Token:  token,
	}
	response := auth.AuthorizationRequest{}

	status, err := post(cl.APIEndpoint()+authEndpoint, &authRequest, &response)
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

	fmt.Println("logout successful")
	return nil
}
