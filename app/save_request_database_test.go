package main

import (
	"context"
	"errors"
	"testing"
)

func TestSaveRequestInDatabase(t *testing.T) {
	type args struct {
		db             DBExecutorInterface
		ctx            context.Context
		request        PayloadRequest
		hashesAreEqual bool
	}
	mock := &mockDB{
		execCount: 0,
		execError: nil,
	}
	mockFalty := &mockDB{
		execCount: 0,
		execError: errors.New("Error executing query"),
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		execCount int
	}{
		{
			name: "Test Case Success",
			args: args{
				db:  mock,
				ctx: context.Background(),
				request: PayloadRequest{
					RequestType:    "level",
					RequestVersion: "1.2.3",
					RequestHash:    &[]string{"123abc321"}[0],
				},
				hashesAreEqual: true,
			},
			wantErr:   false,
			execCount: 1,
		},
		{
			name: "Test Case Error",
			args: args{
				db:  mockFalty,
				ctx: context.Background(),
				request: PayloadRequest{
					RequestType:    "level",
					RequestVersion: "1.2.3",
					RequestHash:    &[]string{"123abc321"}[0],
				},
				hashesAreEqual: true,
			},
			wantErr:   true,
			execCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveRequestInDatabase(tt.args.ctx, tt.args.db, tt.args.request, tt.args.hashesAreEqual); (err != nil) != tt.wantErr {
				t.Errorf("SaveRequestInDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
			if mock.execCount != tt.execCount {
				t.Errorf("DB.Exec Times = %v, expected %v", mock.execCount, tt.execCount)
			}
		})
	}
}
