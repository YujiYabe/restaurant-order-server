package domain

type (
	// Stock は、全商品、販売可能商品、および注文可能なJANコードを保持します。
	Stock struct {
		AllProductList     AllProductList
		ProductList        ProductList
		IsValidJANCodeList []int
	}

	// AllProductList は、有効・無効を含む全商品の一覧です。
	AllProductList []Product
	// ProductList は、現在販売可能な商品の一覧です。
	ProductList []Product
	// InventoryMap は、JANコードごとの現在注文可能な在庫数を保持します。
	InventoryMap map[int]int
	// MiniProductMap は、注文・会計計算に必要な商品項目をJANコードごとに保持します。
	MiniProductMap map[int]MiniProduct

	// MiniProduct は、注文表示と価格計算に必要な商品の概要です。
	MiniProduct struct {
		JANCode       int
		ProductId     int
		Price         int
		Place         string
		Recipe        string
		Tag           string
		NameJa        string
		Allergy       string
		InformationJa string
	}
)

// NewInventoryMap は、空のJANコード別在庫マップを生成します。
func NewInventoryMap() (
	inventoryMap InventoryMap,
) {
	inventoryMap = InventoryMap{}
	return
}

// NewMiniProductMap は、空のJANコード別商品概要マップを生成します。
func NewMiniProductMap() (
	miniProductMap MiniProductMap,
) {
	miniProductMap = MiniProductMap{}
	return
}

// NewStock は、商品一覧が未設定の空のStockを生成します。
func NewStock() (
	stock Stock,
) {
	stock = Stock{}
	return
}

// Product は、商品マスターを表します。
type Product struct {
	JANCode       int
	ProductId     int
	IsValid       bool
	Stock         int
	Place         string
	Calorie       int
	CostPrice     int
	RetailPrice   int
	OperationCost int
	EnergyCost    int
	Recipe        string
	Tag           string
	Allergy       string
	NameJa        string
	NameJaShort   string
	InformationJa string
	UpdatedAt     string
}
