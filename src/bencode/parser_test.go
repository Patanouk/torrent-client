package bencode

import (
	"bufio"
	"strings"
	"testing"
)

func Test_decodeString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Input empty", args{input: ""}, "", true},
		{"Input no separator", args{input: "aa"}, "", true},
		{"Input integer only", args{input: "4"}, "", true},
		{"Input empty string after semicolon", args{input: "4:"}, "", true},
		{"Input length not integer", args{input: "-aa:"}, "", true},
		{"Input negative length", args{input: "-4:"}, "", true},
		{"Input length too long", args{input: "4:aa"}, "", true},

		{"Input zero length", args{input: "0:"}, "", false},
		{"Input one length", args{input: "1:a"}, "a", false},
		{"Input seven length", args{input: "7:abcdefg"}, "abcdefg", false},
		//Invalid format globally, but valid in the scope of the function
		{"Input length too short", args{input: "3:abcj"}, "abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeString(bufio.NewReader(strings.NewReader(tt.args.input)))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeInteger(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"Input empty", args{input: ""}, 0, true},
		{"Input integer only", args{input: "10"}, 0, true},
		{"Input no end char", args{input: "i10"}, 0, true},
		{"Input no start char", args{input: "10e"}, 0, true},
		{"Input invalid integer", args{input: "10e"}, 0, true},

		{"Input 0", args{input: "i0e"}, 0, false},
		{"Input 10", args{input: "i10e"}, 10, false},
		{"Input 99", args{input: "i99e"}, 99, false},
		{"Input -1", args{input: "i-1e"}, -1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeInteger(bufio.NewReader(strings.NewReader(tt.args.input)))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeInteger() got = %v, want %v", got, tt.want)
			}
		})
	}
}
