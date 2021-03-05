package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"google.golang.org/appengine"
)

// Authorized verifies that clientID can access resource kind/GUID
func Authorized(c echo.Context, role string) (int, error) {
	//return fmt.Errorf("Not allowed to access '%s/%s'", kind, guid)
	return http.StatusOK, nil // FIXME this is just a placeholder
}

// GetBearerToken extracts the bearer token
func GetBearerToken(c echo.Context) string {

	auth := c.Request().Header.Get("Authorization")
	if len(auth) == 0 {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return ""
	}

	if parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}

// GetClientID extracts the ClientID from the token
func GetClientID(c echo.Context) (string, error) {
	token := GetBearerToken(c)
	if token == "" {
		return "", a.ErrNoToken
	}
	auth, err := FindAuthorization(appengine.NewContext(c.Request()), token)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", a.ErrNotAuthorized
	}

	return auth.ClientID, nil
}
