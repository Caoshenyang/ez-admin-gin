package middleware

import "testing"

func TestTruncateOperationLogText(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		maxLength int
		want      string
	}{
		{name: "trim whitespace", value: "  created user  ", maxLength: 20, want: "created user"},
		{name: "truncate ascii", value: "abcdef", maxLength: 3, want: "abc"},
		{name: "truncate utf8 by rune", value: "中文错误信息", maxLength: 4, want: "中文错误"},
		{name: "non positive length", value: "anything", maxLength: 0, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateOperationLogText(tt.value, tt.maxLength); got != tt.want {
				t.Fatalf("truncateOperationLogText(%q, %d) = %q, want %q", tt.value, tt.maxLength, got, tt.want)
			}
		})
	}
}
