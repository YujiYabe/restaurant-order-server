package pkg

import (
	"context"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

// GetNewContext は、親コンテキストへリクエストIDを設定した子コンテキストを返します。
func GetNewContext(
	ctx context.Context,
	headerXRequestID string,
) (
	newCtx context.Context,
) {
	// echoのX-Request-IDをcontextに埋め込み
	newCtx = context.WithValue(
		ctx,              // 親コンテキスト
		RequestIDKey,     // キー
		headerXRequestID, // 値
	)

	return
}
