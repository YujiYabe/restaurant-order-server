package pkg

import "testing"

// TestDecodeBase64DataURL は、対象仕様の回帰動作を検証します。
func TestDecodeBase64DataURL(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{name: "valid", source: "data:image/jpeg;base64,aW1hZ2U="},
		{name: "missing delimiter", source: "aW1hZ2U=", wantErr: true},
		{name: "empty payload", source: "data:image/jpeg;base64,", wantErr: true},
		{name: "invalid payload", source: "data:image/jpeg;base64,%%%", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBase64DataURL(tt.source)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DecodeBase64DataURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
