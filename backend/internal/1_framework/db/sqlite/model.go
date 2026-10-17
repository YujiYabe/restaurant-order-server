package sqlite

import domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"

// productModel は、商品のSQLite永続化表現です。
type productModel struct {
	JANCode       int    `gorm:"primaryKey;column:jan_code"`
	ProductId     int    `gorm:"column:product_id"`
	IsValid       bool   `gorm:"column:is_valid"`
	Stock         int    `gorm:"column:stock"`
	Place         string `gorm:"column:place"`
	Calorie       int    `gorm:"column:calorie"`
	CostPrice     int    `gorm:"column:cost_price"`
	RetailPrice   int    `gorm:"column:retail_price"`
	OperationCost int    `gorm:"column:operation_cost"`
	EnergyCost    int    `gorm:"column:energy_cost"`
	Recipe        string `gorm:"column:recipe"`
	Tag           string `gorm:"column:tag"`
	Allergy       string `gorm:"column:allergy"`
	NameJa        string `gorm:"column:name_ja"`
	NameJaShort   string `gorm:"column:name_ja_short"`
	InformationJa string `gorm:"column:information_ja"`
	UpdatedAt     string `gorm:"column:updated_at"`
}

func (productModel) TableName() string { return "product" }

func productModelFromDomain(product domain.Product) productModel {
	return productModel(product)
}

func (product productModel) domain() domain.Product { return domain.Product(product) }

type tagPrimaryModel struct {
	ID               int                 `gorm:"primaryKey;column:id"`
	Name             string              `gorm:"column:name"`
	TagSecondaryList []tagSecondaryModel `gorm:"foreignKey:TagPrimaryID"`
}

func (tagPrimaryModel) TableName() string { return "tag_primary" }

func (tag tagPrimaryModel) domain() domain.TagPrimary {
	secondaries := make(domain.TagSecondaryList, len(tag.TagSecondaryList))
	for index, secondary := range tag.TagSecondaryList {
		secondaries[index] = secondary.domain()
	}
	return domain.TagPrimary{ID: tag.ID, Name: tag.Name, TagSecondaryList: secondaries}
}

type tagSecondaryModel struct {
	ID           int    `gorm:"primaryKey;column:id"`
	TagPrimaryID int    `gorm:"column:primary_tag_id"`
	Name         string `gorm:"column:name"`
}

func (tagSecondaryModel) TableName() string { return "tag_secondary" }

func (tag tagSecondaryModel) domain() domain.TagSecondary {
	return domain.TagSecondary(tag)
}
