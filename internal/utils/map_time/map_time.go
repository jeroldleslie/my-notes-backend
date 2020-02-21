package map_time

import (
	"github.com/jeroldleslie/my-notes-backend/internal/utils/jsn"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type MapTime struct {
	mtx   sync.Mutex
	times map[string]time.Time
}

func (mp *MapTime) Add(key string) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	if mp.times == nil {
		mp.times = make(map[string]time.Time)
	}
	mp.times[key] = time.Now()
}

func (mp *MapTime) Since(key string) (time.Duration, error) {
	t2 := time.Now()
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	t1, ok := mp.times[key]
	if !ok {
		return 0, errors.Errorf("there is no key '%s' in map %+v", key, jsn.B(mp.times))
	}
	return t2.Sub(t1), nil
}
