package smaregi

const (
	transactionHeadDivision = "1" // 固定 取引区分 仮販売登録APIは、1:通常のみ対応しています。
	status                  = "0" // 固定 ステータス 仮販売登録APIは、0:通常のみ対応しています。
	terminalId              = "1" // 端末ID
)

type credentialInfo struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// GetTransactionTemporaryReq は、仮販売取引の検索条件を表します。
type GetTransactionTemporaryReq struct {
	TransactionDateTimeFrom string `json:"transaction_date_time-from"` // 取引区分
	TransactionDateTimeTo   string `json:"transaction_date_time-to"`   // 取引区分
}

type GetTransactionTemporaryRes struct {
	TransactionHeadId              string  `json:"transactionHeadId"`
	TransactionDateTime            string  `json:"transactionDateTime"`
	TransactionHeadDivision        string  `json:"transactionHeadDivision"`
	CancelDivision                 string  `json:"cancelDivision"`
	Status                         string  `json:"status"`
	Barcode                        string  `json:"barcode"`
	SequentialNumber               string  `json:"sequentialNumber"`
	PreRegistrationName            *string `json:"preRegistrationName"`
	UnitNonDiscountsubtotal        string  `json:"unitNonDiscountsubtotal"`
	UnitDiscountsubtotal           string  `json:"unitDiscountsubtotal"`
	Subtotal                       string  `json:"subtotal"`
	SubtotalForDiscount            string  `json:"subtotalForDiscount"`
	SubtotalDiscountPrice          string  `json:"subtotalDiscountPrice"`
	SubtotalDiscountRate           string  `json:"subtotalDiscountRate"`
	SubtotalDiscountDivision       *string `json:"subtotalDiscountDivision"`
	PointDiscount                  string  `json:"pointDiscount"`
	Total                          string  `json:"total"`
	TaxInclude                     string  `json:"taxInclude"`
	TaxExclude                     string  `json:"taxExclude"`
	RoundingDivision               string  `json:"roundingDivision"`
	RoundingPrice                  string  `json:"roundingPrice"`
	Amount                         string  `json:"amount"`
	ReturnAmount                   string  `json:"returnAmount"`
	CostTotal                      string  `json:"costTotal"`
	NewPoint                       string  `json:"newPoint"`
	SpendPoint                     string  `json:"spendPoint"`
	Point                          string  `json:"point"`
	TotalPoint                     string  `json:"totalPoint"`
	StoreId                        string  `json:"storeId"`
	StoreCode                      string  `json:"storeCode"`
	TerminalId                     string  `json:"terminalId"`
	CustomerId                     *string `json:"customerId"`
	CustomerCode                   string  `json:"customerCode"`
	TerminalTranId                 string  `json:"terminalTranId"`
	TerminalTranDateTime           string  `json:"terminalTranDateTime"`
	CustomerGroupId                *string `json:"customerGroupId"`
	CustomerGroupId2               *string `json:"customerGroupId2"`
	CustomerGroupId3               *string `json:"customerGroupId3"`
	CustomerGroupId4               *string `json:"customerGroupId4"`
	CustomerGroupId5               *string `json:"customerGroupId5"`
	StaffId                        *string `json:"staffId"`
	StaffCode                      *string `json:"staffCode"`
	Memo                           *string `json:"memo"`
	ReceiptMemo                    *string `json:"receiptMemo"`
	GuestNumbers                   string  `json:"guestNumbers"`
	GuestNumbersMale               *string `json:"guestNumbersMale"`
	GuestNumbersFemale             *string `json:"guestNumbersFemale"`
	GuestNumbersUnknown            *string `json:"guestNumbersUnknown"`
	EnterDateTime                  string  `json:"enterDateTime"`
	TaxFreeSalesDivision           string  `json:"taxFreeSalesDivision"`
	NetTaxFreeGeneralTaxInclude    string  `json:"netTaxFreeGeneralTaxInclude"`
	NetTaxFreeGeneralTaxExclude    string  `json:"netTaxFreeGeneralTaxExclude"`
	NetTaxFreeConsumableTaxInclude string  `json:"netTaxFreeConsumableTaxInclude"`
	NetTaxFreeConsumableTaxExclude string  `json:"netTaxFreeConsumableTaxExclude"`
	Tags                           string  `json:"tags"`
	MileageDivision                string  `json:"mileageDivision"`
	MileageLabel                   string  `json:"mileageLabel"`
	ReturnSales                    string  `json:"returnSales"`
	SellDivision                   string  `json:"sellDivision"`
	TaxRate                        string  `json:"taxRate"`
	TaxRounding                    string  `json:"taxRounding"`
	DiscountRoundingDivision       string  `json:"discountRoundingDivision"`
}

