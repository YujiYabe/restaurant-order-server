// Package transport は、JSONおよびYAML境界で使用する表現を定義します。
// ドメインモデルには、意図的にシリアライズ用タグを持たせていません。
package transport

import (
	"encoding/json"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

func MarshalJSON(value any) ([]byte, error) { return json.Marshal(JSONValue(value)) }

type Order struct {
	OrderNumber int              `json:"order_number" yaml:"order_number"`
	OrderStatus domain.StatusNum `json:"order_status" yaml:"order_status"`
	JANCode     int              `json:"jan_code" yaml:"jan_code"`
	Count       int              `json:"count" yaml:"count"`
	OrderTime   time.Time        `json:"order_time" yaml:"order_time"`
	OfferTime   time.Time        `json:"offer_time" yaml:"offer_time"`
}

func (v Order) Domain() domain.Order {
	return domain.Order(v)
}

func OrderFromDomain(v domain.Order) Order { return Order(v) }

type Alcohol struct {
	UUID            string `json:"uuid" yaml:"uuid"`
	Base64Source    string `json:"base64Source" yaml:"base64Source"`
	CertificateType string `json:"certificate_type" yaml:"certificate_type"`
}

func (v Alcohol) Domain() domain.Alcohol { return domain.Alcohol(v) }

type OrderAlcohol struct {
	Order
	Alcohol
}

type Place struct {
	PlaceGroup        string              `json:"place_group" yaml:"place_group"`
	PlaceName         string              `json:"place_name" yaml:"place_name"`
	Capacity          int                 `json:"capacity" yaml:"capacity"`
	CustomerId        int                 `json:"customer_id" yaml:"customer_id"`
	LineNumber        int                 `json:"line_number" yaml:"line_number"`
	TransactionHeadId string              `json:"transaction_head_id" yaml:"transaction_head_id"`
	BillingStage      domain.BillingStage `json:"billing_stage" yaml:"billing_stage"`
	PlaceStatus       domain.StatusNum    `json:"place_status" yaml:"place_status"`
	EnterTime         time.Time           `json:"enter_time" yaml:"enter_time"`
	ExitTime          time.Time           `json:"exit_time" yaml:"exit_time"`
	Allergy           string              `json:"allergy" yaml:"allergy"`
	Language          int                 `json:"language" yaml:"language"`
	Permission        string              `json:"permission" yaml:"permission"`
	BaseCount         int                 `json:"base_count" yaml:"base_count"`
	ChargePrice       int                 `json:"charge_price" yaml:"charge_price"`
	OrderPrice        int                 `json:"order_price" yaml:"order_price"`
	OrderList         map[int]Order       `json:"order_list" yaml:"order_list"`
}

func (v Place) Domain() domain.Place {
	orders := make(domain.OrderList, len(v.OrderList))
	for key, order := range v.OrderList {
		orders[key] = order.Domain()
	}
	return domain.Place{
		PlaceGroup: v.PlaceGroup, PlaceName: v.PlaceName, Capacity: v.Capacity,
		CustomerId: v.CustomerId, LineNumber: v.LineNumber,
		TransactionHeadId: v.TransactionHeadId, BillingStage: v.BillingStage,
		PlaceStatus: v.PlaceStatus, EnterTime: v.EnterTime, ExitTime: v.ExitTime,
		Allergy: v.Allergy, Language: v.Language, Permission: v.Permission,
		BaseCount: v.BaseCount, ChargePrice: v.ChargePrice, OrderPrice: v.OrderPrice,
		OrderList: orders,
	}
}

func PlaceFromDomain(v domain.Place) Place {
	orders := make(map[int]Order, len(v.OrderList))
	for key, order := range v.OrderList {
		orders[key] = OrderFromDomain(order)
	}
	return Place{
		PlaceGroup: v.PlaceGroup, PlaceName: v.PlaceName, Capacity: v.Capacity,
		CustomerId: v.CustomerId, LineNumber: v.LineNumber,
		TransactionHeadId: v.TransactionHeadId, BillingStage: v.BillingStage,
		PlaceStatus: v.PlaceStatus, EnterTime: v.EnterTime, ExitTime: v.ExitTime,
		Allergy: v.Allergy, Language: v.Language, Permission: v.Permission,
		BaseCount: v.BaseCount, ChargePrice: v.ChargePrice, OrderPrice: v.OrderPrice,
		OrderList: orders,
	}
}

type PlaceList map[string]*Place

func (v PlaceList) Domain() domain.PlaceList {
	result := make(domain.PlaceList, len(v))
	for key, place := range v {
		if place != nil {
			converted := place.Domain()
			result[key] = &converted
		}
	}
	return result
}

func PlaceListFromDomain(v domain.PlaceList) PlaceList {
	result := make(PlaceList, len(v))
	for key, place := range v {
		if place != nil {
			converted := PlaceFromDomain(*place)
			result[key] = &converted
		}
	}
	return result
}

type Product struct {
	JANCode       int    `json:"jan_code"`
	ProductId     int    `json:"product_id"`
	IsValid       bool   `json:"is_valid"`
	Stock         int    `json:"stock"`
	Place         string `json:"place"`
	Calorie       int    `json:"calorie"`
	CostPrice     int    `json:"cost_price"`
	RetailPrice   int    `json:"retail_price"`
	OperationCost int    `json:"operation_cost"`
	EnergyCost    int    `json:"energy_cost"`
	Recipe        string `json:"recipe"`
	Tag           string `json:"tag"`
	Allergy       string `json:"allergy"`
	NameJa        string `json:"name_ja"`
	NameJaShort   string `json:"name_ja_short"`
	InformationJa string `json:"information_ja"`
	UpdatedAt     string `json:"updated_at"`
}

func (v Product) Domain() domain.Product         { return domain.Product(v) }
func ProductFromDomain(v domain.Product) Product { return Product(v) }

type MiniProduct struct {
	JANCode       int    `json:"jan_code"`
	ProductId     int    `json:"product_id"`
	Price         int    `json:"price"`
	Place         string `json:"place"`
	Recipe        string `json:"recipe"`
	Tag           string `json:"tag"`
	NameJa        string `json:"name_ja"`
	Allergy       string `json:"allergy"`
	InformationJa string `json:"information_ja"`
}

type TagSecondary struct {
	ID           int    `json:"id"`
	TagPrimaryID int    `json:"primary_tag_id"`
	Name         string `json:"name"`
}
type TagPrimary struct {
	ID               int            `json:"id"`
	Name             string         `json:"name"`
	TagSecondaryList []TagSecondary `json:"tag_secondary_list"`
}

// JSONValue は、HTTPハンドラーが返す主なドメイン値を通信用DTOへ変換します。
func JSONValue(value any) any {
	switch v := value.(type) {
	case domain.PlaceList:
		return PlaceListFromDomain(v)
	case *domain.Place:
		if v == nil {
			return (*Place)(nil)
		}
		converted := PlaceFromDomain(*v)
		return &converted
	case domain.Place:
		return PlaceFromDomain(v)
	case domain.AllProductList:
		result := make([]Product, len(v))
		for i := range v {
			result[i] = ProductFromDomain(v[i])
		}
		return result
	case domain.MiniProductMap:
		result := make(map[int]MiniProduct, len(v))
		for key, item := range v {
			result[key] = MiniProduct(item)
		}
		return result
	case domain.TagPrimaryList:
		result := make([]TagPrimary, len(v))
		for i, primary := range v {
			secondaries := make([]TagSecondary, len(primary.TagSecondaryList))
			for j, secondary := range primary.TagSecondaryList {
				secondaries[j] = TagSecondary(secondary)
			}
			result[i] = TagPrimary{ID: primary.ID, Name: primary.Name, TagSecondaryList: secondaries}
		}
		return result
	case domain.OrderList:
		result := make(map[int]Order, len(v))
		for key, order := range v {
			result[key] = OrderFromDomain(order)
		}
		return result
	default:
		return value
	}
}
