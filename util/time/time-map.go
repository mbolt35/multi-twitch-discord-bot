package time

import (
	"errors"
	"time"

	"github.com/mbolt35/multi-twitch-discord-bot/storage"
)

type TimeMap struct {
	backingStore storage.BackingStore
	timeFormat   string
}

func NewTimeMap(backingStore storage.BackingStore, timeFormat string) *TimeMap {
	instance := TimeMap{
		backingStore: backingStore,
		timeFormat:   timeFormat,
	}

	return &instance
}

func (tm *TimeMap) Exists(key string) bool {
	result, err := tm.backingStore.Get(key)
	return "" == result || nil != err
}

func (tm *TimeMap) Get(key string) (time.Time, error) {
	result, err := tm.backingStore.Get(key)
	if nil != err {
		return time.Time{}, err
	}

	if "" == result {
		return time.Time{}, errors.New("No entry for key")
	}

	return time.Parse(tm.timeFormat, result)
}

func (tm *TimeMap) Set(key string, t string) error {
	_, err := time.Parse(tm.timeFormat, t)
	if nil != err {
		return err
	}

	return tm.backingStore.Set(key, t)
}