type PostTransactionTemporaryReq struct {
	// 必須項目 =================================
	TransactionHeadDivision string `json:"transactionHeadDivision"` // 取引区分
	// 取引を識別する区分。
	// 1:通常、2:入金、3:出金、4:預かり金、 5:預かり金返金、6:ポイント加算、
	// 7:ポイント減算、8:ポイント失効、10:取置き、13:マイル加算、
	// 14:マイル減算、15:バリューカード入金、16:領収証
	// 仮販売登録APIは、1:通常のみ対応しています。

	Status string `json:"status"` // ステータス
	// 値：0
	// 仮販売のステータス。（0:通常、1:完了、2:取消）
	// 仮販売登録APIは、0:通常のみ対応しています。

	Subtotal string `json:"subtotal"` // 小計
	// 取引明細の値引き後小計の合計。
	// 明細合計と異なる場合エラー

	Total string `json:"total"` // 合計
	// 小計     - 小計値引き            - ポイント値引き + 端数値引額    + 外税額     - 免税額 を設定
	// 計算式：Subtotal - SubtotalDiscountPrice - PointDiscount + RoundingPrice + TaxExclude - TaxFree
	// 端数値引き機能は現在ご利用いただけません

	StoreId string `json:"storeId"` // 店舗ID
	// 店舗毎に付与するID。
	// 店舗情報が存在しない場合エラー
	// 休止中の店舗IDを指定するとエラー

	TerminalId string `json:"terminalId"` // 端末ID
	// 店舗の端末毎に付与するID。
	// 存在しない端末IDも設定可能。

	TerminalTranId string `json:"terminalTranId"` // 端末取引ID
	// 10文字以内
	// 端末で設定された取引ID。(連携元システムで任意の連番を設定してください）

	TerminalTranDateTime string `json:"terminalTranDateTime"` // 端末取引日時
	// 端末で設定された取引日時。[YYYY-MM-DDThh:mm:ssTZD]
	// （連携元システムで処理日時を設定してください）

	Details []Detail `json:"details"` // 取引明細

	CustomerCode string `json:"customerCode,omitempty"` // 会員コード
	// 20文字以内 /^[ -~]+$/
	// 会員毎に付与する一意のコード。
	// 設定された場合、会員IDは必須。
	// 会員マスタに存在する会員IDが設定されている場合、会員情報の会員コードと一致しないとエラー。
	// 省略時、会員IDが設定されている場合、会員マスタの会員コードを設定。
	// （スマレジ管理画面ver.4.12.0リリース以降

	// 使用必須項目 =================================
	CustomerId string `json:"customerId"` // 会員ID
	// 会員毎に付与するID。
	// 存在しない会員IDの場合、そのまま会員IDは登録されますが、会員ポイント情報等は更新しません。

	// 任意項目 =================================
	Barcode *string `json:"barcode,omitempty"` // バーコード
	// 20文字以内
	// ユニーク値 バーコード
	// 省略時はシステムで一意のバーコードを設定

	SequentialNumber *string `json:"sequentialNumber,omitempty"` // 連番
	// 連番
	// 省略時はシステムで日毎の連番を設定

	PreRegistrationName *string `json:"preRegistrationName,omitempty"` // 仮販売名
	// 20文字以内
	// 仮販売の名称を設定。20文字以内。
	// 設定した場合、スマレジアプリ取引機能において、仮販売一覧に表示されます。
	// 設定しない場合は連番が表示されます。

	SubtotalDiscountPrice *string `json:"subtotalDiscountPrice,omitempty"` // 小計値引き
	// 小計の値引き金額を設定。

	SubtotalDiscountRate *string `json:"subtotalDiscountRate,omitempty"` // 小計割引率
	// 小計の割引を設定 （単位％） デフォルトは未設定。
	// 設定された場合、小計値引は必須。
	// 小計値引率から小計値引は計算しないため小計値引に値引額を設定してください。

	SubtotalDiscountDivision *string `json:"subtotalDiscountDivision,omitempty"` // 小計値引き/割引区分
	// 小計値引き/割引時に指定された区分を設定。
	// 管理画面->設定->販売設定->値引区分/割引区分のページで確認できるID。

	PointDiscount *string `json:"pointDiscount,omitempty"` // ポイント値引き
	// 使用ポイント * ポイント還元率。

	TaxInclude *string `json:"taxInclude,omitempty"` // 内税額
	// 内税商品計 / 110 * 10。
	// 税率が10%の場合
	// 小数点以下は税の丸め方式で指定の通り

	TaxExclude *string `json:"taxExclude,omitempty"` // 外税額
	// 外税商品計 * 10 / 100。
	// 税率が10%の場合
	// 小数点以下は税の丸め方式で指定の通り

	RoundingDivision *string `json:"roundingDivision,omitempty"` // 端数値引区分
	// 列挙値："00" "11" "12" "13" "21" "22" "23" "31" "32" "33" "41" "42" "43" "99"
	// 販売時の合計金額の端数を丸めるか否かを設定する。
	// 00:なし
	// 11:一の位を四捨五入 12:一の位を切捨 13:一の位を切上
	// 21:十の位を四捨五入 22:十の位を切捨 23:十の位を切上
	// 31:百の位を四捨五入 32:百の位を切捨 33:百の位を切上
	// 41:千の位を四捨五入 42:千の位を切捨 43:千の位を切上
	// 99:手入力

	RoundingPrice *string `json:"roundingPrice,omitempty"` // 端数値引額

	NewPoint *string `json:"newPoint,omitempty"` // 付与ポイント
	// 既定値：0
	// 取引で付与されるポイント。 合計 * ポイント付与率
	// 会員IDが設定されている場合、会員ポイント情報を更新。
	// 省略時は0を設定

	SpendPoint *string `json:"spendPoint,omitempty"` // 使用ポイント
	// 取引で使用したポイント。
	// 会員IDが設定されている場合、会員ポイント情報を更新。
	// 省略時は0を設定

	Point *string `json:"point,omitempty"` // 現在ポイント
	// 取引前のポイント。

	TotalPoint *string `json:"totalPoint,omitempty"` // 合計ポイント
	// 取引後のポイント 現在ポイント + 付与ポイント - 使用ポイント。

	CustomerGroupId *string `json:"customerGroupId,omitempty"` // 客層ID
	// 客層ID。客層セクションIDが1の客層IDを設定。
	// 客層情報が存在しない場合エラー

	CustomerGroupId2 *string `json:"customerGroupId2,omitempty"` // 客層ID2
	// 客層ID2。客層セクションIDが2の客層IDを設定。
	// 客層情報が存在しない場合エラー

	CustomerGroupId3 *string `json:"customerGroupId3,omitempty"` // 客層ID3
	// 客層ID3。客層セクションIDが3の客層IDを設定。
	// 客層情報が存在しない場合エラー

	CustomerGroupId4 *string `json:"customerGroupId4,omitempty"` // 客層ID4
	// 客層ID4。客層セクションIDが4の客層IDを設定。
	// 客層情報が存在しない場合エラー

	CustomerGroupId5 *string `json:"customerGroupId5,omitempty"` // 客層ID5
	// 客層ID5。客層セクションIDが5の客層IDを設定。
	// 客層情報が存在しない場合エラー

	StaffId *string `json:"staffId,omitempty"` // 販売員ID
	// 販売員ID（スタッフID）
	// スタッフ情報が存在しない場合エラー
	// ユーザーアクセストークンを利用する場合、ユーザーの所属する店舗の販売員IDを指定してください。

	Memo *string `json:"memo,omitempty"` // メモ
	// 入出金時のメモなどを設定。 1000文字以内。

	ReceiptMemo *string `json:"receiptMemo,omitempty"` // レシートメモ
	// レシートに印字するメモを設定。100文字以内。

	GuestNumbers *string `json:"guestNumbers,omitempty"` // 客数
	// 客数。

	GuestNumbersMale *string `json:"guestNumbersMale,omitempty"` // 客数（男）
	// 客数（男）。

	GuestNumbersFemale *string `json:"guestNumbersFemale,omitempty"` // 客数（女
	// 客数（女）。

	GuestNumbersUnknown *string `json:"guestNumbersUnknown,omitempty"` // 客数（不明
	// 客数（不明）。

	EnterDateTime *string `json:"enterDateTime,omitempty"` // 入店日時
	// 入店日時。[YYYY-MM-DDThh:mm:ssTZD]

	TaxFreeSalesDivision *string `json:"taxFreeSalesDivision,omitempty"` // 免税販売区分
	// 列挙値：0 1 2 3 4 5、既定値：0
	// 免税販売区分
	// 0:通常販売
	// 1:免税販売(自動計算)
	// 2:免税販売(強制免税適用)
	// 3:免税販売(一般品のみ強制免税適用)
	// 4:免税販売(消耗品のみ強制免税適用)
	// 5:免税販売(自動計算・合算)
	// 未設定の場合、0:通常販売

	NetTaxFreeGeneralTaxInclude *string `json:"netTaxFreeGeneralTaxInclude,omitempty"` // 内税商品の一般品免税額
	// 内税商品の一般品免税額。

	NetTaxFreeGeneralTaxExclude *string `json:"netTaxFreeGeneralTaxExclude,omitempty"` // 外税商品の一般品免税額
	// 外税商品の一般品免税額。

	NetTaxFreeConsumableTaxInclude *string `json:"netTaxFreeConsumableTaxInclude,omitempty"` // 内税商品の消耗品免税額
	// 内税商品の消耗品免税額。

	NetTaxFreeConsumableTaxExclude *string `json:"netTaxFreeConsumableTaxExclude,omitempty"` // 外税商品の消耗品免税額
	// 外税商品の消耗品免税額。

	Tags *string `json:"tags,omitempty"` // 取引タグ
	// レジ端末よりクイックコマンドを押下した値を設定。 和名を設定する。
	// 1取引タグは最大64文字。複数ある場合は、カンマ区切りで設定。

	MileageDivision *string `json:"mileageDivision,omitempty"` // マイレージ利用区分
	// 既定値：0
	// 列挙値：0 1 2
	// マイレージ機能を利用するか否か。
	// (0: 利用しない 1: 利用する 2: ポイントの代わりに貯める)
	// 未設定の場合は店舗マスタから設定。 初期値:0。

	MileageLabel *string `json:"mileageLabel,omitempty"` // マイレージラベル
	// マイレージの名称。未設定の場合は「マイル」

	SellDivision *string `json:"sellDivision,omitempty"` // 販売区分
	// 列挙値：0 1
	// 取引が外税販売、内税販売かを識別する区分。
	// （0:内税販売、1:外税販売）
	// 省略時は0:内税販売を設定

	TaxRate *string `json:"taxRate,omitempty"` // 税率
	// 税率 （単位％）
	// 省略時は、スマレジに設定されている消費税を設定（端末取引日時より算出）

	TaxRounding *string `json:"taxRounding,omitempty"` // 税の丸め方式
	// 列挙値：0 1 2
	// 税の丸めを判定。 0:四捨五入、1:切り捨て、2:切り上げ
	// 省略時はスマレジに設定されている税の丸め方式を設定

	DiscountRoundingDivision *string `json:"discountRoundingDivision,omitempty"` // 割引丸め区分
	// 列挙値：0 1 2
	// 割引の端数を丸める方法を指定する区分。 0:四捨五入、1:切り捨て、2:切り上げ
	// 省略時は店舗マスタに設定されている税の丸め方式を設定

}

