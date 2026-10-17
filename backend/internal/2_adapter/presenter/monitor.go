package presenter

import (
	"context"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// DistributeAgent は、ドメイン状態を画面配信ゲートウェイへ渡します。
func (receiver *Presenter) DistributeAgent(
	ctx context.Context,
	domainDomain *domain.Domain,
) {
	receiver.ToMonitor.DistributeAgent(
		ctx,
		domainDomain,
	)
}
