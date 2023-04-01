package action

import (
	"reflect"
	"testing"
)

func TestTruncateString(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Truncate string",
			args: args{
				s:      "This is a string",
				length: 10,
			},
			want: "This is a ...",
		},
		{
			name: "Truncate string",
			args: args{
				s:      "This is a string",
				length: 20,
			},
			want: "This is a string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("TruncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "Remove duplicates",
			args: args{
				elements: []string{"a", "b", "c", "a", "b", "c"},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "Remove duplicates",
			args: args{
				elements: []string{"a", "b", "c", "a", "b", "c", "d", "e", "f", "d", "e", "f"},
			},
			want: []string{"a", "b", "c", "d", "e", "f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedLengthString(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Fixed length string",
			args: args{
				s:      "This is a string",
				length: 10,
			},
			want: "This is a ",
		},
		{
			name: "Fixed length string",
			args: args{
				s:      "This is a string",
				length: 20,
			},
			want: "This is a string    ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FixedLengthString(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("FixedLengthString() = %v, want %v", got, tt.want)
			}
		})
	}
}
