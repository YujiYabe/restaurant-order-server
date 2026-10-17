package domain

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

// GetPlaceList は、現在の席一覧のスナップショットを返します。
func (receiver *PlaceList) GetPlaceList(
	ctx context.Context,
) (
	placeList PlaceList,
) {
	placeList = *receiver
	return
}

// GetPlace は、指定された席のスナップショットを返します。
func (receiver *PlaceList) GetPlace(
	ctx context.Context,
	placeName string,
) (
	place *Place,
) {
	place = (*receiver)[placeName]
	return
}

// EditPlace は、指定席に席状態の遷移ルールを適用します。
//
// 許可する遷移と主な副作用は次のとおりです。
//
//   - Cleaning: Opening、Using、Billing、Stoppingから遷移でき、時刻と料金をリセットして利用番号を加算します。
//   - Opening: Cleaningから遷移でき、利用者別情報、取引情報、注文一覧、および権限キーを初期化します。
//   - Using: Openingから遷移でき、入店時刻を記録します。
//   - Billing: Usingから遷移でき、退店時刻、取引ID、および会計進捗を記録します。
//   - Stopping: Cleaningからのみ遷移できます。
//
// 席が存在しない場合または遷移元が条件を満たさない場合は、状態を変更せずfalseを返します。
func (receiver *PlaceList) EditPlace(
	ctx context.Context,
	newPlace Place,
) (ok bool) {
	place, exists := (*receiver)[newPlace.PlaceName]
	if !exists || place == nil {
		return
	}

	defaultTime := new(time.Time)
	defaultLanguage := 1041
	defaultAllergy := "{\"egg\":false, \"milk\":false, \"shrimp\":false, \"crab\":false, \"squid\":false, \"salmon\":false, \"mackerel\":false, \"salmon_roe\":false, \"abalone\":false, \"peanut\":false, \"walnut\":false, \"soybean\":false, \"almond\":false, \"cashew_nut\":false, \"sesame\":false, \"wheat\":false, \"buckwheat\":false, \"apple\":false, \"orange\":false, \"banana\":false, \"peach\":false, \"kiwi_fruit\":false, \"pork\":false, \"chicken\":false, \"beef\":false, \"wild_yam\":false, \"matsutake_mushroom\":false, \"gelatin\":false}"

	placeStatus := place.PlaceStatus
	isStatusCleaning := placeStatus == PlaceStatusCleaning
	isStatusOpening := placeStatus == PlaceStatusOpening
	isStatusUsing := placeStatus == PlaceStatusUsing
	isStatusBilling := placeStatus == PlaceStatusBilling
	isStatusStopping := placeStatus == PlaceStatusStopping

	switch newPlace.PlaceStatus {
	// ___________________________________
	case PlaceStatusCleaning:
		// 状態が opening か billing か stopping であれば可能
		// using→cleaning 正常ルート
		// billing→cleaning 正常ルート
		if isStatusOpening || isStatusUsing || isStatusBilling || isStatusStopping {
			(*receiver)[newPlace.PlaceName].PlaceStatus = newPlace.PlaceStatus
			(*receiver)[newPlace.PlaceName].EnterTime = *defaultTime
			(*receiver)[newPlace.PlaceName].ExitTime = *defaultTime
			(*receiver)[newPlace.PlaceName].BaseCount = 0
			(*receiver)[newPlace.PlaceName].ChargePrice = 0
			(*receiver)[newPlace.PlaceName].OrderPrice = 0
			(*receiver)[newPlace.PlaceName].LineNumber++
			ok = true
		}

	case PlaceStatusOpening:
		// 状態が cleaning であれば可能
		// 店頭から席予約ができる状態
		if isStatusCleaning {
			// placeをリセット
			(*receiver)[newPlace.PlaceName].PlaceStatus = newPlace.PlaceStatus
			(*receiver)[newPlace.PlaceName].Allergy = defaultAllergy
			(*receiver)[newPlace.PlaceName].Language = defaultLanguage
			(*receiver)[newPlace.PlaceName].TransactionHeadId = ""
			(*receiver)[newPlace.PlaceName].BillingStage = BillingStageNone
			(*receiver)[newPlace.PlaceName].Permission = uuid.New().String()
			(*receiver)[newPlace.PlaceName].OrderList = NewOrderList()
			ok = true
		}

	case PlaceStatusUsing:
		// 状態が opening であれば可能
		// 注文が受け付けれる状態
		if isStatusOpening {
			(*receiver)[newPlace.PlaceName].PlaceStatus = newPlace.PlaceStatus
			(*receiver)[newPlace.PlaceName].EnterTime = time.Now()
			ok = true
		}

	case PlaceStatusBilling:
		// 状態が using であれば可能
		// 支払い状態 注文受付ができない
		if isStatusUsing {
			(*receiver)[newPlace.PlaceName].PlaceStatus = newPlace.PlaceStatus
			(*receiver)[newPlace.PlaceName].ExitTime = time.Now()
			if newPlace.TransactionHeadId != "" {
				(*receiver)[newPlace.PlaceName].TransactionHeadId = newPlace.TransactionHeadId
			}
			(*receiver)[newPlace.PlaceName].BillingStage = newPlace.BillingStage
			ok = true
		}

	case PlaceStatusStopping:
		// 状態がcleaningであれば可能
		if isStatusCleaning {
			(*receiver)[newPlace.PlaceName].PlaceStatus = newPlace.PlaceStatus
			ok = true
		}

	default:
		//
	}
	return
}

