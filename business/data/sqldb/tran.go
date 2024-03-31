package sqldb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/jmoiron/sqlx"
)

// dbBeginner implements the core coreTransaction interface.
type dbBeginner struct {
	sqlxDB *sqlx.DB
}

// NewBeginner constructs a value that implements the database
// beginner interface.
func NewBeginner(sqlxDB *sqlx.DB) transaction.Beginner {
	return &dbBeginner{
		sqlxDB: sqlxDB,
	}
}

// Begin starts a transaction and returns a value that implements
// the core transactor interface.
func (db *dbBeginner) Begin() (transaction.Transaction, error) {
	return db.sqlxDB.Beginx()
}

// GetExtContext is a helper function extracts the sqlx value
// from the core transactor interface for transactional use.
func GetExtContext(tx transaction.Transaction) (sqlx.ExtContext, error) {
	ec, ok := tx.(sqlx.ExtContext)
	if !ok {
		return nil, fmt.Errorf("Transactor(%T) not of a type *sqlx.Tx", tx)
	}

	return ec, nil
}