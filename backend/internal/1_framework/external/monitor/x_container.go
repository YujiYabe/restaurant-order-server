package monitor

import (
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/presenter"
)

// Monitor は、画面への状態配信を管理します。
type Monitor struct{}

// NewMonitor は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewMonitor() (
	monitor *Monitor,
) {
	monitor = &Monitor{}
	return
}

// NewToMonitor は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToMonitor() (
	toMonitor presenter.ToMonitor,
) {
	toMonitor = new(Monitor)
	return
}
