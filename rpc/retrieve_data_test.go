package rpc

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"
	"testing"

	_db "heroiclabs.com/go-setup-demo/db"
	_pl "heroiclabs.com/go-setup-demo/payload"
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

type mockDB struct {
	execCount int
	execError error
}

func (db *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	db.execCount++
	return nil, db.execError
}

func TestRpcRetrieveData(t *testing.T) {
	type args struct {
		ctx     context.Context
		logger  LoggerInterface
		db      _db.DBExecutorInterface
		nk      NakamaModuleInterface
		reader  io.Reader
		payload _pl.PayloadRequest
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
				payload: _pl.PayloadRequest{},
			},
			want:    `{"type":"","version":"","hash":"","content":null}`,
			wantErr: false,
		},
		{
			name: "Test Case Success Same Hash",
			args: args{
				ctx:    context.Background(),
				logger: mockLogger,
				db:     mock,
				nk:     nakamaMock,
				reader: strings.NewReader("Hello World"),
				payload: _pl.PayloadRequest{
					RequestType:    "core",
					RequestVersion: "1.0.0",
					RequestHash:    &[]string{"a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"}[0],
				},
			},
			want:    `{"type":"core","version":"1.0.0","hash":"a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e","content":"Hello World"}`,
			wantErr: false,
		},
		{
			name: "Test Case Success Diff Hash",
			args: args{
				ctx:    context.Background(),
				logger: mockLogger,
				db:     mock,
				nk:     nakamaMock,
				reader: strings.NewReader("Hello World"),
				payload: _pl.PayloadRequest{
					RequestType:    "core",
					RequestVersion: "1.0.0",
					RequestHash:    &[]string{"123abc321"}[0],
				},
			},
			want:    `{"type":"core","version":"1.0.0","hash":"123abc321","content":null}`,
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
				payload: _pl.PayloadRequest{},
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
				payload: _pl.PayloadRequest{},
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
