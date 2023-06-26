package main

import (
	"context"
	"database/sql"

	"io"
	"os"

	"github.com/heroiclabs/nakama-common/runtime"
)

type NakamaModuleInterface interface {
}

type LoggerInterface interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

const (
	INVALID_ARGUMENT = 3
	NOT_FOUND        = 5
	INTERNAL         = 13
)

var (
	errBadInput      = runtime.NewError("Unable to read arguments, please verify input", INVALID_ARGUMENT)
	errInternalError = runtime.NewError("Internal Server Error", INTERNAL)
	errFileFound     = runtime.NewError("File not found", NOT_FOUND)
)

func RpcRetrieveData(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("RetrieveData RPC called")
	logger.Info("Payload: %s", payload)

	request, err := DeserializePayload(payload)
	if err != nil {
		return "", errBadInput
	}

	filePath := GetFilePath(request)
	file, err := os.Open(filePath)
	if err != nil {
		return "", errFileFound
	}

	return ExecuteRpcRetrieveData(ctx, logger, db, nk, file, request)
}

func ExecuteRpcRetrieveData(ctx context.Context, logger LoggerInterface, db DBExecutorInterface, nk NakamaModuleInterface, reader io.Reader, request PayloadRequest) (string, error) {

	content, err := ReadFileFromDisk(reader)
	if err != nil {
		return "", errInternalError
	}

	var contentHash = CalculateHash(content)
	request_hash := ConvertNullablePointerToString(request.RequestHash)
	var equalHashes = contentHash == request_hash

	if err := SaveRequestInDatabase(ctx, db, request, equalHashes); err != nil {
		return "", errInternalError
	}

	response, err := GenerateResponse(request, request_hash, content, equalHashes)
	if err != nil {
		return "", errInternalError
	}

	return response, nil
}
