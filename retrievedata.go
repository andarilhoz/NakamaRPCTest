package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"

	"io"
	"os"

	"github.com/heroiclabs/nakama-common/runtime"
)

type PayloadRequest struct {
	RequestType    string  `json:"type"`
	RequestVersion string  `json:"version"`
	RequestHash    *string `json:"hash"`
}

type Response struct {
	DataType    string  `json:"type"`
	DataVersion string  `json:"version"`
	DataHash    string  `json:"hash"`
	DataContent *string `json:"content"`
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

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

func ExecuteRpcRetrieveData(ctx context.Context, logger LoggerInterface, db DBExecutor, nk NakamaModuleInterface, reader io.Reader, request PayloadRequest) (string, error) {

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

func DeserializePayload(payload string) (PayloadRequest, error) {
	var request PayloadRequest

	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return PayloadRequest{}, err
	}

	request.PopulateDefaultValues()
	return request, nil
}

func (request *PayloadRequest) PopulateDefaultValues() {
	if request.RequestType == "" {
		request.RequestType = "core"
	}

	if request.RequestVersion == "" {
		request.RequestVersion = "1.0.0"
	}
}

func ReadFileFromDisk(reader io.Reader) (string, error) {
	//if file is too large, this should be replaced to read line by line
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetFilePath(request PayloadRequest) string {
	var builder strings.Builder

	builder.WriteString("/nakama/json_test_files/")
	builder.WriteString(request.RequestType)
	builder.WriteString("/")
	builder.WriteString(request.RequestVersion)
	builder.WriteString(".json")

	return builder.String()
}

func CalculateHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func SaveRequestInDatabase(ctx context.Context, db DBExecutor, request PayloadRequest, hashesAreEqual bool) error {
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

func ConvertNullablePointerToString(pointer *string) string {
	if pointer != nil {
		return *pointer
	} else {
		return ""
	}
}

func GenerateResponse(request PayloadRequest, request_hash string, content string, equalHashes bool) (string, error) {
	var dataContent *string
	if equalHashes {
		dataContent = &content
	} else {
		dataContent = nil
	}

	responseObject := Response{
		DataType:    request.RequestType,
		DataVersion: request.RequestVersion,
		DataHash:    request_hash,
		DataContent: dataContent,
	}

	response, err := json.Marshal(responseObject)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
