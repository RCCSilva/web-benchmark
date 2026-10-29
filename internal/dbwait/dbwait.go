// Package dbwait polls a freshly-created Postgres instance until it accepts
// connections and creates the "users" table the benchmark apps expect,
// mirroring the previous src/db.ts.
package dbwait

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Conn struct {
	User     string
	Password string
	Host     string
	Port     int32
	Name     string
}

const createTableSQL = `
CREATE TABLE "users" (
  ID    BIGSERIAL         PRIMARY KEY,
  NAME  VARCHAR(64)                  ,
  AGE   INT
)`

// WaitAndCreateTable retries connecting until the database is reachable,
// then creates the users table it expects. It gives up when ctx is done.
func WaitAndCreateTable(ctx context.Context, conn Conn, debug bool) error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conn.User, conn.Password, conn.Host, conn.Port, conn.Name)

	for {
		err := tryCreateTable(ctx, url)
		if err == nil {
			return nil
		}
		if debug {
			fmt.Printf("   -> %v\n", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func tryCreateTable(ctx context.Context, url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, createTableSQL)
	return err
}
