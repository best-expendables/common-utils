package _example

import (
	"github.com/best-expendables/common-utils/connection"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type PostgresConfig struct {
	Host          string
	Port          int
	DbName        string
	User          string
	Pass          string
	LogEnable     bool
	MaxConnection int
	DefaultRetry  int
	DefaultDelay  time.Duration
}

func main() {
	postgresConfig := PostgresConfig{
		Host:         "127.0.0.1",
		Port:         5432,
		DbName:       "app_db",
		User:         "app_user",
		Pass:         "app_pass",
		LogEnable:    true,
		DefaultRetry: 10,
		DefaultDelay: 10,
	}

	rawDb := CreateRawDB(postgresConfig)
	db := CreateDBWithRetry(rawDb, postgresConfig)
	fmt.Println(db)
}

func CreateRawDB(conf PostgresConfig) *sql.DB {
	dbSource := fmt.Sprintf(
		"host=%s user=%s port=%d dbname=%s sslmode=disable password=%s",
		conf.Host,
		conf.User,
		conf.Port,
		conf.DbName,
		conf.Pass,
	)
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxConnection)
	return db
}

func CreateDBWithRetry(db *sql.DB, conf PostgresConfig) *gorm.DB {
	dbWithRetryConf := connection.DBWithRetryConf{
		DefaultRetry: conf.DefaultRetry,
		DefaultDelay: conf.DefaultDelay,
	}
	retryDB := connection.NewDBWithRetry(db, dbWithRetryConf)
	c, err := gorm.Open("postgres", retryDB)
	if err != nil {
		panic(err)
	}
	return c
}