type Detail struct {
	// 必須項目 =================================
	TransactionDetailId string `json:"transactionDetailId"` // 取引明細ID
	// ユニーク値 取引毎の明細に付与するID。
	// 設定する場合は1〜999まで。
	// 取引明細を指定する場合必須。

	TransactionDetailDivision string `json:"transactionDetailDivision"` // 取引明細区分
	// 列挙値：1 2 3
	// 取引明細を識別する区分。
	// （1:通常、2:返品、3:部門売り）
	// 取引明細を指定する場合必須。

	SalesPrice string `json:"salesPrice"` // 販売単価
	// 販売時の単価を設定。
	// 取引明細を指定する場合必須。

	Quantity string `json:"quantity"` // 数量
	// 購入数
	// 取引明細を指定する場合必須。

	// 使用必須項目 =================================
	ProductCode string `json:"productCode,omitempty"` // 商品コード
	// 商品コード。（半角英数記号 20桁以内）
	// 省略時、商品IDが設定されている場合、マスタから設定。以外はnull

	ProductId string `json:"productId"` // 商品ID
	// 商品毎に付与するID
	// 存在しない商品IDが設定されている場合、nullを登録。
	// 商品IDで指定された商品の持つ商品区分が 1:回数券 の明細は登録できません。
	// ユーザーアクセストークンを利用する場合、ユーザーの所属する店舗で販売している商品IDを指定してください。

	// 任意項目 =================================
	ParentTransactionDetailId *string `json:"parentTransactionDetailId,omitempty"` // 親取引明細ID
	// 紐付く親の取引明細がある場合、その取引明細IDを設定。
	// 商品区分が、0（通常販売）, 4（バンドル(親)）の場合は設定できず、
	// 2（オプション商品）の場合はオプション元の取引明細IDを、
	// 7（バンドル子）, 8（バンドル値引）の時は、バンドル(親)の取引明細IDを設定しなければならない。

	ProductName *string `json:"productName,omitempty"` // 商品名
	// 省略時、商品IDが設定されている場合、マスタから設定。以外はnull

	PrintReceiptProductName *string `json:"printReceiptProductName,omitempty"` // レシート印字商品名
	// 省略時、商品IDが設定されている場合、マスタから設定。以外はnull

	TaxDivision *string `json:"taxDivision,omitempty"` // 税区分
	// 列挙値：0 1 2
	// (0:税込、1:税抜、2:非課税)
	// 設定時そのまま設定
	// 省略時、商品IDが設定されている場合、マスタから設定。以外は0:税込

	Price *string `json:"price,omitempty"` // 商品単価
	// 設定時そのまま設定
	// 省略時、商品IDが設定されている場合、マスタから設定。以外は0円

	UnitDiscountPrice *string `json:"unitDiscountPrice,omitempty"` // 単品値引
	// 販売単価の値引き金額を設定。

	UnitDiscountRate *string `json:"unitDiscountRate,omitempty"` // 単品割引率
	// 単品の割引率を設定。（単位:％）
	// このフィールドはラベルであり、販売単価の値引き金額の算出には利用されません。
	// 設定された場合、単品値引は必須

	UnitDiscountDivision *string `json:"unitDiscountDivision,omitempty"` // 単品値引き/割引区分
	// 単品値引き/割引時に指定された割引区分を設定。

	CategoryId *string `json:"categoryId,omitempty"` // 部門ID
	// 商品に紐付く部門毎に付与するID。数字9桁以内。
	// 部門情報が存在しない場合、エラー
	// 商品IDが設定されている場合、商品情報の部門IDと一致しないとエラー
	// 省略時、商品IDが設定されている場合、商品マスタの部門IDを設定

	CategoryName *string `json:"categoryName,omitempty"` // 部門名
	// 商品に紐付く部門名。(85文字以内)
	// 設定時、そのまま設定。
	// 省略時、部門IDが設定されている場合、マスタから設定。以外はnull

	SalesDivision *string `json:"salesDivision,omitempty"` // 売上区分
	// 列挙値：0 1
	// 売上計上するか否かの区分。（0:売上対象、1:売上対象外）
	// 省略時、商品IDが設定されている場合、商品マスタから設定。以外は0:売上対象を設定

	ProductDivision *string `json:"productDivision,omitempty"` // 商品区分
	// 列挙値：0 2 4 7 8
	// 商品区分 (0:通常、1:回数券、2:オプション商品、3:部門販売、4:バンドル (親)、5:セット売(親)、6:簡易セット(親)、7:バンドル商品(子)、8:バンドル値 引(子)、A:プラン(子)、B:コース(子)) 0, 2, 4, 7, 8のみ許可。商品区分を0, 4にする場合、親取引明細IDは設定不可。また、商品区分を2（オプション商品）にする場合は親取引明細IDにオプション元の取引明細IDの、7, 8にする場合はバンドル(親)の取引明細IDの設定が必要
	// A,Bはウェイターでのみ利用する区分。

	CalcDiscount *string `json:"calcDiscount,omitempty"` // 値引割引計算対象
	// 列挙値：0 1
	// 小計値引／割引(クーポン、ポイント値引き含む)の対象かどうかを設定。（0:対象外、1:対象）
	// 省略時、商品IDが設定されている場合、商品マスタから設定。その他の場合は1:対象を設定。
	// スマレジ・アプリver.4.27.0以降で取引登録すると、取引明細区分が2:返品の取引明細は仮販売登録時の値引割引計算対象設定値に関わらず、0:対象外として再計算します。

	TaxFreeDivision *string `json:"taxFreeDivision,omitempty"` // 免税区分
	// 列挙値：0 1 2
	// 免税区分の設定。(0:対象外、1:一般品、2:消耗品)
	// 省略時、商品区分が0:通常、2:オプション商品、7:バンドル商品(子)、8:バンドル値引(子)で商品IDが設定されている場合、商品マスタから設定。商品区分が4:バンドル (親)で商品バンドルグループIDが設定されている場合、商品バンドルグループマスタから設定。以外は0:対象外を設定

	TaxFreeCommodityPrice *string `json:"taxFreeCommodityPrice,omitempty"` // 免税対象額
	// 免税後の税抜金額を設定。（税抜価格x数量）
	// 免税販売時に未設定だと、日報の表示がずれる事があるので注意してください。

	TaxFree *string `json:"taxFree,omitempty"` // 免税額
	// 免税額を設定。（免税対象額 x 税率）
	// 免税販売時に未設定だと、日報の表示がずれる事があるので注意してください。

	ProductBundleGroupId *string `json:"productBundleGroupId,omitempty"` // 商品バンドルグループID
	// 文字列形式の整数
	// (商品バンドルグループID)
	// [ 1 .. 999999999 ]
	// 商品バンドル販売の明細の場合、商品バンドルグループのIDを設定。
	// 商品バンドルグループに存在するIDを指定しなければならない。
	// この項目を設定する明細の商品区分は「4:バンドル(親)」でなければならない

	ReduceTaxId *string `json:"reduceTaxId,omitempty"` // 軽減税率ID（税設定）
	// 文字列形式の整数
	// [ 1 .. 999999999 ]
	// 軽減税率ID
	// 標準:null（または未設定）
	// 軽減:10000001 (特定商品の軽減税率適用）
	// 選択[標準]:10000002（状態による適用[適用しない]）
	// 選択[軽減]:10000003（状態による適用[適用する]）
	// 選択[選択]:10000004（状態による適用[都度選択する]）
	// 上記以外の場合、カスタム軽減税率で登録されている軽減税率ID

	// 軽減税率販売時のみ設定。
	// 10000001〜10000004における軽減税率は日本国における軽減税率(8%)です。
	// 存在しない軽減税率IDを指定した場合エラー。

	Memo *string `json:"memo,omitempty"` // メモ
	// 文字列
	// 1000文字以内
	// 返品販売時などのメモなどを設定。 1000文字以内。
}

