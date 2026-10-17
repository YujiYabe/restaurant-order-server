//go:build e2e

package e2e

type placeResponse struct {
	PlaceName         string                `json:"place_name"`
	TransactionHeadID string                `json:"transaction_head_id"`
	BillingStage      int                   `json:"billing_stage"`
	PlaceStatus       int                   `json:"place_status"`
	OrderList         map[int]orderResponse `json:"order_list"`
}

type orderResponse struct {
	OrderNumber int `json:"order_number"`
	OrderStatus int `json:"order_status"`
	JANCode     int `json:"jan_code"`
	Count       int `json:"count"`
}
