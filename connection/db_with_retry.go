package connection

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

type DBRetryConf struct {
	ErrMsgKey string
	Retry     int
	Delay     time.Duration
}

type DBWithRetry struct {
	retryConf []DBRetryConf
	*sql.DB
}

type PostgresConfig struct {
	DefaultRetry int
	DefaultDelay time.Duration
}

//Define DBRetryConf
func makeRetryConfig(conf PostgresConfig) []DBRetryConf {
	var retryConf []DBRetryConf
	// Connection Reset By Peer
	retryConf = append(retryConf, DBRetryConf{"connection reset by peer", conf.DefaultRetry, conf.DefaultDelay})
	retryConf = append(retryConf, DBRetryConf{"write: broken pipe", conf.DefaultRetry, conf.DefaultDelay})
	//....
	return retryConf
}

func NewDBWithRetry(db *sql.DB, conf PostgresConfig) DBWithRetry {
	return DBWithRetry{
		retryConf: makeRetryConfig(conf),
		DB:        db,
	}
}

func (d DBWithRetry) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	err = d.retry(func() error {
		result, err = d.DB.Exec(query, args...)
		return err
	}, err)
	return result, err
}

func (d DBWithRetry) Prepare(query string) (*sql.Stmt, error) {
	result, err := d.DB.Prepare(query)
	err = d.retry(func() error {
		result, err = d.DB.Prepare(query)
		return err
	}, err)
	return result, err
}

func (d DBWithRetry) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := d.DB.Query(query, args...)
	err = d.retry(func() error {
		result, err = d.DB.Query(query, args...)
		return err
	}, err)
	return result, err
}
func (d DBWithRetry) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

func (d DBWithRetry) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	result, err := d.DB.BeginTx(ctx, opts)
	err = d.retry(func() error {
		result, err = d.DB.BeginTx(ctx, opts)
		return err
	}, err)
	return result, err
}

func (d DBWithRetry) retry(f func() error, err error) error {
	if err == nil {
		return nil
	}
	for _, retryConf := range d.retryConf {
		if strings.Contains(err.Error(), retryConf.ErrMsgKey) {
			for i := 0; i < retryConf.Retry; i++ {
				err = f()
				if err == nil {
					return nil
				}
				log.Fatal(errors.Wrap(err, "DB Retry Error"))
				time.Sleep(retryConf.Delay)
			}
		}
	}
	return err
}
