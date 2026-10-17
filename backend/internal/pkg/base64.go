package pkg

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// DecodeBase64DataURL は、Base64データURLのペイロードを検証してデコードします。
func DecodeBase64DataURL(source string) ([]byte, error) {
	_, payload, found := strings.Cut(source, "base64,")
	if !found || payload == "" {
		return nil, fmt.Errorf("invalid Base64 data URL")
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("decode Base64 data URL: %w", err)
	}
	if len(decoded) == 0 {
		return nil, fmt.Errorf("Base64 data URL payload is empty")
	}

	return decoded, nil
}
