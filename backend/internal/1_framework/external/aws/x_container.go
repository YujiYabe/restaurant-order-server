package aws

import "github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"

// Aws は、AWS Rekognitionとの接続を管理します。
type Aws struct{}

// NewToAws は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToAws() (
	toAws gateway.ToAws,
) {
	toAws = new(Aws)
	return
}
