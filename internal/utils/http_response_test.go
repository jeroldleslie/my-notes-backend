package utils

import (
	"errors"
	"net/http"
	"testing"

	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	var tests = []struct {
		Response *Response
	}{
		{&Response{0, http.StatusOK, nil, "", 0}},
		{&Response{0, http.StatusNotFound, nil, "", 0}},
		{&Response{0, http.StatusOK, errors.New("some error"), "", 0}},
	}

	for _, test := range tests {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := test.Response.Send(&c)

		assert.Equal(t, nil, err)
		if test.Response.Error != nil {
			assert.Equal(t, test.Response.Error.Error(), test.Response.ErrorString)
		}
		assert.NotEqual(t, 0, test.Response.Timestamp)
	}
}
