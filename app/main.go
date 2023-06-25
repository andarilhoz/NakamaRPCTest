package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	createDBErr := CreateRequestTableIfNotExists(db, ctx, logger)
	if createDBErr != nil {
		return createDBErr
	}

	initializeRPCErr := initializer.RegisterRpc("retrievedata", RpcRetrieveData)
	if initializeRPCErr != nil {
		return initializeRPCErr
	}

	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
	return nil
}

func CreateRequestTableIfNotExists(db *sql.DB, ctx context.Context, logger runtime.Logger) error {

	content, readErr := os.ReadFile("/nakama/scripts/create_request_table.sql")
	if readErr != nil {
		logger.Error("Error opening sql file, %v", readErr)
		return readErr
	}

	_, execErr := db.ExecContext(ctx, string(content))
	if execErr != nil {
		logger.Error("Error: %s", execErr.Error())
	}

	return nil
}
