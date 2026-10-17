package monitor

import (
	"context"

	ssechannel "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/sse/channel"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// DistributeAgent は、最新のドメイン状態をSSE配信キューへ登録します。
func (receiver *Monitor) DistributeAgent(
	ctx context.Context,
	domainDomain *domain.Domain,
) {
	// 最新のスナップショットだけを保持し、
	// 配信遅延で注文・会計を止めない。
	ssechannel.Publish(domainDomain)
}
