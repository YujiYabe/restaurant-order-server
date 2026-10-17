package gateway

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// VerifyCertificate は、本人確認書類画像をAzureへ送信し、生年月日と書類番号を取得します。
func (receiver *Gateway) VerifyCertificate(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	birthday string,
	number string,
	err error,
) {
	birthday, number, err = receiver.ToAzure.VerifyCertificate(
		ctx,
		alcohol,
	)
	return
}
