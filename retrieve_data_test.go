package main

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type MockLogger struct {
}

func (m *MockLogger) Debug(format string, v ...interface{}) {}
func (m *MockLogger) Info(format string, v ...interface{})  {}
func (m *MockLogger) Warn(format string, v ...interface{})  {}
func (m *MockLogger) Error(format string, v ...interface{}) {}

type MockNakamaModule struct {
}

type BrokenReader struct{}

func (br *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Testing error")
}

func TestRpcRetrieveData(t *testing.T) {
	type args struct {
		ctx     context.Context
		logger  LoggerInterface
		db      DBExecutor
		nk      NakamaModuleInterface
		reader  io.Reader
		payload PayloadRequest
	}
	mockLogger := new(MockLogger)
	nakamaMock := new(MockNakamaModule)
	mock := &mockDB{
		execCount: 0,
		execError: nil,
	}
	mockFalty := &mockDB{
		execCount: 0,
		execError: errors.New("Error executing query"),
	}
	brokenReader := &BrokenReader{}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test Case Success",
			args: args{
				ctx:     context.Background(),
				logger:  mockLogger,
				db:      mock,
				nk:      nakamaMock,
				reader:  strings.NewReader("content"),
				payload: PayloadRequest{},
			},
			want:    `{"type":"","version":"","hash":"","content":null}`,
			wantErr: false,
		},
		{
			name: "Test Case DB Error",
			args: args{
				ctx:     context.Background(),
				logger:  mockLogger,
				db:      mockFalty,
				nk:      nakamaMock,
				reader:  strings.NewReader("content"),
				payload: PayloadRequest{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Test Case Reader Error",
			args: args{
				ctx:     context.Background(),
				logger:  mockLogger,
				db:      mockFalty,
				nk:      nakamaMock,
				reader:  brokenReader,
				payload: PayloadRequest{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteRpcRetrieveData(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.nk, tt.args.reader, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RpcRetrieveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RpcRetrieveData() = %v, want %v", got, tt.want)
			}
		})
	}
}
