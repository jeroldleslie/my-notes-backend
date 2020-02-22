package api

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

type DBMapper struct {
	DB    *pg.DB
	Cache *cache.Cache
}

func pageLimitQ(query *orm.Query, page, pageLimit int) *orm.Query {
	if pageLimit == 0 {
		logrus.Warnf("pageLimit is 0")
		return query
	}
	if page == 0 {
		logrus.Warnf("page is 0")
		return query
	}
	return query.Limit(pageLimit).
		Offset((page - 1) * pageLimit)
}

func (m *DBMapper) InitCache() {
	m.Cache = cache.New(5*time.Minute, 10*time.Minute)
}


