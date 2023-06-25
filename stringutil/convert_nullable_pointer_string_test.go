package stringutil

import "testing"

func TestConvertNullablePointerToString(t *testing.T) {
	type args struct {
		pointer *string
	}
	var testPointerNil *string
	var testPointerFilled *string = &[]string{"testing"}[0]
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case Nil",
			args: args{
				pointer: testPointerNil,
			},
			want: "",
		},
		{
			name: "Test Case Filled",
			args: args{
				pointer: testPointerFilled,
			},
			want: "testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertNullablePointerToString(tt.args.pointer); got != tt.want {
				t.Errorf("ConvertNullablePointerToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