type PostTransactionTemporaryRes struct {
	TransactionHeadId        string  `json:"transactionHeadId"`
	TransactionDateTime      string  `json:"transactionDateTime"`
	TransactionHeadDivision  string  `json:"transactionHeadDivision"`
	CancelDivision           string  `json:"cancelDivision"`
	Status                   string  `json:"status"`
	Barcode                  string  `json:"barcode"`
	SequentialNumber         string  `json:"sequentialNumber"`
	PreRegistrationName      *string `json:"preRegistrationName"`
	UnitNonDiscountsubtotal  string  `json:"unitNonDiscountsubtotal"`
	UnitDiscountsubtotal     string  `json:"unitDiscountsubtotal"`
	Subtotal                 string  `json:"subtotal"`
	SubtotalForDiscount      string  `json:"subtotalForDiscount"`
	SubtotalDiscountPrice    string  `json:"subtotalDiscountPrice"`
	SubtotalDiscountRate     string  `json:"subtotalDiscountRate"`
	SubtotalDiscountDivision *string `json:"subtotalDiscountDivision"`
	PointDiscount            string  `json:"pointDiscount"`
	Total                    string  `json:"total"`
	TaxInclude               string  `json:"taxInclude"`
	TaxExclude               string  `json:"taxExclude"`
	RoundingDivision         string  `json:"roundingDivision"`
	RoundingPrice            string  `json:"roundingPrice"`
	Amount                   string  `json:"amount"`
	ReturnAmount             string  `json:"returnAmount"`
	CostTotal                string  `json:"costTotal"`
	NewPoint                 string  `json:"newPoint"`
	SpendPoint               string  `json:"spendPoint"`
	Point                    string  `json:"point"`
	TotalPoint               string  `json:"totalPoint"`
	StoreId                  string  `json:"storeId"`
	// StoreCode                      string  `json:"storeCode"`
	TerminalId           string  `json:"terminalId"`
	CustomerId           *string `json:"customerId"`
	CustomerCode         string  `json:"customerCode"`
	TerminalTranId       string  `json:"terminalTranId"`
	TerminalTranDateTime string  `json:"terminalTranDateTime"`
	CustomerGroupId      *string `json:"customerGroupId"`
	CustomerGroupId2     *string `json:"customerGroupId2"`
	CustomerGroupId3     *string `json:"customerGroupId3"`
	CustomerGroupId4     *string `json:"customerGroupId4"`
	CustomerGroupId5     *string `json:"customerGroupId5"`
	StaffId              *string `json:"staffId"`
	// StaffCode                      *string `json:"staffCode"`
	Memo                           *string `json:"memo"`
	ReceiptMemo                    *string `json:"receiptMemo"`
	GuestNumbers                   string  `json:"guestNumbers"`
	GuestNumbersMale               *string `json:"guestNumbersMale"`
	GuestNumbersFemale             *string `json:"guestNumbersFemale"`
	GuestNumbersUnknown            *string `json:"guestNumbersUnknown"`
	EnterDateTime                  string  `json:"enterDateTime"`
	TaxFreeSalesDivision           string  `json:"taxFreeSalesDivision"`
	NetTaxFreeGeneralTaxInclude    string  `json:"netTaxFreeGeneralTaxInclude"`
	NetTaxFreeGeneralTaxExclude    string  `json:"netTaxFreeGeneralTaxExclude"`
	NetTaxFreeConsumableTaxInclude string  `json:"netTaxFreeConsumableTaxInclude"`
	NetTaxFreeConsumableTaxExclude string  `json:"netTaxFreeConsumableTaxExclude"`
	Tags                           string  `json:"tags"`
	MileageDivision                string  `json:"mileageDivision"`
	MileageLabel                   string  `json:"mileageLabel"`
	ReturnSales                    string  `json:"returnSales"`
	SellDivision                   string  `json:"sellDivision"`
	TaxRate                        string  `json:"taxRate"`
	TaxRounding                    string  `json:"taxRounding"`
	DiscountRoundingDivision       string  `json:"discountRoundingDivision"`
	// Details                        []Detail `json:"details"` // 取引明細
}
