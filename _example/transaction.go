package _example

import (
	"github.com/best-expendables/common-utils/transaction"
	"context"
)

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

	// Init Transaction
	tnxManager := transaction.NewTxManager(db)

	// Use transaction
	var err error
	// Start db transaction
	tnx, ctx := tnxManager.Start(context.Background())
	defer func() {
		err = tnx.Finish(err)
	}()

	// Call service
	err = putRecord(ctx)
}

func putRecord(ctx context.Context) error {
	return nil
}
