package azure

import "github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"

// Azure は、Azure Document Intelligenceとの接続を管理します。
type Azure struct{}

// NewToAzure は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToAzure() (
	toAzure gateway.ToAzure,
) {
	toAzure = new(Azure)
	return
}
