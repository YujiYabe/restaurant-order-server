package controller

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// 注文 -----------------------
// ResetPlaceList は、席一覧の初期化をユースケースへ依頼します。
func (receiver *controller) ResetPlaceList(
	ctx context.Context,
) (
	err error,
) {
	err = receiver.ToUseCase.ResetPlaceList(
		ctx,
	)

	return
}

// GetPlaceList は、現在の席一覧のスナップショットを返します。
func (receiver *controller) GetPlaceList(
	ctx context.Context,
) (
	placeList domain.PlaceList,
) {
	placeList = receiver.ToUseCase.GetPlaceList(
		ctx,
	)

	return
}

// GetPlace は、指定された席のスナップショットを返します。
func (receiver *controller) GetPlace(
	ctx context.Context,
	placeName string,
) (
	place domain.Place,
) {
	place = receiver.ToUseCase.GetPlace(
		ctx,
		placeName,
	)

	return
}

// EditPlace は、席の状態遷移を検証し、必要な会計処理を実行して状態を更新します。
func (receiver *controller) EditPlace(
	ctx context.Context,
	place domain.Place,
) (
	err error,
) {
	err = receiver.ToUseCase.EditPlace(
		ctx,
		place,
	)

	return
}

// IsAlcohol は、注文商品が酒類として登録されているかを返します。
func (receiver *controller) IsAlcohol(
	ctx context.Context,
	order domain.Order,
) (
	ok bool,
) {
	ok = receiver.ToUseCase.IsAlcohol(
		ctx,
		order,
	)

	return
}

// IsIncludeTwentyList は、利用者UUIDが有効な成人確認済み一覧に含まれるかを返します。
func (receiver *controller) IsIncludeTwentyList(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	ok bool,
) {
	ok = receiver.ToUseCase.IsIncludeTwentyList(
		ctx,
		alcohol,
	)

	return
}

// AddOrder は、席状態と在庫を検証して注文を追加します。
func (receiver *controller) AddOrder(
	ctx context.Context,
	placeName string,
	order domain.Order,
) (
	ok bool,
) {
	ok = receiver.ToUseCase.AddOrder(
		ctx,
		placeName,
		order,
	)

	return
}

// GetPermission は、指定席のインターネット注文用権限キーを返します。
func (receiver *controller) GetPermission(
	ctx context.Context,
	placeName string,
) (
	permission string,
) {
	permission = receiver.ToUseCase.GetPermission(
		ctx,
		placeName,
	)

	return
}

// EditOrder は、注文番号と状態遷移を検証して注文状態を更新します。
func (receiver *controller) EditOrder(
	ctx context.Context,
	placeName string,
	order domain.Order,
) (
	ok bool,
) {
	ok = receiver.ToUseCase.EditOrder(
		ctx,
		placeName,
		order,
	)

	return
}
