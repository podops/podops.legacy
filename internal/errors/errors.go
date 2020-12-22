package errors

//
// This is a copy of txsvc/commons/pkg/error, with lots of changes and improvements.
// Will be eventually backported.
//
//
// References
// https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
// https://godoc.org/github.com/pkg/errors
//

import (
	"net/http"
	"runtime"
	"strings"
)

type (

	// StatusObject is used to report status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status  int    `json:"status" binding:"required"`
		Message string `json:"message" binding:"required"`
		// location of an error
		Pkg string `json:"pkg,omitempty"`
		Fn  string `json:"fn,omitempty"`
	}
)

func (so *StatusObject) Error() string {
	return so.Message
}

// New returns an error that formats as the given text
func New(msg string, s int) error {
	p, f := packageAndFunc()
	return &StatusObject{
		Status:  s,
		Message: msg,
		Pkg:     p,
		Fn:      f,
	}
}

// Wrap adds some context to an error
func Wrap(e error) error {
	p, f := packageAndFunc()
	return &StatusObject{
		Status:  http.StatusInternalServerError,
		Message: e.Error(),
		Pkg:     p,
		Fn:      f,
	}
}

// see https://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func packageAndFunc() (string, string) {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkg := ""
	fn := parts[pl-1]
	if parts[pl-2][0] == '(' {
		fn = parts[pl-2] + "." + fn
		pkg = strings.Join(parts[0:pl-2], ".")
	} else {
		pkg = strings.Join(parts[0:pl-1], ".")
	}
	return pkg, fn
}
