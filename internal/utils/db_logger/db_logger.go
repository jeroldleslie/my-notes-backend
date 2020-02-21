package db_logger

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/go-pg/pg"
	"github.com/jeroldleslie/my-notes-backend/internal/utils/map_time"
	"github.com/sirupsen/logrus"
)

type dbLogger struct {
	LogExecutionTime bool
}

var times map_time.MapTime

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	if !d.LogExecutionTime {
		return
	}
	key, err := queryKey(q)
	if err != nil {
		logrus.Error(err)
	}
	times.Add(key)
	qStr, err := q.FormattedQuery() // Needed when PG Query is failing
	if err != nil {
		logrus.Error(err)
	}
	logrus.Tracef("query to execute: %s", qStr)
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	qStr, err := q.FormattedQuery()
	if err != nil {
		logrus.Error(err)
	}
	var aboutTimeSpent string
	aboutTimeSpent, err = func() (string, error) {
		if !d.LogExecutionTime {
			return "", nil
		}
		key, err := queryKey(q)
		if err != nil {
			return "", errors.WithStack(err)
		}
		since, err := times.Since(key)
		if err != nil {
			err = errors.WithStack(err)
			return "", err
		}
		return fmt.Sprintf("%s - ", since), nil
	}()
	if err != nil {
		logrus.Errorf("%+v", err)
		aboutTimeSpent = "couldn't get query time - "
	}
	logrus.Tracef("%s executed query: %s", aboutTimeSpent, qStr)
}

func AddDbLogger(db *pg.DB, logExecutionTime bool) {
	logrus.Infof("adding db logger, logExecutionTime: %+v", logExecutionTime)
	db.AddQueryHook(dbLogger{LogExecutionTime: logExecutionTime})
}

func queryKey(q *pg.QueryEvent) (string, error) {
	return fmt.Sprintf("%p", q), nil
}
