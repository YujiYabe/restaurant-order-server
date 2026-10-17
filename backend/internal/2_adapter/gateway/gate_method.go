package gateway

import (
	"context"
)

// SendPostTransactionInfo は、取引IDと店舗IDを会計連携サービスへ通知します。
func (receiver *Gateway) SendPostTransactionInfo(
	ctx context.Context,
	transactionHeadId string,
) (
	err error,
) {
	err = receiver.ToGate.SendPostTransactionInfo(
		ctx,
		transactionHeadId,
	)

	return
}
