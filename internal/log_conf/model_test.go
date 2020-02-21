package log_conf_test

import (
	"testing"

	"github.com/jeroldleslie/my-notes-backend/internal/log_conf"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAllLevelFiles(t *testing.T) {
	as := assert.New(t)
	if err := log_conf.AllLevelFiles("/tmp", "log_conf_test", logrus.TraceLevel); !as.NoError(err) {
		return
	}
}
