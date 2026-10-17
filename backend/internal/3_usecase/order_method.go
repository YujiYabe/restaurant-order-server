package usecase

import (
	"context"
	"fmt"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// ResetPlaceList は、Redisの席情報を初期設定へ戻し、メモリと接続画面へ反映します。
func (receiver *UseCase) ResetPlaceList(
	ctx context.Context,
) (
	err error,
) {
	receiver.orderMu.Lock()
	defer receiver.orderMu.Unlock()

	err = receiver.ToGateway.ResetPlaceListInRedis(
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

	receiver.ToDomain.SetPlaceList(
		ctx,
		placeList,
	)

	receiver.DistributeAgent(ctx)

	return
}

// GetPlaceList は、現在の席一覧のスナップショットを返します。
func (receiver *UseCase) GetPlaceList(
	ctx context.Context,
) (
	placeList domain.PlaceList,
) {
	placeList = receiver.ToDomain.GetPlaceList(
		ctx,
	)

	return
}

// GetPlace は、指定された席のスナップショットを返します。
func (receiver *UseCase) GetPlace(
	ctx context.Context,
	placeName string,
) (
	place domain.Place,
) {
	targetPlace := receiver.ToDomain.GetPlace(
		ctx,
		placeName,
	)
	if targetPlace == nil {
		return
	}

	place = *targetPlace

	return
}

// EditPlace は、席の状態遷移を検証し、必要な会計処理を実行して状態を更新します。
func (receiver *UseCase) EditPlace(
	ctx context.Context,
	newPlace domain.Place,
) (
	err error,
) {
	receiver.orderMu.Lock()
	defer receiver.orderMu.Unlock()

	nowPlace := receiver.ToDomain.GetPlace(
		ctx,
		newPlace.PlaceName,
	)
	if nowPlace == nil {
		err = fmt.Errorf("place %q not found", newPlace.PlaceName)
		return
	}

	isNowStatusUsing := nowPlace.PlaceStatus == domain.PlaceStatusUsing
	isNewStatusBilling := newPlace.PlaceStatus == domain.PlaceStatusBilling

	// 支払い処理
	if isNowStatusUsing && isNewStatusBilling {
		err = receiver.operateBilling(
			ctx,
			nowPlace,
		)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
		newPlace.TransactionHeadId = nowPlace.TransactionHeadId
		newPlace.BillingStage = nowPlace.BillingStage
	}

	if !receiver.ToDomain.EditPlace(
		ctx,
		newPlace,
	) {
		return fmt.Errorf("invalid place status transition: place=%s from=%d to=%d", newPlace.PlaceName, nowPlace.PlaceStatus, newPlace.PlaceStatus)
	}

	receiver.DistributeAgent(ctx)

	return
}

// operateBilling は、外部連携の進捗を記録しながら会計処理を段階的に実行します。
func (receiver *UseCase) operateBilling(
	ctx context.Context,
	nowPlace *domain.Place,
) (
	err error,
) {
	// smaregiに登録
	miniProductMap := receiver.ToDomain.GetMiniProductMap(
		ctx,
	)

	sendPostTransactionTemporary, janCodeCountMap := receiver.ToDomain.BuildPostTransactionTemporary(
		ctx,
		nowPlace.PlaceName,
		miniProductMap,
	)

	if nowPlace.BillingStage < domain.BillingStageTransactionCreated {
		var response domain.SmaregiPostResponseEnquiry
		response, err = receiver.ToGateway.SendPostTransactionTemporary(ctx, sendPostTransactionTemporary)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
		if nowPlace.CustomerId != response.CustomerId {
			err = fmt.Errorf("customer ID mismatch: place=%d, response=%d", nowPlace.CustomerId, response.CustomerId)
			return
		}
		if response.TransactionHeadId == "" {
			err = fmt.Errorf("transaction head ID is empty")
			return
		}
		nowPlace.TransactionHeadId = response.TransactionHeadId
		if err = receiver.saveBillingProgress(ctx, nowPlace, domain.BillingStageTransactionCreated); err != nil {
			return
		}
	}
	if nowPlace.TransactionHeadId == "" {
		return fmt.Errorf("transaction head ID is empty at billing stage %d", nowPlace.BillingStage)
	}

	billingPlace := *nowPlace
	billingPlace.TransactionHeadId = nowPlace.TransactionHeadId

	// firestoreに登録
	if nowPlace.BillingStage < domain.BillingStageOrderSaved {
		err = receiver.ToGateway.SaveOrder(ctx, billingPlace, sendPostTransactionTemporary)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
		if err = receiver.saveBillingProgress(ctx, nowPlace, domain.BillingStageOrderSaved); err != nil {
			return
		}
	}

	// 取引ヘッダーID：店舗ID
	if nowPlace.BillingStage < domain.BillingStageTransactionInfoSent {
		err = receiver.ToGateway.SendPostTransactionInfo(ctx, nowPlace.TransactionHeadId)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
		if err = receiver.saveBillingProgress(ctx, nowPlace, domain.BillingStageTransactionInfoSent); err != nil {
			return
		}
	}

	// DB更新
	if nowPlace.BillingStage < domain.BillingStageStockUpdated {
		err = receiver.ToGateway.UpdateStock(ctx, janCodeCountMap)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
		if err = receiver.saveBillingProgress(ctx, nowPlace, domain.BillingStageStockUpdated); err != nil {
			return
		}
	}

	return
}

// saveBillingProgress は、会計処理の進捗をメモリとRedisへ保存します。
func (receiver *UseCase) saveBillingProgress(
	ctx context.Context,
	place *domain.Place,
	stage domain.BillingStage,
) (
	err error,
) {
	place.BillingStage = stage
	receiver.ToDomain.UpdateBillingProgress(ctx, place.PlaceName, place.TransactionHeadId, stage)
	if err := receiver.ToGateway.SetPlaceListToRedis(ctx, receiver.ToDomain.GetPlaceList(ctx)); err != nil {
		return fmt.Errorf("save billing progress: %w", err)
	}

	return nil
}

// AddOrder は、席状態と在庫を検証して注文を追加します。
func (receiver *UseCase) AddOrder(
	ctx context.Context,
	placeName string,
	order domain.Order,
) (
	ok bool,
) {
	receiver.orderMu.Lock()
	defer receiver.orderMu.Unlock()

	// 0個以下の注文は在庫と金額を不正に変更するため受け付けない
	if order.Count <= 0 {
		return
	}

	// 席状態の確認、在庫調整、注文追加を不可分に実行する
	ok = receiver.ToDomain.TryAddOrder(
		ctx,
		placeName,
		order,
	)
	if !ok {
		return
	}

	// オーダーリストの更新を通知
	receiver.DistributeAgent(ctx)

	return
}

// IsAlcohol は、注文商品が酒類として登録されているかを返します。
func (receiver *UseCase) IsAlcohol(
	ctx context.Context,
	order domain.Order,
) (
	ok bool,
) {
	ok = receiver.ToDomain.IsAlcohol(
		ctx,
		order,
	)

	return
}

// IsIncludeTwentyList は、利用者UUIDが有効な成人確認済み一覧に含まれるかを返します。
func (receiver *UseCase) IsIncludeTwentyList(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	ok bool,
) {
	ok = receiver.ToDomain.IsIncludeTwentyList(
		ctx,
		alcohol,
	)

	return
}

// GetPermission は、指定席のインターネット注文用権限キーを返します。
func (receiver *UseCase) GetPermission(
	ctx context.Context,
	placeName string,
) (
	permission string,
) {
	permission = receiver.ToDomain.GetPermission(
		ctx,
		placeName,
	)

	return
}

// EditOrder は、注文番号と状態遷移を検証して注文状態を更新します。
func (receiver *UseCase) EditOrder(
	ctx context.Context,
	placeName string,
	order domain.Order,
) (
	ok bool,
) {
	receiver.orderMu.Lock()
	defer receiver.orderMu.Unlock()

	ok = receiver.ToDomain.EditOrder(
		ctx,
		placeName,
		order,
		receiver.ToDomain.GetMiniProductMap(ctx),
	)

	if ok {
		receiver.DistributeAgent(ctx)
	}

	return
}
