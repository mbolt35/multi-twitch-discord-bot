package time

import (
	"errors"
	"time"

	"github.com/mbolt35/multi-twitch-discord-bot/storage"
)

// TimeMap is an object capable of storing and retrieving time.Time values by key
type TimeMap struct {
	backingStore storage.BackingStore
	timeFormat   string
}

// NewTimeMap creates a new TimeMap using the provided backing storage and time format
func NewTimeMap(backingStore storage.BackingStore, timeFormat string) *TimeMap {
	instance := TimeMap{
		backingStore: backingStore,
		timeFormat:   timeFormat,
	}

	return &instance
}

// Exists returns true if the key exists in storage
func (tm *TimeMap) Exists(key string) bool {
	result, err := tm.backingStore.Get(key)
	return "" == result || nil != err
}

// Get returns the time.Time value for the key in storage.
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
