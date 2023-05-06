package handlers

import (
	"database/sql"
	"fmt"
)

func commitOrRollback(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error rolling back transaction: %w", rollbackErr)
		}
		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("error committing transaction: %w", commitErr)
	}
	return nil
}
