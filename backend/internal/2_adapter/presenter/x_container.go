package presenter

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// Presenter は、ユースケースの出力を画面配信用に変換します。
type Presenter struct {
	ToMonitor ToMonitor
}

// ToMonitor は、モニターへ公開する処理を定義します。
type ToMonitor interface {
	DistributeAgent(
		ctx context.Context,
		domainDomain *domain.Domain,
	)
}

type NewPresenterOption struct {
	ToMonitor
}

// NewPresenter は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewPresenter(
	newPresenterOption NewPresenterOption,
) (
	presenter *Presenter,
) {
	presenter = &Presenter{
		ToMonitor: newPresenterOption.ToMonitor,
	}

	return
}
