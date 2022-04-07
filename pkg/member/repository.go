package member

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

func (p postgres) GetAll(ctx context.Context, limit, offset, name string) ([]Member, error) {
	var memberCollection []Member
	var q strings.Builder
	q.WriteString(`
		SELECT * FROM member
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
	err := pgxscan.Select(ctx, p.db, &memberCollection, q.String())
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return memberCollection, nil
}

func (p postgres) GetByID(ctx context.Context, id string) (Member, error) {
	var member Member
	sql := `
			SELECT * FROM member
			WHERE id = $1
		`
	err := pgxscan.Get(ctx, p.db, &member, sql, id)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return Member{}, err
	}
	return member, nil
}

func (p postgres) GetTotalCount(ctx context.Context) (int64, error) {
	sql := "SELECT COUNT(*) FROM member"
	var total int64
	err := p.db.QueryRow(ctx, sql).Scan(&total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return 0, err
	}
	//fmt.Println(total)
	return total, nil
}

func (p postgres) Save(ctx context.Context, member Member) (id int, err error) {
	// saving with pessimistic concurrency control
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	tsql := `
		 INSERT INTO member (name)
		 VALUES ($1)
		 RETURNING id;
	`
	// QueryRow is used instead of Exec because of postgres returning property
	err = tx.QueryRow(ctx, tsql,
		member.Name,
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

func (p postgres) Update(ctx context.Context, id string, member Member) error {
	// updating with pessimistic concurrency control
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(ctx)

	tsql := `
		UPDATE member
		SET name = $1
		WHERE id = $2;`

	tran, err := tx.Exec(ctx, tsql,
		member.Name,
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
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT COMMIT \n", err)
		return err
	} else {
		//fmt.Print("INSERT COMMITED")
	}
	return nil
}

func (p postgres) Delete(ctx context.Context, id string) error {
	sql := `
			DELETE FROM member
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
