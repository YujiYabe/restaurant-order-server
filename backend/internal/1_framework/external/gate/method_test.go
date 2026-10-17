package gate

import (
	"net/http"
	"testing"
)

func TestTransactionInfoURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseURL string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "configured endpoint",
			baseURL: "https://example.com",
			path:    "/transactions",
			want:    "https://example.com/transactions",
		},
		{name: "missing base URL", path: "/transactions", wantErr: true},
		{name: "missing path", baseURL: "https://example.com", wantErr: true},
		{name: "relative base URL", baseURL: "example.com", path: "/transactions", wantErr: true},
		{name: "absolute path URL", baseURL: "https://example.com", path: "https://other.example.com/transactions", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := transactionInfoURL(test.baseURL, test.path)
			if test.wantErr {
				if err == nil {
					t.Fatalf("transactionInfoURL() returned no error, URL = %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("transactionInfoURL() error = %v", err)
			}
			if got.String() != test.want {
				t.Fatalf("transactionInfoURL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestIsSuccessfulStatus(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent} {
		if !isSuccessfulStatus(status) {
			t.Fatalf("isSuccessfulStatus(%d) = false", status)
		}
	}
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		if isSuccessfulStatus(status) {
			t.Fatalf("isSuccessfulStatus(%d) = true", status)
		}
	}
}

func TestShouldSkipTransactionInfo(t *testing.T) {
	t.Parallel()
	if !shouldSkipTransactionInfo("local") {
		t.Error("local transaction info must not be sent to an external service")
	}
	if shouldSkipTransactionInfo("prod") {
		t.Error("production transaction info must be sent")
	}
}
