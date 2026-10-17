package gateway

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// GetPlaceListFromRedis は、Redisに保存された席一覧を取得します。
func (receiver *Gateway) GetPlaceListFromRedis(
	ctx context.Context,
) (
	placeList domain.PlaceList,
	err error,
) {
	placeList, err = receiver.ToRedis.GetPlaceListFromRedis(ctx)
	return
}

// SetPlaceListToRedis は、席一覧をRedisへ永続化します。
func (receiver *Gateway) SetPlaceListToRedis(
	ctx context.Context,
	placeList domain.PlaceList,
) (
	err error,
) {
	err = receiver.ToRedis.SetPlaceListToRedis(
		ctx,
		placeList,
	)
	return
}

// ResetPlaceListInRedis は、Redisに保存された席一覧を削除します。
func (receiver *Gateway) ResetPlaceListInRedis(
	ctx context.Context,
) (
	err error,
) {
	err = receiver.ToRedis.ResetPlaceListInRedis(
		ctx,
	)

	return
}
