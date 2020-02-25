package api

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-pg/pg"
)

type DBMapper struct {
	DB    *pg.DB
	Cache *cache.Cache
}

func (m *DBMapper) InitCache() {
	m.Cache = cache.New(5*time.Minute, 10*time.Minute)
}
