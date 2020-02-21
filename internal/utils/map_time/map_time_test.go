package map_time_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapTime(t *testing.T) {
	as := assert.New(t)
	mt := map_time.MapTime{}
	mt.Add("key1")
	sleep := 10 * time.Millisecond
	time.Sleep(sleep)
	since, err := mt.Since("key1")
	as.NoError(err)
	t.Logf("since: %s", since)
	as.Truef(since > sleep, "expected to be greater than %s but was %s", sleep, since)
	as.Truef(since < 2*sleep, "expected to be less than %s but was %s", sleep, since)
}
