package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func RollbackTxPgx(tx pgx.Tx, err error) {
	if rollbackErr := tx.Rollback(context.TODO()); rollbackErr != nil {
		log.Printf("\n[ERROR]: UNABLE TO ROLLBACK \n", rollbackErr)
	}
	log.Printf("\n[ERROR]: TRANSACTION COULD NOT EXEC CONTEXT \n", err)
}
