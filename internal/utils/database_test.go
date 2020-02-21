package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToPostgres(t *testing.T) {
	var tests = []struct {
		ConnectionString string
		ExpectedError    error
	}{
		{"", nil},
		{"invalid", errors.New("connecting to postgres: pg: invalid scheme: ")},
		{TestDB, nil},
	}
	connString := os.Getenv("POSTGRESQL_CRYPTO_ADDRESS")
	for _, test := range tests {
		if test.ConnectionString == "" {
			os.Setenv("POSTGRESQL_CRYPTO_ADDRESS", TestDB)
		}
		_, err := ConnectToPostgres(test.ConnectionString)
		if test.ExpectedError != nil {
			assert.EqualError(t, err, test.ExpectedError.Error())
		} else {
			assert.Equal(t, nil, err)
		}
	}
	os.Setenv("POSTGRESQL_CRYPTO_ADDRESS", connString)
}
