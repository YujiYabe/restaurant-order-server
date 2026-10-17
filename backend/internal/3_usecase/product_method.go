package usecase

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// GetProduct は、指定JANコードの商品を返します。
func (receiver *UseCase) GetProduct(
	ctx context.Context,
	productNumber int,
) (
	product domain.Product,
) {
	product = receiver.ToDomain.GetProduct(
		ctx,
		productNumber,
	)

	return
}

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *UseCase) GetAllProductList(
	ctx context.Context,
) (
	allProductList domain.AllProductList,
) {
	allProductList = receiver.ToDomain.GetAllProductList(
		ctx,
	)

	return
}

// GetMiniProductMap は、注文計算に使用する商品概要マップのコピーを返します。
func (receiver *UseCase) GetMiniProductMap(
	ctx context.Context,
) (
	miniProductMap domain.MiniProductMap,
) {
	miniProductMap = receiver.ToDomain.GetMiniProductMap(
		ctx,
	)

	return
}

// UpdateProduct は、商品情報を更新し、メモリ上のマスター情報を再構築します。
func (receiver *UseCase) UpdateProduct(
	ctx context.Context,
	product domain.Product,
) (
	err error,
) {
	// DB更新
	err = receiver.ToGateway.UpdateProduct(
		ctx,
		product,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// インメモリの情報を更新
	receiver.SetUpInMemory(ctx)

	// 在庫の更新を通知
	receiver.DistributeAgent(ctx)

	return
}

// SetUpInMemory は、SQLiteとRedisからデータを読み込み、ドメインのインメモリ状態を構築します。
func (receiver *UseCase) SetUpInMemory(
	ctx context.Context,
) {
	// localDBから全商品を取得
	allProductList, err := receiver.ToGateway.GetAllProductList(
		ctx,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// localDBから全tagを取得
	allTag, err := receiver.ToGateway.GetAllTag(
		ctx,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}
	// redisから取得
	placeList, err := receiver.ToGateway.GetPlaceListFromRedis(
		ctx,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	receiver.ToDomain.SaveInMemory(
		ctx,
		allProductList,
		allTag,
		placeList,
	)

}
