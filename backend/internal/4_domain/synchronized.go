package domain

import (
	"context"
	"time"
)

// GetPlaceList は、現在の席一覧のスナップショットを返します。
func (receiver *Domain) GetPlaceList(
	context.Context,
) PlaceList {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return clonePlaceList(receiver.PlaceList)
}

// GetPlace は、指定された席のスナップショットを返します。
func (receiver *Domain) GetPlace(
	_ context.Context,
	placeName string,
) *Place {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	place := receiver.PlaceList[placeName]
	if place == nil {
		return nil
	}
	return clonePlace(place)
}

// EditPlace は、席一覧をロックして状態遷移ルールを適用します。
// 遷移が許可されて席情報を更新した場合だけtrueを返します。外部会計連携はこのメソッドでは実行しません。
func (receiver *Domain) EditPlace(
	ctx context.Context,
	place Place,
) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	return receiver.PlaceList.EditPlace(ctx, place)
}

// UpdateBillingProgress は、指定席の取引IDと会計処理段階を不可分に更新します。
// 席が存在しない場合はfalseを返し、状態を変更しません。
func (receiver *Domain) UpdateBillingProgress(
	_ context.Context,
	placeName string,
	transactionHeadID string,
	stage BillingStage,
) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	place := receiver.PlaceList[placeName]
	if place == nil {
		return false
	}
	place.TransactionHeadId = transactionHeadID
	place.BillingStage = stage
	return true
}

// AddOrder は、指定席へ注文を追加し、席の注文金額を再計算します。
// このメソッドは在庫と席状態を検証しないため、通常の注文受付ではTryAddOrderを使用します。
func (receiver *Domain) AddOrder(
	ctx context.Context,
	placeName string,
	order Order,
	miniProductMap MiniProductMap,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	receiver.PlaceList.AddOrder(ctx, placeName, order, miniProductMap)
}

// TryAddOrder は、席状態と在庫をロック下で検証し、在庫減算と注文追加を不可分に実行します。
// Using中の席、正数の数量、0以外の登録済みJANコード、および十分な在庫数が
// すべて満たされた場合だけtrueを返します。JANコード0は内部生成するチャージ明細専用です。
func (receiver *Domain) TryAddOrder(
	ctx context.Context,
	placeName string,
	order Order,
) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	place := receiver.PlaceList[placeName]
	if place == nil || place.PlaceStatus != PlaceStatusUsing || order.Count <= 0 || order.JANCode == 0 {
		return false
	}
	stock, exists := receiver.InventoryMap[order.JANCode]
	if !exists || stock < order.Count {
		return false
	}
	receiver.InventoryMap[order.JANCode] = stock - order.Count

	receiver.PlaceList.AddOrder(ctx, placeName, order, receiver.MiniProductMap)
	return true
}

// GetPermission は、指定席のインターネット注文用権限キーを返します。
func (receiver *Domain) GetPermission(_ context.Context, placeName string) string {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	place := receiver.PlaceList[placeName]
	if place == nil {
		return ""
	}
	return place.Permission
}

// EditOrder は、注文番号と状態遷移を検証して注文状態を更新します。
// 取消への変更時は在庫を復元し、取消から戻す場合は在庫を再確保します。
// 存在しない席・注文、または不正な状態遷移の場合はfalseを返します。
func (receiver *Domain) EditOrder(
	ctx context.Context,
	placeName string,
	order Order,
	miniProductMap MiniProductMap,
) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	place := receiver.PlaceList[placeName]
	if place == nil {
		return false
	}
	currentOrder, exists := place.OrderList[order.OrderNumber]
	if !exists {
		return false
	}
	if currentOrder.OrderStatus == OrderStatusCancel && order.OrderStatus != OrderStatusCancel && currentOrder.JANCode != 0 {
		stock, stockExists := receiver.InventoryMap[currentOrder.JANCode]
		if !stockExists || stock < currentOrder.Count {
			return false
		}
		receiver.InventoryMap[currentOrder.JANCode] = stock - currentOrder.Count
	}

	updated := receiver.PlaceList.EditOrder(ctx, placeName, order, miniProductMap)
	if !updated && currentOrder.OrderStatus == OrderStatusCancel && order.OrderStatus != OrderStatusCancel && currentOrder.JANCode != 0 {
		receiver.InventoryMap[currentOrder.JANCode] += currentOrder.Count
	}
	if updated && currentOrder.OrderStatus != OrderStatusCancel && order.OrderStatus == OrderStatusCancel && currentOrder.JANCode != 0 {
		receiver.InventoryMap[currentOrder.JANCode] += currentOrder.Count
	}
	return updated
}

// BuildPostTransactionTemporary は、席の未取消注文からスマレジ仮販売登録データと商品数量を生成します。
func (receiver *Domain) BuildPostTransactionTemporary(
	ctx context.Context,
	placeName string,
	miniProductMap MiniProductMap,
) (ReqPostTransactionTemporary, map[int]int) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return receiver.PlaceList.BuildPostTransactionTemporary(
		ctx,
		placeName,
		miniProductMap,
	)
}

// GetProduct は、指定JANコードの商品を返します。
func (receiver *Domain) GetProduct(
	ctx context.Context,
	janCode int,
) Product {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return receiver.Stock.GetProduct(ctx, janCode)
}

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *Domain) GetAllProductList(context.Context) AllProductList {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return append(AllProductList(nil), receiver.Stock.AllProductList...)
}

// UpdateProduct は、ドメインの商品一覧に含まれる対象商品をロック下で更新します。
func (receiver *Domain) UpdateProduct(
	ctx context.Context,
	product Product,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	receiver.Stock.UpdateProduct(ctx, product)
}

// GetIsValidJANCodeList は、注文可能なJANコード一覧のコピーを返します。
func (receiver *Domain) GetIsValidJANCodeList() []int {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return append([]int(nil), receiver.Stock.IsValidJANCodeList...)
}

// VerifyJANCodeList は、指定JANコードと言語コードを有効な値だけに補正します。
func (receiver *Domain) VerifyJANCodeList(
	janCodeList []int,
	isValidLangCodeList []int,
	defaultLangCode int,
) ([]int, int) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	return receiver.Stock.VerifyJANCodeList(
		janCodeList,
		isValidLangCodeList,
		defaultLangCode,
	)
}

// autoUpdate は、コンテキストが終了するまで席のチャージ料金を定期更新します。
func (receiver *Domain) autoUpdate(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receiver.mutex.Lock()
		isUpdated := false
		for _, place := range receiver.PlaceList {
			if place.AutoUpdate(ctx, receiver.MiniProductMap) {
				isUpdated = true
			}
		}
		receiver.mutex.Unlock()

		if isUpdated {
			select {
			case IsAutoUpdate <- true:
			case <-ctx.Done():
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// clonePlaceList は、席一覧と各注文一覧のディープコピーを生成します。
func clonePlaceList(source PlaceList) PlaceList {
	cloned := make(PlaceList, len(source))
	for name, place := range source {
		cloned[name] = clonePlace(place)
	}
	return cloned
}

// clonePlace は、席と注文一覧のディープコピーを生成します。
func clonePlace(source *Place) *Place {
	if source == nil {
		return nil
	}
	cloned := *source
	cloned.OrderList = make(OrderList, len(source.OrderList))
	for number, order := range source.OrderList {
		cloned.OrderList[number] = order
	}
	return &cloned
}
