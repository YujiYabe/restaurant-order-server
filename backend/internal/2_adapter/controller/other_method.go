package controller

import (
	"context"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// Start は、必要な初期処理とバックグラウンド処理を開始します。
func (receiver *controller) Start() {
	ctx := context.Background()

	receiver.ToUseCase.Start(ctx) // 初期処理 DBをインメモリに保存
}

// アレルギー -----------------------
// GetAllergyDefault は、既定のアレルギー設定をユースケースから取得します。
func (receiver *controller) GetAllergyDefault(
	ctx context.Context,
) (
	allergy domain.Allergy,
) {
	allergy = receiver.ToUseCase.GetAllergyDefault(
		ctx,
	)

	return
}

// SSE配信 -----------------------
// DistributeAgent は、現在の状態を接続中の画面へ配信するようユースケースへ依頼します。
func (receiver *controller) DistributeAgent(
	ctx context.Context,
) {
	receiver.ToUseCase.DistributeAgent(
		ctx,
	)
}

// GetTagList は、商品タグ一覧のコピーを返します。
func (receiver *controller) GetTagList(
	ctx context.Context,
) (
	tagPrimaryList domain.TagPrimaryList,
) {
	tagPrimaryList = receiver.ToUseCase.GetTagList(
		ctx,
	)

	return
}

// GetInventoryMap は、現在の在庫数マップのコピーを返します。
func (receiver *controller) GetInventoryMap(
	ctx context.Context,
) (
	inventoryMap domain.InventoryMap,
) {
	inventoryMap = receiver.ToUseCase.GetInventoryMap(
		ctx,
	)

	return
}

// VerifyTwenty は、書類から取得した生年月日を検証し、成人であることを確認します。
func (receiver *controller) VerifyTwenty(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	certificateID string,
	age int,
	birthday time.Time,
	err error,
) {
	certificateID, age, birthday, err = receiver.ToUseCase.VerifyTwenty(
		ctx,
		alcohol,
	)

	return
}

// RegisterFace は、許可済み利用者の顔画像をFirestoreへ登録します。
func (receiver *controller) RegisterFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {
	err = receiver.ToUseCase.RegisterFace(
		ctx,
		alcohol,
	)

	return
}

// VerifyFace は、登録済み画像との顔照合に成功した利用者へ当日の酒類注文を許可します。
func (receiver *controller) VerifyFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {
	err = receiver.ToUseCase.VerifyFace(
		ctx,
		alcohol,
	)

	return
}
