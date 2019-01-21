package storage

// PostgresBackingStore is the implementation of BackingStore with Postgres SQL
type MemoryBackingStore struct {
	memory map[string]string
}

// Ensure we correctly implement BackingStore
var _ BackingStore = &MemoryBackingStore{}

// NewMemoryStore creates a new BackingStore implementation using an in memory map
func NewMemoryStore() BackingStore {
	instance := MemoryBackingStore{}
	return &instance
}

// isMapEntry determines if the input key exists in the map
func isMapEntry(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}

func (p *MemoryBackingStore) Init() error {
	p.memory = make(map[string]string)
	return nil
}

func (p *MemoryBackingStore) Get(key string) (string, error) {
	if !isMapEntry(p.memory, key) {
		return "", nil
	}

	return p.memory[key], nil
}

func (p *MemoryBackingStore) Set(key string, value string) error {
	p.memory[key] = value
	return nil
}
