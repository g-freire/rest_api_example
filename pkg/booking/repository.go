package booking

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	pg "gym/internal/db/postgres"
	"gym/pkg/class"
	"gym/pkg/member"
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

func (p postgres) GetAll(ctx context.Context, limit, offset string) ([]Booking, error) {
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

func (p postgres) GetByID(ctx context.Context, id string) (Booking, error) {
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

func (p postgres) GetTotalCount(ctx context.Context) (int64, error) {
	sql := "SELECT COUNT(*) FROM booking"
	var total int64
	err := p.db.QueryRow(context.Background(), sql).Scan(&total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return 0, err
	}
	//fmt.Println(total)
	return total, nil
}

func (p postgres) GetByDateRange(ctx context.Context, startDate, endDate string) ([]Booking, error) {
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

func (p postgres) Save(ctx context.Context, booking Booking) (id int, err error) {
	// saving with pessimistic concurrency control
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Print("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	tsql := `
		 INSERT INTO booking (
				  class_id,
				  member_id, 
				  date
				)
		 VALUES ($1, $2, $3)
		 RETURNING ID 
	`
	// QueryRow is used instead of Exec because of postgres returning property
	err = tx.QueryRow(ctx, tsql,
		booking.ClassId,
		booking.MemberId,
		booking.Date,
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

func (p postgres) Update(ctx context.Context, id string, booking Booking) error {
	// updating with pessimistic concurrency control
	tx, err := p.db.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: "serializable"})
	if err != nil {
		log.Printf("\n[ERROR]: TRANSACTION COULD NOT BEGIN", err)
	}
	defer tx.Rollback(context.TODO())

	tsql := `
		UPDATE booking
		SET class_id = $1,
			member_id = $2,
			date = $3
		WHERE id = $4;
		;`

	tran, err := tx.Exec(ctx, tsql,
		booking.ClassId,
		booking.MemberId,
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
		//fmt.Print("INSERT COMMITED")
	}
	return nil
}

func (p postgres) Delete(ctx context.Context, id string) error {
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

func (p postgres) GetAllClassesByMemberId(ctx context.Context, memberID string) ([]class.Class, error) {
	var classCollection []class.Class
	sql := `
		SELECT name, start_date, end_date, capacity, c.id, c.creation_time
		FROM class c
		INNER JOIN booking b 
		ON c.id = b.class_id
		WHERE b.member_id =$1
		`
	err := pgxscan.Select(ctx, p.db, &classCollection, sql, memberID)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return classCollection, nil
}

func (p postgres) GetAllMembersByClassId(ctx context.Context, id string) ([]member.Member, error) {
	var memberCollection []member.Member
	sql := `
		SELECT name, m.id, m.creation_time
		FROM member m
		INNER JOIN booking b 
		ON m.id = b.member_id
		WHERE b.class_id =$1
		`
	err := pgxscan.Select(ctx, p.db, &memberCollection, sql, id)
	if err != nil {
		log.Printf("\n[ERROR]:", err)
		return nil, err
	}
	return memberCollection, nil
}
