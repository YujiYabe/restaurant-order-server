package gate

import (
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
)

// Gate は、会計連携サービスとの接続を管理します。
type Gate struct {
}

// NewToGate は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToGate() (
	toGate gateway.ToGate,
) {
	toGate = new(Gate)
	return
}
