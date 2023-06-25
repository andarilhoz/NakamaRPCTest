package payload

import "testing"

func TestGenerateResponse(t *testing.T) {
	type args struct {
		request      PayloadRequest
		request_hash string
		content      string
		equalHashes  bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test Case Success Equal Hashes",
			args: args{
				request: PayloadRequest{
					RequestType:    "level",
					RequestVersion: "1.2.3",
				},
				request_hash: "123abc321",
				content:      "Hello World",
				equalHashes:  true,
			},
			want:    `{"type":"level","version":"1.2.3","hash":"123abc321","content":"Hello World"}`,
			wantErr: false,
		},
		{
			name: "Test Case Different Hashes",
			args: args{
				request: PayloadRequest{
					RequestType:    "level",
					RequestVersion: "1.2.3",
				},
				request_hash: "123abc321",
				content:      "Hello World",
				equalHashes:  false,
			},
			want:    `{"type":"level","version":"1.2.3","hash":"123abc321","content":null}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateResponse(tt.args.request, tt.args.request_hash, tt.args.content, tt.args.equalHashes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
