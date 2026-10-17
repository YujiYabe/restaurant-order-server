package domain

import (
	"time"
)

// StatusNum は、注文または席の業務状態を表します。
type StatusNum int

// BillingStage は、会計連携の完了段階を表し、失敗後の再開位置として使用します。
type BillingStage int

const (
	// OrderStatusCancel は、注文が取り消された状態です。
	OrderStatusCancel StatusNum = 0
	// OrderStatusPreparing は、注文した商品を準備中の状態です。
	OrderStatusPreparing StatusNum = 1
	// OrderStatusPassed は、注文した商品を提供済みの状態です。
	OrderStatusPassed StatusNum = 2

	// PlaceStatusCleaning は、席を清掃中または次の利用に向けて整備中の状態です。
	PlaceStatusCleaning StatusNum = 0 // 掃除
	// PlaceStatusOpening は、利用者を受け入れられる空席状態です。
	PlaceStatusOpening StatusNum = 1 // 店頭で席を確保する 開始前の状況
	// PlaceStatusUsing は、利用者が着席して注文できる状態です。
	PlaceStatusUsing StatusNum = 2 // 利用中
	// PlaceStatusBilling は、注文受付を停止して会計中の状態です。
	PlaceStatusBilling StatusNum = 3 // 支払い中
	// PlaceStatusStopping は、端末故障などにより席の稼働を停止した状態です。
	PlaceStatusStopping StatusNum = 4 // 端末故障などで稼働できない場合
)

const (
	// BillingStageNone は、会計連携をまだ開始していない状態です。
	BillingStageNone BillingStage = iota
	// BillingStageTransactionCreated は、スマレジの仮販売登録が完了した状態です。
	BillingStageTransactionCreated
	// BillingStageOrderSaved は、注文履歴のFirestore保存が完了した状態です。
	BillingStageOrderSaved
	// BillingStageTransactionInfoSent は、会計連携サービスへの取引通知が完了した状態です。
	BillingStageTransactionInfoSent
	// BillingStageStockUpdated は、SQLiteの在庫更新まで完了した状態です。
	BillingStageStockUpdated
)

type (
	// PlaceList は、席名をキーとする席情報のマップです。
	PlaceList map[string]*Place

	// Place は、席の設定、利用状態、会計進捗、および注文一覧を保持します。
	Place struct {
		PlaceGroup        string
		PlaceName         string
		Capacity          int
		CustomerId        int
		LineNumber        int
		TransactionHeadId string
		BillingStage      BillingStage
		PlaceStatus       StatusNum
		EnterTime         time.Time
		ExitTime          time.Time
		Allergy           string
		Language          int
		Permission        string
		BaseCount         int
		ChargePrice       int
		OrderPrice        int
		OrderList         OrderList
	}

	// SmaregiPostResponseEnquiry は、スマレジが返した顧客IDと取引IDを保持します。
	SmaregiPostResponseEnquiry struct {
		CustomerId        int
		TransactionHeadId string
	}

	// OrderList は、注文番号をキーとする席ごとの注文マップです。
	OrderList map[int]Order

	// Order は、商品、数量、提供状態、および時刻を表す注文明細です。
	Order struct {
		OrderNumber int
		OrderStatus StatusNum
		JANCode     int
		Count       int
		OrderTime   time.Time
		OfferTime   time.Time
	}

	// Alcohol は、成人確認用UUID、画像、および書類種別を保持します。
	Alcohol struct {
		UUID            string
		Base64Source    string
		CertificateType string
	}

	// OrderAlcohol は、酒類注文と成人確認情報を組み合わせた入力です。
	OrderAlcohol struct {
		Order   Order
		Alcohol Alcohol
	}

	// ReqPostTransactionTemporary は、スマレジの仮販売登録に使用する取引情報です。
	ReqPostTransactionTemporary struct {
		Subtotal int // 小計
		// 取引明細の値引き後小計の合計。
		// 明細合計と異なる場合エラー

		Total int // 合計
		// 小計     - 小計値引き            - ポイント値引き + 端数値引額    + 外税額     - 免税額 を設定
		// 計算式：Subtotal - SubtotalDiscountPrice - PointDiscount + RoundingPrice + TaxExclude - TaxFree
		// 端数値引き機能は現在ご利用いただけません

		TerminalTranDateTime string // 端末取引日時
		// 端末で設定された取引日時。[YYYY-MM-DDThh:mm:ssTZD]
		// （連携元システムで処理日時を設定してください）

		TerminalTranId int

		CustomerId int // 会員ID テーブル
		// 会員毎に付与するID。
		// 存在しない会員IDの場合、そのまま会員IDは登録されますが、会員ポイント情報等は更新しません。
		Details []Detail
	}

	// Detail は、スマレジ仮販売取引の商品別明細です。
	Detail struct {
		TransactionDetailId int // 取引明細ID
		// ユニーク値 取引毎の明細に付与するID。
		// 設定する場合は1〜999まで。
		// 取引明細を指定する場合必須。

		SalesPrice int // 販売単価
		// 販売時の単価を設定。
		// 取引明細を指定する場合必須。

		Quantity int // 数量
		// 購入数
		// 取引明細を指定する場合必須。

		ProductCode int // 商品コード jan code
		ProductId   int // 商品ID
	}
)

// NewPlaceList は、初期席X0を含む席一覧を生成します。
func NewPlaceList() (
	placeList PlaceList,
) {
	placeList = PlaceList{
		"X0": &Place{},
	}
	return
}

// NewOrderList は、注文が未登録の空の注文マップを生成します。
func NewOrderList() (
	orderList OrderList,
) {
	orderList = OrderList{}
	return
}
