package class

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	pg "gym/internal/db/postgres"
	"log"
	"os"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *postgres {
	return &postgres{
		db: db,
	}
}

var ctx = context.TODO()

func (p postgres) GetAll(limit, offset string) ([]Class, error) {
	var classCollection []Class
	// request with pagination arguments
	if limit != "" && offset != "" {
		sql := `
			SELECT * FROM class 
			ORDER BY created_at DESC 
			LIMIT $1 
			OFFSET $2
		`
		err := pgxscan.Select(ctx, p.db, &classCollection, sql, limit, offset)
		if err != nil {
			log.Printf("\n[ERROR]:", err)
			return nil, err
		}
	} else { // request without pagination arguments
		sql := `
			SELECT * FROM class 
			ORDER BY created_at DESC 
		`
		err := pgxscan.Select(ctx, p.db, &classCollection, sql)
		if err != nil {
			log.Printf("\n[ERROR]:", err)
			return nil, err
		}
	}
	return classCollection, nil
}

func (p postgres) GetTotalCount() (int64, error) {
	sql := "SELECT COUNT(*) FROM class"
	var total int64
	err := p.db.QueryRow(context.Background(), sql).Scan(&total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return 0, err
	}
	fmt.Println(total)
	return total, nil
}

func (p postgres) GetByDateRange(startDate, endDate string) ([]Class, error) {
	var classCollection []Class
	sql := `
		SELECT * FROM class 
		WHERE start_date >= $1 AND end_date <= $2 
		ORDER BY created_at DESC`
	err := pgxscan.Select(ctx, p.db, &classCollection, sql, startDate, endDate)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return classCollection, nil
}

func (p postgres) Save(class Class) error {
	// saving with pessimistic concurrency control
	tx, err := p.db.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(context.TODO())

	tsql := `
		 INSERT INTO class (
				  name,
				  start_date, 
				  end_date, 
				  capacity
				  )
		 VALUES ($1, $2, $3, $4)`

	tran, err := tx.Exec(ctx, tsql,
		class.Name,
		class.StartDate,
		class.EndDate,
		class.Capacity,
	)
	if err != nil {
		pg.RollbackTxPgx(tx, err)
		return err
	}
	rowsAffected := tran.RowsAffected()
	if err != nil || rowsAffected != 1 {
		log.Print(err)
		pg.RollbackTxPgx(tx, err)
		return err
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT COMMIT \n", err)
		return err
	} else {
		fmt.Print("\n INSERT COMMITED")
	}
	return nil
}

func (p postgres) Update(id string, class Class) error {
	// updating with pessimistic concurrency control
	tx, err := p.db.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(context.TODO())

	tsql := `
		UPDATE class
		SET name = $1,
			start_date = $2,
			end_date = $3,
			capacity = $4
		WHERE id = $5
		;`

	tran, err := tx.Exec(ctx, tsql,
		class.Name,
		class.StartDate,
		class.EndDate,
		class.Capacity,
		id,
	)
	if err != nil {
		pg.RollbackTxPgx(tx, err)
		return err
	}
	rowsAffected := tran.RowsAffected()
	if err != nil || rowsAffected != 1 {
		log.Print(err)
		pg.RollbackTxPgx(tx, err)
		return err
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT COMMIT \n", err)
		return err
	} else {
		fmt.Print("\n INSERT COMMITED")
	}
	return nil
}

func (p postgres) GetByID(id string) (Class, error) {
	var class Class
	sql := `
			SELECT * FROM class
			WHERE ID = $1
		`
	err := pgxscan.Get(ctx, p.db, &class, sql, id)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return Class{}, err
	}
	return class, nil
}

func (p postgres) GetByName(name string) ([]Class, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
