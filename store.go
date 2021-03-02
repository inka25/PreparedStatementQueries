package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

const (
	MYSQL_DUPLICATE_PRIMARY_KEY = 1062
)

type StoreImpl struct {
	db     *sqlx.DB
}

func NewStore(db *sqlx.DB) *StoreImpl {
	return &StoreImpl{
		db:     db,
	}
}

func (s StoreImpl) StartTransaction(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
}

func (s StoreImpl) CommitTransaction(ctx context.Context, tx *sqlx.Tx) error {
	return tx.Commit()
}

func (s StoreImpl) RollbackTransaction(ctx context.Context, tx *sqlx.Tx) error {
	return tx.Rollback()
}





