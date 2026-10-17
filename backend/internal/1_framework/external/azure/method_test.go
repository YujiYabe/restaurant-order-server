package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

func TestVerifyCertificateRejectsMissingAnalyzeURL(t *testing.T) {
	originalURL := env.Env.AzureDocumentIntelligenceURL
	env.Env.AzureDocumentIntelligenceURL = ""
	t.Cleanup(func() {
		env.Env.AzureDocumentIntelligenceURL = originalURL
	})

	receiver := &Azure{}
	if _, _, err := receiver.VerifyCertificate(context.Background(), domain.Alcohol{}); err == nil {
		t.Fatal("VerifyCertificate() returned no error for an empty analyze URL")
	}
}

// TestExtractCertificateFieldsRejectsMissingData は、対象仕様の回帰動作を検証します。
func TestExtractCertificateFieldsRejectsMissingData(t *testing.T) {
	tests := []struct {
		name   string
		result Result
	}{
		{name: "no documents", result: Result{}},
		{
			name: "no birthday",
			result: Result{AnalyzeResult: AnalyzeResult{Documents: []Document{
				{Fields: DocumentFields{Number: FieldDetail{ValueString: "123"}}},
			}}},
		},
		{
			name: "no certificate number",
			result: Result{AnalyzeResult: AnalyzeResult{Documents: []Document{
				{Fields: DocumentFields{Birthday: FieldDetail{ValueString: "昭和54年9月4日"}}},
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := extractCertificateFields(tt.result); err == nil {
				t.Fatal("extractCertificateFields() returned no error")
			}
		})
	}
}

func TestPollCertificateResultReturnsSuccessfulAnalysis(t *testing.T) {
	wantBirthday := "昭和54年9月4日"
	wantNumber := "123456"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(Result{
			Status: "succeeded",
			AnalyzeResult: AnalyzeResult{Documents: []Document{{Fields: DocumentFields{
				Birthday: FieldDetail{ValueString: wantBirthday},
				Number:   FieldDetail{ValueString: wantNumber},
			}}}},
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()

	receiver := &Azure{}
	birthday, number, err := receiver.pollCertificateResult(
		context.Background(),
		server.URL,
		time.Nanosecond,
	)
	if err != nil {
		t.Fatalf("pollCertificateResult() error = %v", err)
	}
	if birthday != wantBirthday || number != wantNumber {
		t.Fatalf("pollCertificateResult() = (%q, %q), want (%q, %q)", birthday, number, wantBirthday, wantNumber)
	}
}
