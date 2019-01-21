package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	// DatabasePGType is the SQL driver type used when opening a connection
	DatabasePGType string = "postgres"

	// CreateTableSql creates a new stream
	CreateTableSql string = `CREATE TABLE IF NOT EXISTS store (
                                key varchar(255) not null,
                                value varchar(255) not null,
                                PRIMARY KEY(key));`

	// GetQuery is the SQL which looks up a value from storage given a key
	GetQuery string = "SELECT value FROM store WHERE key=$1"

	// SetStatement is the SQL which inserts a new value for the provided key
	SetStatement string = `INSERT INTO store(key, value) VALUES($1, $2)
                                ON CONFLICT (key) DO UPDATE SET value=$2`
)

// PostgresBackingStore is the implementation of BackingStore with Postgres SQL
type PostgresBackingStore struct {
	databaseHost string
	db           *sql.DB
	getQuery     *sql.Stmt
	setStatement *sql.Stmt
}

// Ensure we correctly implement BackingStore
var _ BackingStore = &PostgresBackingStore{}

// NewPostgresStore creates a new BackingStore implementation using Postgres SQL
func NewPostgresStore(databaseHost string) BackingStore {
	instance := PostgresBackingStore{
		databaseHost: databaseHost,
	}

	return &instance
}

func (p *PostgresBackingStore) Init() error {
	db, err := sql.Open(DatabasePGType, p.databaseHost)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(CreateTableSql)
	if nil != err {
		return err
	}

	p.db = db
	p.getQuery, _ = db.Prepare(GetQuery)
	p.setStatement, _ = db.Prepare(SetStatement)
	return nil
}

func (p *PostgresBackingStore) Get(key string) (string, error) {
	rows, err := p.getQuery.Query(key)
	if nil != err {
		return "", err
	}

	defer rows.Close()

	// No Entries
	if !rows.Next() {
		return "", nil
	}

	var value string
	err = rows.Scan(&value)
	if nil != err {
		return "", err
	}

	return value, nil
}

func (p *PostgresBackingStore) Set(key string, value string) error {
	_, err := p.setStatement.Exec(key, value)
	return err
}
