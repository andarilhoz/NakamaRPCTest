package main

import (
	"io"
	"strings"
)

func ReadFileFromDisk(reader io.Reader) (string, error) {
	//if file is too large, this should be replaced to read line by line
	content, err := io.ReadAll(reader)
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
