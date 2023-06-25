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

func RpcRetrieveData(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("RetrieveData RPC called")
	logger.Info("Payload: %s", payload)

	request, err := DeserializePayload(payload)
	if err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 3)
	}

	filePath := GetFilePath(request)
	file, err := os.Open(filePath)
	if err != nil {
		return "", runtime.NewError("Unable to open file", 5)
	}

	return ExecuteRpcRetrieveData(ctx, logger, db, nk, file, request)
}

func ExecuteRpcRetrieveData(ctx context.Context, logger LoggerInterface, db DBExecutorInterface, nk NakamaModuleInterface, reader io.Reader, request PayloadRequest) (string, error) {

	content, err := ReadFileFromDisk(reader)
	if err != nil {
		return "", runtime.NewError("Unable to read file", 13)
	}

	var contentHash = CalculateHash(content)
	request_hash := ConvertNullablePointerToString(request.RequestHash)
	var equalHashes = contentHash == request_hash

	if err := SaveRequestInDatabase(ctx, db, request, equalHashes); err != nil {
		return "", runtime.NewError("Unable to save to database", 13)
	}

	response, err := GenerateResponse(request, request_hash, content, equalHashes)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return response, nil
}
