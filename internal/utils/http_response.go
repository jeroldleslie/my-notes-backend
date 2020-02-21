package utils

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo"
)

// Response maps response into detailed error report
type Response struct {
	Timestamp   int64       `json:"timestamp,omitempty"`
	StatusCode  int         `json:"status_code,omitempty"`
	Error       error       `json:"-"`
	ErrorString string      `json:"error,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// Send handles sending response to the client
// If HTTP status code is 200, it will include only data as response (since that's what webapp expects at the moment)
// For all other errors it will contain the timestamp, code and error message as response
func (res *Response) Send(c *echo.Context) error {
	res.Timestamp = time.Now().UTC().Unix()
	(*c).Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
	logrus.Tracef("res.Error: %+v", res.Error)
	if res.Error != nil {
		logrus.Tracef("res.Error: %+v, res.Error.Error(): %+v", res.Error, res.Error.Error())
		res.ErrorString = res.Error.Error()
	}
	if res.StatusCode != http.StatusOK {
		logrus.Tracef("res: %+v", res)
		return (*c).JSON(res.StatusCode, res.ErrorString)
	}
	if reflect.TypeOf(res.Data) == reflect.TypeOf("") {
		logrus.Tracef("res.Data: %+v", res.Data)
		return (*c).JSONBlob(res.StatusCode, []byte(res.Data.(string)))
	}
	logrus.Tracef("res.Data: %+v", res.Data)
	return (*c).JSON(res.StatusCode, res.Data)
}
