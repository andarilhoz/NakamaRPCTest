package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

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

func RpcRetrieveData(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("RetrieveData RPC called")
	logger.Info("Payload: %s", payload)

	var request PayloadRequest

	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 3)
	}

	request.PopulateDefaultValues()

	logger.Info("Payload Version: %s", request.RequestVersion)

	content, err := ReadFileFromDisk(logger, request.RequestType, request.RequestVersion)
	if err != nil {
		return "", runtime.NewError("Unable to open file", 5)
	}

	logger.Info("Content: %s", content)

	var contentHash = CalculateHash(content)

	logger.Info("HashResult: %s", contentHash)

	var hashesAreEqual bool
	if request.RequestHash != nil {
		hashesAreEqual = contentHash == *request.RequestHash
	} else {
		hashesAreEqual = false
	}

	if err := SaveRequestInDatabase(db, ctx, logger, request, hashesAreEqual); err != nil {
		return "", runtime.NewError("Unable to save to database", 13)
	}

	var dataContent *string
	if hashesAreEqual {
		dataContent = &content
	} else {
		dataContent = nil
	}

	request_hash := ConvertNullablePointerToString(request.RequestHash)
	responseObject := Response{
		DataType:    request.RequestType,
		DataVersion: request.RequestVersion,
		DataHash:    request_hash,
		DataContent: dataContent,
	}

	response, err := json.Marshal(responseObject)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(response), nil
}

func (request *PayloadRequest) PopulateDefaultValues() {
	if request.RequestType == "" {
		request.RequestType = "core"
	}

	if request.RequestVersion == "" {
		request.RequestVersion = "1.0.0"
	}
}

func ReadFileFromDisk(logger runtime.Logger, requestType string, requestVersion string) (string, error) {
	var builder strings.Builder

	builder.WriteString("/nakama/json_test_files/")
	builder.WriteString(requestType)
	builder.WriteString("/")
	builder.WriteString(requestVersion)
	builder.WriteString(".json")

	var path string = builder.String()

	//if file is too large, this should be replaced to read line by line
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Error opening file, %v", err)
		return "", err
	}

	return string(content), nil
}

func CalculateHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func SaveRequestInDatabase(db *sql.DB, ctx context.Context, logger runtime.Logger, request PayloadRequest, hashesAreEqual bool) error {
	request_hash := ConvertNullablePointerToString(request.RequestHash)
	_, err := db.ExecContext(ctx, `
	INSERT INTO requests (request_type, request_version, request_hash, request_hashes_are_equal)
	VALUES ($1,$2, $3,$4)
	`, request.RequestType, request.RequestVersion, request_hash, hashesAreEqual)
	if err != nil {
		logger.Error("Error: %s", err.Error())
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
