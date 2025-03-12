package tests

import (
	"net/http"
	"testing"

	"github.com/divakaivan/rssagg/internal/auth"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		want      string
		expectErr bool
	}{
		{
			name:      "No Authorization header",
			headers:   http.Header{},
			want:      "",
			expectErr: true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"Bearer token"},
			},
			want:      "",
			expectErr: true,
		},
		{
			name: "Invalid scheme in Authorization header",
			headers: http.Header{
				"Authorization": []string{"Bearer token"},
			},
			want:      "",
			expectErr: true,
		},
		{
			name: "Valid Authorization header",
			headers: http.Header{
				"Authorization": []string{"ApiKey validapikey"},
			},
			want:      "validapikey",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.GetAPIKey(tt.headers)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAPIKey() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
