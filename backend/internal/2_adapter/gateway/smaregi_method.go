package gateway

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// SendPostTransactionTemporary は、注文内容をスマレジへ仮販売登録し、取引情報を返します。
func (receiver *Gateway) SendPostTransactionTemporary(
	ctx context.Context,
	reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
) (
	smaregiPostResponseEnquiry domain.SmaregiPostResponseEnquiry,
	err error,
) {
	smaregiPostResponseEnquiry, err = receiver.ToSmaregi.SendPostTransactionTemporary(
		ctx,
		reqPostTransactionTemporary,
	)

	return
}
