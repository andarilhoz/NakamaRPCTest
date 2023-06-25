package ioutil

import (
	"testing"

	_pl "heroiclabs.com/go-setup-demo/payload"
)

func TestGetFilePath(t *testing.T) {
	type args struct {
		request _pl.PayloadRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case Filled",
			args: args{
				_pl.PayloadRequest{
					RequestType:    "level",
					RequestVersion: "1.2.3",
				},
			},
			want: "/nakama/json_test_files/level/1.2.3.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilePath(tt.args.request); got != tt.want {
				t.Errorf("GetFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
