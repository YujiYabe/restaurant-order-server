package domain

import (
	"context"
	"math"
	"strings"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

// GetDomain は、現在のドメイン状態のディープコピーを返します。
// 返されたマップ、スライス、席、または注文を呼び出し側が変更しても、保持中の状態には影響しません。
func (receiver *Domain) GetDomain(
	ctx context.Context,
) (
	domain *Domain,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	domain = receiver.clone()
	return
}

// GetAllergyDefault は、既定のアレルギー設定のコピーを返します。
func (receiver *Domain) GetAllergyDefault(
	ctx context.Context,
) (
	allergy Allergy,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	allergy = cloneAllergy(receiver.AllergyDefault)
	return
}

// GetTagList は、商品タグ一覧のコピーを返します。
func (receiver *Domain) GetTagList(
	ctx context.Context,
) (
	tagPrimaryList TagPrimaryList,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	tagPrimaryList = cloneTagPrimaryList(receiver.TagPrimaryList)
	return
}

// GetInventoryMap は、現在の在庫数マップのコピーを返します。
func (receiver *Domain) GetInventoryMap(
	ctx context.Context,
) (
	inventoryMap InventoryMap,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	inventoryMap = make(InventoryMap, len(receiver.InventoryMap))
	for janCode, count := range receiver.InventoryMap {
		inventoryMap[janCode] = count
	}
	return
}

// GetMiniProductMap は、注文計算に使用する商品概要マップのコピーを返します。
func (receiver *Domain) GetMiniProductMap(
	ctx context.Context,
) (
	miniProductMap MiniProductMap,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	miniProductMap = make(MiniProductMap, len(receiver.MiniProductMap))
	for janCode, product := range receiver.MiniProductMap {
		miniProductMap[janCode] = product
	}
	return
}

// AdjustmentInventory は、正数の注文数と在庫数を検証して在庫を減算します。
// 数量が0以下、または在庫不足の場合は状態を変更せずfalseを返します。
// このメソッドは席への注文追加を行わないため、通常の注文受付ではTryAddOrderを使用します。
func (receiver *Domain) AdjustmentInventory(
	ctx context.Context,
	newOrder Order,
) (
	ok bool,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if newOrder.Count <= 0 {
		return
	}

	// 在庫以上の注文であればエラーとする
	if receiver.InventoryMap[newOrder.JANCode] < newOrder.Count {
		return
	}

	receiver.InventoryMap[newOrder.JANCode] -= newOrder.Count
	ok = true

	return
}

// SaveInMemory は、永続データから商品・在庫・席情報を再構築し、未会計注文を在庫へ反映します。
//
// 有効商品だけをProductListとIsValidJANCodeListへ登録し、すべての商品から在庫と注文計算用の
// MiniProductMapを構築します。酒類タグを含む商品はAlcoholListへ登録されます。
//
// SQLite在庫から、Using中の未取消注文と、Billing中かつBillingStageStockUpdatedより前の
// 未取消注文を差し引きます。計算結果が負数になる場合は0に補正します。
//
// チャージ機能が有効な場合は旧自動更新ループを停止し、渡されたctxを親とする新しいループを開始します。
func (receiver *Domain) SaveInMemory(
	ctx context.Context,
	allProductList AllProductList,
	tagPrimaryList TagPrimaryList,
	placeList PlaceList,
) {
	productList := ProductList{}
	isValidJANCodeList := []int{}
	inventoryMap := map[int]int{}
	miniProductMap := MiniProductMap{}
	alcoholList := []int{}

	for _, allProduct := range allProductList {
		if allProduct.IsValid {
			productList = append(
				productList,
				allProduct,
			)

			isValidJANCodeList = append(
				isValidJANCodeList,
				allProduct.JANCode,
			)

			inventoryMap[allProduct.JANCode] = allProduct.Stock

			if strings.Contains(
				allProduct.Tag,
				"alcohol",
			) {
				alcoholList = append(
					alcoholList,
					allProduct.JANCode,
				)
			}

		} else {
			inventoryMap[allProduct.JANCode] = 0
		}

		price := math.Floor(
			env.Env.EnvUnique.Commission * float64(allProduct.RetailPrice+allProduct.OperationCost+allProduct.EnergyCost),
		)

		miniProductMap[allProduct.JANCode] = MiniProduct{
			JANCode:       allProduct.JANCode,
			ProductId:     allProduct.ProductId,
			Price:         int(price),
			Place:         allProduct.Place,
			Recipe:        allProduct.Recipe,
			Tag:           allProduct.Tag,
			NameJa:        allProduct.NameJa,
			Allergy:       allProduct.Allergy,
			InformationJa: allProduct.InformationJa,
		}
	}

	// SQLiteの在庫には未会計の予約分が反映されていない。再起動や商品更新時に
	// Redis上の稼働中注文が再度販売可能にならないよう、予約数を差し引く。
	for _, place := range placeList {
		if place == nil {
			continue
		}
		needsReservation := place.PlaceStatus == PlaceStatusUsing ||
			(place.PlaceStatus == PlaceStatusBilling && place.BillingStage < BillingStageStockUpdated)
		if !needsReservation {
			continue
		}
		for _, order := range place.OrderList {
			if order.OrderStatus == OrderStatusCancel || order.JANCode == 0 || order.Count <= 0 {
				continue
			}
			inventoryMap[order.JANCode] -= order.Count
			if inventoryMap[order.JANCode] < 0 {
				inventoryMap[order.JANCode] = 0
			}
		}
	}

	receiver.mutex.Lock()
	previousAutoUpdateCancel := receiver.autoUpdateCancel
	receiver.autoUpdateCancel = nil
	receiver.Stock.AllProductList = append(AllProductList(nil), allProductList...)
	receiver.Stock.ProductList = productList
	receiver.Stock.IsValidJANCodeList = isValidJANCodeList

	receiver.TagPrimaryList = tagPrimaryList
	receiver.InventoryMap = inventoryMap
	receiver.MiniProductMap = miniProductMap
	receiver.PlaceList = clonePlaceList(placeList)
	receiver.AlcoholList = alcoholList
	autoUpdateContext := ctx
	isAutoUpdateEnabled := env.Env.EnvUnique.IsEnableCharge
	if isAutoUpdateEnabled {
		var cancel context.CancelFunc
		autoUpdateContext, cancel = context.WithCancel(ctx)
		receiver.autoUpdateCancel = cancel
	}
	receiver.mutex.Unlock()

	if previousAutoUpdateCancel != nil {
		previousAutoUpdateCancel()
	}
	if isAutoUpdateEnabled {
		go receiver.autoUpdate(autoUpdateContext)
	}
}

// SetPlaceList は、席一覧と各席の注文一覧をディープコピーしてドメインへ設定します。
func (receiver *Domain) SetPlaceList(
	ctx context.Context,
	placeList PlaceList,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	receiver.PlaceList = clonePlaceList(placeList)
}

// clone は、ドメイン状態のディープコピーを生成します。
func (receiver *Domain) clone() *Domain {
	inventoryMap := make(InventoryMap, len(receiver.InventoryMap))
	for janCode, count := range receiver.InventoryMap {
		inventoryMap[janCode] = count
	}
	miniProductMap := make(MiniProductMap, len(receiver.MiniProductMap))
	for janCode, product := range receiver.MiniProductMap {
		miniProductMap[janCode] = product
	}
	twentyList := make(TwentyList, len(receiver.AllowedTwentyList))
	for uuid, expiresAt := range receiver.AllowedTwentyList {
		twentyList[uuid] = expiresAt
	}
	pendingFaceList := make(TwentyList, len(receiver.PendingFaceList))
	for uuid, expiresAt := range receiver.PendingFaceList {
		pendingFaceList[uuid] = expiresAt
	}

	return &Domain{
		PlaceList: clonePlaceList(receiver.PlaceList),
		Stock: Stock{
			AllProductList:     append(AllProductList(nil), receiver.Stock.AllProductList...),
			ProductList:        append(ProductList(nil), receiver.Stock.ProductList...),
			IsValidJANCodeList: append([]int(nil), receiver.Stock.IsValidJANCodeList...),
		},
		AllergyDefault:    cloneAllergy(receiver.AllergyDefault),
		TagPrimaryList:    cloneTagPrimaryList(receiver.TagPrimaryList),
		InventoryMap:      inventoryMap,
		MiniProductMap:    miniProductMap,
		AlcoholList:       append(AlcoholList(nil), receiver.AlcoholList...),
		AllowedTwentyList: twentyList,
		PendingFaceList:   pendingFaceList,
	}
}

func cloneAllergy(source Allergy) Allergy {
	cloned := make(Allergy, len(source))
	for name, label := range source {
		cloned[name] = label
	}
	return cloned
}

func cloneTagPrimaryList(source TagPrimaryList) TagPrimaryList {
	cloned := make(TagPrimaryList, len(source))
	for index, primary := range source {
		cloned[index] = primary
		cloned[index].TagSecondaryList = append(TagSecondaryList(nil), primary.TagSecondaryList...)
	}
	return cloned
}
