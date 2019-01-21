package storage

// BackingStore implementation prototype for an object capable of retrieving and
// storing strings
type BackingStore interface {
	Init() error
	Get(key string) (string, error)
	Set(key string, value string) error
}
