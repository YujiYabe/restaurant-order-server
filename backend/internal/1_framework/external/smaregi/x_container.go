package smaregi

import (
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
	"time"
)

// Smaregi は、スマレジとの接続を管理します。
type Smaregi struct {
	accessToken        string
	accessTokenExpired time.Time
}

// NewToSmaregi は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToSmaregi() (
	toSmaregi gateway.ToSmaregi,
) {
	toSmaregi = new(Smaregi)
	return
}
