package postgres

import (
	"context"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gym/internal/constants"
	"log"
	"sync"
)

var (
	postgresClientPool *ClientPool
	postgresOnce       sync.Once
)

type ClientPool struct {
	Conn *pgxpool.Pool
}

func NewPostgresConnectionPool(dbHost string) *pgxpool.Pool {
	//fmt.Println("Connecting to Postgres at", dbHost)
	postgresOnce.Do(func() {
		config, err := pgxpool.ParseConfig(dbHost)
		//config.MinConns = 25
		pool, err := pgxpool.ConnectConfig(context.Background(), config)
		if err != nil {
			log.Fatal(constants.Red, "Couldn't connect to the database. Reason ", err, constants.Reset)
		}
		pool.Stat()
		err = Ping(pool)
		if err != nil {
			log.Fatal(constants.Red, "Couldn't ping the database. Reason ", err, constants.Reset)
		}
		postgresClientPool = &ClientPool{Conn: pool}
	})
	log.Print(constants.Blue, "DB CONNECTION CREATED at: ", dbHost, constants.Reset)
	return postgresClientPool.Conn
}

// Ping acquires a connection from the Pool and executes an empty sql statement against it.
// If the sql returns without error, the database Ping is considered successful, otherwise, the error is returned.
func Ping(pool *pgxpool.Pool) error {
	c, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer c.Release()
	return c.Conn().Ping(context.Background())
}