// AddOrder は、席の注文一覧の末尾に新しい注文番号を割り当てて注文を追加します。
// JANコード0のチャージ明細は初期状態を提供済みとし、通常商品は準備中とします。
// 追加後は取消済み注文を除外して席の注文金額を再計算します。
func (receiver *PlaceList) AddOrder(
	ctx context.Context,
	placeName string,
	newOrder Order,
	miniProductMap MiniProductMap,
) {
	place, ok := (*receiver)[placeName]
	if !ok || place == nil {
		return
	}

	index := 1 + len(place.OrderList)

	orderStatus := OrderStatusPreparing

	// チャージ料金は渡し済みにする
	if newOrder.JANCode == 0 {
		orderStatus = OrderStatusPassed
	}

	order := Order{
		OrderNumber: index,
		OrderStatus: orderStatus,
		JANCode:     newOrder.JANCode,
		Count:       newOrder.Count,
		OrderTime:   time.Now(),
	}

	(*receiver)[placeName].OrderList[index] = order
	(*receiver)[placeName].CalcPrice(ctx, miniProductMap)
}

// GetPermission は、指定席のインターネット注文用権限キーを返します。
func (receiver *PlaceList) GetPermission(
	ctx context.Context,
	placeName string,
) (
	permission string,
) {
	place, ok := (*receiver)[placeName]
	if !ok || place == nil {
		return
	}

	permission = place.Permission
	return
}

// EditOrder は、注文を取消、準備中、提供済みのいずれかへ変更します。
// 更新時に提供時刻を記録し、取消済み注文を除外して席の注文金額を再計算します。
// 注文が存在しない場合、同じ状態への変更、または遷移先が不正な場合はfalseを返します。
func (receiver *PlaceList) EditOrder(
	ctx context.Context,
	placeName string,
	targetOrder Order,
	miniProductMap MiniProductMap,
) (
	ok bool,
) {
	place, exists := (*receiver)[placeName]
	if !exists || place == nil {
		return
	}

	order, exists := place.OrderList[targetOrder.OrderNumber]
	if !exists {
		return
	}

	isValidTargetStatus := targetOrder.OrderStatus == OrderStatusCancel ||
		targetOrder.OrderStatus == OrderStatusPreparing ||
		targetOrder.OrderStatus == OrderStatusPassed
	if order.OrderStatus == targetOrder.OrderStatus || !isValidTargetStatus {
		return
	}

	order.OrderStatus = targetOrder.OrderStatus
	if targetOrder.OrderStatus == OrderStatusPreparing {
		order.OfferTime = time.Time{}
	} else {
		order.OfferTime = time.Now()
	}

	place.OrderList[targetOrder.OrderNumber] = order
	place.CalcPrice(ctx, miniProductMap)

	ok = true

	return
}

