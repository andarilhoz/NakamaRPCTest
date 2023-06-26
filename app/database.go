package main

import (
	"context"
	"database/sql"
)

type DBExecutorInterface interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func SaveRequestInDatabase(ctx context.Context, db DBExecutorInterface, request PayloadRequest, hashesAreEqual bool) error {
	request_hash := ConvertNullablePointerToString(request.RequestHash)
	_, err := db.ExecContext(ctx, `
	INSERT INTO requests (request_type, request_version, request_hash, request_hashes_are_equal)
	VALUES ($1,$2,$3,$4)
	`, request.RequestType, request.RequestVersion, request_hash, hashesAreEqual)
	if err != nil {
		return err
	}

	return nil
}
