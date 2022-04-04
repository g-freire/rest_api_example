package booking

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

var ctx = context.TODO()

func (p postgres) GetAll(limit, offset string) ([]Booking, error) {
	var bookingCollection []Booking
	var q strings.Builder
	q.WriteString(`
		SELECT * FROM booking
	`)
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
	err := pgxscan.Select(ctx, p.db, &bookingCollection, q.String())
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return bookingCollection, nil
}

func (p postgres) GetByID(id string) (Booking, error) {
	var booking Booking
	sql := `
			SELECT * FROM booking
			WHERE id = $1
		`
	err := pgxscan.Get(ctx, p.db, &booking, sql, id)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return Booking{}, err
	}
	return booking, nil
}

func (p postgres) GetTotalCount() (int64, error) {
	sql := "SELECT COUNT(*) FROM booking"
	var total int64
	err := p.db.QueryRow(context.Background(), sql).Scan(&total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return 0, err
	}
	fmt.Println(total)
	return total, nil
}

func (p postgres) GetByDateRange(startDate, endDate string) ([]Booking, error) {
	var bookingCollection []Booking
	sql := `
		SELECT * FROM booking 
		WHERE start_date >= $1 AND end_date <= $2 
		ORDER BY creation_time DESC`
	err := pgxscan.Select(ctx, p.db, &bookingCollection, sql, startDate, endDate)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return bookingCollection, nil
}

func (p postgres) Save(booking Booking) error {
	// saving with pessimistic concurrency control
	tx, err := p.db.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(context.TODO())

	tsql := `
		 INSERT INTO booking (
				  class_id,
				  member_id, 
				  date
				)
		 VALUES ($1, $2, $3)`

	tran, err := tx.Exec(ctx, tsql,
		booking.ClassId,
		booking.MemberId,
		booking.Date,
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

func (p postgres) Update(id string, booking Booking) error {
	// updating with pessimistic concurrency control
	tx, err := p.db.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(context.TODO())

	tsql := `
		UPDATE booking
		SET name = $1,
			start_date = $2,
			end_date = $3,
			capacity = $4
		WHERE id = $5
		;`

	tran, err := tx.Exec(ctx, tsql,
		booking.Date,
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

func (p postgres) Delete(id string) error {
	sql := `
			DELETE FROM booking
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