// BuildPostTransactionTemporary は、席の未取消注文をJANコードごとに集計し、
// スマレジ仮販売登録データとSQLite在庫更新用の数量マップを生成します。
// 商品小計がチャージ料金を下回る場合は、差額をJANコード0のチャージ明細として追加します。
// TerminalTranIdには席の利用番号、TerminalTranDateTimeには生成時刻を設定します。
func (receiver *PlaceList) BuildPostTransactionTemporary(
	ctx context.Context,
	placeName string,
	miniProductMap MiniProductMap,
) (
	reqPostTransactionTemporary ReqPostTransactionTemporary,
	janCodeCountMap map[int]int,
) {

	janCodeCountMap = make(map[int]int)
	for _, order := range (*receiver)[placeName].OrderList {
		if order.OrderStatus != OrderStatusCancel && order.JANCode != 0 {
			janCodeCountMap[order.JANCode] += order.Count
		}
	}

	details := []Detail{}
	transactionDetailId := 0
	subtotal := 0

	for janCode, count := range janCodeCountMap {
		productId := miniProductMap[janCode].ProductId
		salesPrice := miniProductMap[janCode].Price
		quantity := count

		subtotal += salesPrice * quantity

		transactionDetailId++
		detail := Detail{
			TransactionDetailId: transactionDetailId,
			SalesPrice:          salesPrice,
			Quantity:            quantity,
			ProductCode:         janCode,
			ProductId:           productId,
		}

		details = append(
			details,
			detail,
		)
	}

	baseCount := (*receiver)[placeName].BaseCount
	capacity := (*receiver)[placeName].Capacity
	chargeBasePrice := env.Env.EnvUnique.ChargeBasePrice
	placePrice := baseCount * capacity * chargeBasePrice

	if placePrice > subtotal {
		transactionDetailId++
		details = append(
			details,
			Detail{
				TransactionDetailId: transactionDetailId,
				SalesPrice:          placePrice - subtotal,
				Quantity:            1,
				ProductCode:         0,
				ProductId:           1000000000001,
			},
		)
		subtotal = placePrice
	}

	reqPostTransactionTemporary = ReqPostTransactionTemporary{
		Subtotal:             subtotal,
		Total:                subtotal,
		TerminalTranDateTime: time.Now().Format(time.RFC3339),
		TerminalTranId:       (*receiver)[placeName].LineNumber,
		CustomerId:           (*receiver)[placeName].CustomerId,
		Details:              details,
	}

	return
}

var IsAutoUpdate = make(chan bool)

// AutoUpdate は、Using中の席について、入店からの経過分数をChargeBaseMinuteで割った区切り数を求めます。
// 区切り数がBaseCountを上回った場合だけ、ChargeBasePriceと席の収容人数からChargePriceを更新し、
// trueを返します。Using中でない席や区切り数が増えていない席はfalseを返します。
func (receiver *Place) AutoUpdate(
	ctx context.Context,
	miniProductMap map[int]MiniProduct,
) (
	isUpdated bool,
) {
	if receiver.PlaceStatus != PlaceStatusUsing {
		return
	}

	timeNow := time.Now()
	// diffBaseCount := 1 + math.Floor(
	diffBaseCount := math.Floor(
		float64(
			timeNow.Sub(
				receiver.EnterTime,
			).Minutes(),
		) / float64(env.Env.ChargeBaseMinute),
	)

	receiver.CalcPrice(ctx, miniProductMap)

	// if receiver.PlaceName == "A1" {
	// log.Println("== == == == == == == == == == ")
	// log.Printf(
	// 	"%#v : %#v : %#v : %#v : %#v : %#v : %#v \n",
	// 	receiver.EnterTime.Format("15:04:05"),
	// 	timeNow.Format("15:04:05"),
	// 	receiver.PlaceName,
	// 	receiver.Capacity,
	// 	diffBaseCount,
	// 	receiver.ChargePrice,
	// 	price,
	// )
	// log.Println("== == == == == == == == == == ")
	// }

	if receiver.BaseCount < int(diffBaseCount) {
		receiver.ChargePrice = int(diffBaseCount) * env.Env.ChargeBasePrice * receiver.Capacity
		receiver.BaseCount = int(diffBaseCount)
		isUpdated = true
	}

	return
}

// CalcPrice は、取消注文を除外して席の注文金額を再計算します。
func (receiver *Place) CalcPrice(
	ctx context.Context,
	miniProductMap map[int]MiniProduct,
) {
	price := 0
	for _, order := range receiver.OrderList {
		if order.OrderStatus == OrderStatusCancel {
			continue
		}
		price += order.Count * miniProductMap[order.JANCode].Price
	}
	receiver.OrderPrice = price
}
