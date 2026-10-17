package controller

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// GetProduct は、指定JANコードの商品を返します。
func (receiver *controller) GetProduct(
	ctx context.Context,
	productNumber int,
) (
	product domain.Product,
) {
	product = receiver.ToUseCase.GetProduct(
		ctx,
		productNumber,
	)

	return
}

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *controller) GetAllProductList(
	ctx context.Context,
) (
	allProductList domain.AllProductList,
) {
	allProductList = receiver.ToUseCase.GetAllProductList(
		ctx,
	)

	return
}

// GetMiniProductMap は、注文計算に使用する商品概要マップのコピーを返します。
func (receiver *controller) GetMiniProductMap(
	ctx context.Context,
) (
	miniProductMap domain.MiniProductMap,
) {
	miniProductMap = receiver.ToUseCase.GetMiniProductMap(
		ctx,
	)

	return
}

// UpdateProduct は、商品情報を更新し、メモリ上のマスター情報を再構築します。
func (receiver *controller) UpdateProduct(
	ctx context.Context,
	product domain.Product,
) (
	err error,
) {
	return receiver.ToUseCase.UpdateProduct(
		ctx,
		product,
	)
}
