package main

type ProductAttribute int

const (
	JANCode ProductAttribute = iota
	ProductId
	IsValid
	Stock
	Place
	Calorie
	CostPrice
	RetailPrice
	OperationCost
	EnergyCost
	Recipe
	Tag
	Allergy
	NameJa
	NameJaShort
	InformationJa
	UpdatedAt
)

type Product struct {
	JANCode       int    `gorm:"primaryKey;column:jan_code" json:"jan_code"`       // jan_code
	ProductId     int    `gorm:"column:product_id"          json:"product_id"`     // smaregi product id
	IsValid       bool   `gorm:"column:is_valid"            json:"is_valid"`       // 有効無効
	Stock         int    `gorm:"column:stock"               json:"stock"`          // 在庫
	Place         string `gorm:"column:place"               json:"place"`          // 冷蔵庫の場所
	Calorie       int    `gorm:"column:calorie"             json:"calorie"`        // カロリー
	CostPrice     int    `gorm:"column:cost_price"          json:"cost_price"`     // 原価
	RetailPrice   int    `gorm:"column:retail_price"        json:"retail_price"`   // 小売価格
	OperationCost int    `gorm:"column:operation_cost"      json:"operation_cost"` // 提供するまでかかった時間の価格
	EnergyCost    int    `gorm:"column:energy_cost"         json:"energy_cost"`    // 提供するまでかかった電気料金
	Recipe        string `gorm:"column:recipe"              json:"recipe"`         // レシピ
	Tag           string `gorm:"column:tag"                 json:"tag"`            // タグ情報 json:array
	Allergy       string `gorm:"column:allergy"             json:"allergy"`        // アレルギー情報 json:object
	NameJa        string `gorm:"column:name_ja"             json:"name_ja"`        // 名前_ja
	NameJaShort   string `gorm:"column:name_ja_short"       json:"name_ja_short"`  // 名前_ja_short
	InformationJa string `gorm:"column:information_ja"      json:"information_ja"` // 商品のキャッチコピー_ja
	UpdatedAt     string `gorm:"column:updated_at"          json:"updated_at"`     // 更新日時
}

func (Product) TableName() (
	tableName string,
) {
	tableName = "product"
	return
}

type TagPrimary struct {
	ID   int    `gorm:"primaryKey;column:id"    json:"id"`   // id
	Name string `gorm:"column:name"             json:"name"` // 名前
}

func (TagPrimary) TableName() (
	tableName string,
) {
	tableName = "tag_primary"
	return
}

type TagSecondary struct {
	ID           int    `gorm:"primaryKey;column:id"  json:"id"`             // id
	TagPrimaryID int    `gorm:"column:primary_tag_id" json:"primary_tag_id"` // primary_tag_id
	Name         string `gorm:"column:name"           json:"name"`           // 名前
}

func (TagSecondary) TableName() (
	tableName string,
) {
	tableName = "tag_secondary"
	return
}
