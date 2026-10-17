package gateway

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// SaveOrder は、会計対象の注文情報をFirestoreへ保存します。
func (receiver *Gateway) SaveOrder(
	ctx context.Context,
	place domain.Place,
	reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
) (
	err error,
) {
	err = receiver.ToFirestore.SaveOrder(
		ctx,
		place,
		reqPostTransactionTemporary,
	)

	return
}

// SaveClientAttribute は、本人確認書類の属性を利用者UUIDに紐づけてFirestoreへ保存します。
func (receiver *Gateway) SaveClientAttribute(
	ctx context.Context,
	alcohol domain.Alcohol,
	certificateID string,
) (
	err error,
) {
	err = receiver.ToFirestore.SaveClientAttribute(
		ctx,
		alcohol,
		certificateID,
	)

	return
}

// GetClientImage は、利用者UUIDに紐づく最新の顔画像をFirestoreから取得します。
func (receiver *Gateway) GetClientImage(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	imageList []string,
	err error,
) {
	imageList, err = receiver.ToFirestore.GetClientImage(
		ctx,
		alcohol,
	)

	return
}

// SaveClientOrder は、利用者の酒類注文履歴をFirestoreへ保存します。
func (receiver *Gateway) SaveClientOrder(
	ctx context.Context,
	clientUUID string,
	imageUUID string,
	janCode int,
	count int,
) (
	err error,
) {
	err = receiver.ToFirestore.SaveClientOrder(
		ctx,
		clientUUID,
		imageUUID,
		janCode,
		count,
	)

	return
}
