package ioutil

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type faultyReader struct{}

func (r faultyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Unable to Read File")
}

func TestReadFileFromDisk(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test Case Success",
			args: args{
				reader: strings.NewReader("content"),
			},
			want:    "content",
			wantErr: false,
		},
		{
			name: "Test Case Failure",
			args: args{
				reader: faultyReader{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileFromDisk(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFileFromDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadFileFromDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}
