package storage

import (
	_ "github.com/lib/pq"
)

type BackingStore interface {
	Init() error
	Get(key string) (string, error)
	Set(key string, value string) error
}
