package connection

import (
	"github.com/best-expendables/logger"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type dbRetryConf struct {
	ErrMsgKey string
	Retry     int
	Delay     time.Duration
}

type dbWithRetry struct {
	retryConf []dbRetryConf
	*sql.DB
}

type DBWithRetryConf struct {
	DefaultRetry int
	DefaultDelay time.Duration
}

//Define DBRetryConf
func makeRetryConfig(conf DBWithRetryConf) []dbRetryConf {
	var retryConf []dbRetryConf
	// Connection Reset By Peer
	retryConf = append(retryConf, dbRetryConf{"connection reset by peer", conf.DefaultRetry, conf.DefaultDelay})
	retryConf = append(retryConf, dbRetryConf{"write: broken pipe", conf.DefaultRetry, conf.DefaultDelay})
	//....
	return retryConf
}

func NewDBWithRetry(db *sql.DB, conf DBWithRetryConf) dbWithRetry {
	return dbWithRetry{
		retryConf: makeRetryConfig(conf),
		DB:        db,
	}
}

func (d dbWithRetry) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	err = d.retry(func() error {
		result, err = d.DB.Exec(query, args...)
		return err
	}, err)
	return result, err
}

func (d dbWithRetry) Prepare(query string) (*sql.Stmt, error) {
	result, err := d.DB.Prepare(query)
	err = d.retry(func() error {
		result, err = d.DB.Prepare(query)
		return err
	}, err)
	return result, err
}

func (d dbWithRetry) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := d.DB.Query(query, args...)
	err = d.retry(func() error {
		result, err = d.DB.Query(query, args...)
		return err
	}, err)
	return result, err
}
func (d dbWithRetry) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

func (d dbWithRetry) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	result, err := d.DB.BeginTx(ctx, opts)
	err = d.retry(func() error {
		result, err = d.DB.BeginTx(ctx, opts)
		return err
	}, err)
	return result, err
}

func (d dbWithRetry) retry(f func() error, err error) error {
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
				logger.Error(errors.Wrap(err, "DB Retry Error"))
				time.Sleep(retryConf.Delay)
			}
		}
	}
	return err
}
