package cache

import "time"

var nowFunc = time.Now

func expired(t time.Time) bool {
	if t.IsZero() {
		return false
	}

	return nowFunc().After(t)
}
