package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/jeroldleslie/my-notes-backend/internal/utils/db_logger"
	"github.com/jeroldleslie/my-notes-backend/internal/utils/stage"
	"github.com/labstack/echo"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const TestDB = "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
const DevRedisDB = "localhost:6379"
const TestRedisDB = "localhost:6379"
const RedisCacheTime = 300

// ConnectToPostgres connects to postgres instance
func ConnectToPostgres(connectionString string) (*pg.DB, error) {
	if connectionString == "" {
		connectionString = os.Getenv("POSTGRESQL_ADDRESS")
		logrus.Infof("connectionString from env POSTGRESQL_ADDRESS: '%+v'", connectionString)
	}
	if connectionString == "" {
		err := fmt.Errorf("missing connectionString")
		return nil, errors.Wrap(err, "POSTGRESQL_ADDRESS env is empty")
	}
	opt, err := pg.ParseURL(connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to postgres with connection string: "+connectionString)
	}

	db := pg.Connect(opt)
	_, err = db.Exec("SELECT 1")
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	if !stage.IsProd() {
		db_logger.AddDbLogger(db, true)
	} else {
		logrus.Infof("don't add db logger")
	}

	return db, nil
}

// ConnectToPostgres connects to postgres instance
func ConnectToPostgresTimeout(connectionString string, timeout, retry time.Duration) (*pg.DB, error) {
	var (
		connectionError error
		db              *pg.DB
	)
	connected := make(chan bool)
	go func() {
		for {
			db, connectionError = ConnectToPostgres(connectionString)
			if connectionError != nil {
				time.Sleep(retry)
				continue
			}
			connected <- true
			break
		}
	}()
	select {
	case <-time.After(timeout):
		err := errors.Wrapf(connectionError, "timeout %s connecting to db", timeout)
		return nil, err
	case <-connected:
	}
	return db, nil
}

// RedisGet retrieves value by URI and sends response
func RedisGet(c *echo.Context, redisDB *redis.Client) error {
	key := fmt.Sprintf("%s?%s", (*c).Path(), (*c).QueryString())
	cachedResponse, err := redisDB.Get(key).Result()
	if err != nil {
		return nil
	}
	if cachedResponse != "" {
		response := Response{
			StatusCode: http.StatusOK,
			Data:       cachedResponse,
		}
		return response.Send(c)
	}
	return nil
}

// RedisSet sets value by URI
func RedisSet(res interface{}, c *echo.Context, redisDB *redis.Client) error {
	key := fmt.Sprintf("%s?%s", (*c).Path(), (*c).QueryString())
	var val []byte
	if reflect.TypeOf(res) == reflect.TypeOf("") {
		val = []byte(res.(string))
	} else {
		stringed, err := json.Marshal(res)
		if err != nil {
			return err
		}
		val = stringed
	}
	err := redisDB.Set(key, val, RedisCacheTime*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
