package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenConns, maxIdeleConns int, maxIdleTime string) (*sql.DB, error) {

	// open connection
	db, err :=sql.Open("postgres", addr)
	if err!= nil{
		return nil, err 
	}

	// duration
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdeleConns)
	duration, err:= time.ParseDuration(maxIdleTime)
	if err !=nil{
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ping to connection
	if err = db.PingContext(ctx); err!=nil{
		return nil, err
	}

	return db, nil
}
