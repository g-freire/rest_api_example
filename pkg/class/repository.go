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
	"strings"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *postgres {
	return &postgres{
		db: db,
	}
}

func (p postgres) GetAll(ctx context.Context, limit, offset, name string) ([]Class, error) {
	var classCollection []Class
	var q strings.Builder
	q.WriteString(`
		SELECT * FROM class
	`)
	// filter by name
	if name != "" {
		q.WriteString(fmt.Sprintf("WHERE name='%s' ", name))
	}
	q.WriteString("ORDER BY creation_time DESC ")
	// request with pagination - default limit is 200 if no limit is specified
	if limit == "" {
		limit = "200"
	}
	q.WriteString(fmt.Sprintf("LIMIT %s ", limit))
	if offset != "" {
		q.WriteString("OFFSET ")
		q.WriteString(offset)
	}
	err := pgxscan.Select(ctx, p.db, &classCollection, q.String())
	if err != nil {
		log.Print("\n[ERROR]:", err)
		return nil, err
	}
	return classCollection, nil
}

func (p postgres) GetByID(ctx context.Context, id string) (Class, error) {
	var class Class
	sql := `
			SELECT * FROM class
			WHERE id = $1
		`
	err := pgxscan.Get(ctx, p.db, &class, sql, id)
	if err != nil {
		log.Print("\n[ERROR]:", err)
		return Class{}, err
	}
	return class, nil
}

func (p postgres) GetTotalCount(ctx context.Context,) (int64, error) {
	sql := "SELECT COUNT(*) FROM class"
	var total int64
	err := p.db.QueryRow(context.Background(), sql).Scan(&total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return 0, err
	}
	return total, nil
}

func (p postgres) GetByDateRange(ctx context.Context, startDate, endDate string) ([]Class, error) {
	var classCollection []Class
	sql := `
		SELECT * FROM class 
		WHERE start_date >= $1 AND end_date <= $2 
		ORDER BY creation_time DESC`
	err := pgxscan.Select(ctx, p.db, &classCollection, sql, startDate, endDate)
	if err != nil {
		log.Print("\n[ERROR]:", err)
		return nil, err
	}
	return classCollection, nil
}

func (p postgres) Save(ctx context.Context, class Class) (id int, err error) {
	// saving with pessimistic concurrency control
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	tsql := `
		 INSERT INTO class (
				  name,
				  start_date, 
				  end_date, 
				  capacity
				  )
		 VALUES ($1, $2, $3, $4)
		 RETURNING id
`
	// QueryRow is used instead of Exec because of postgres returning property
	err = tx.QueryRow(ctx, tsql,
		class.Name,
		class.StartDate,
		class.EndDate,
		class.Capacity,
	).Scan(&id)
	if err != nil {
		pg.RollbackTxPgx(tx, err)
		log.Print("\n[ERROR]: TRANSACTION ERROR", err)
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT COMMIT \n", err)
		return 0, err
	} else {
		//fmt.Print("INSERT COMMITED")
	}
	return id, nil
}

func (p postgres) Update(ctx context.Context, id string, class Class) error {
	// updating with pessimistic concurrency control
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(ctx)

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
	err = tx.Commit(ctx)
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT COMMIT \n", err)
		return err
	} else {
		//fmt.Print("INSERT COMMITED")
	}
	return nil
}

func (p postgres) Delete(ctx context.Context, id string) error {
	sql := `
			DELETE FROM class
			WHERE id = $1;
	`
	res, err := p.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	rowsAffected := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		log.Print(err)
		return err
	}
	return nil
}
