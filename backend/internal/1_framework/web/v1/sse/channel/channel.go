package ssechannel

import (
	"sync"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

var DomainChan = make(chan *domain.Domain, 1)
var publishMu sync.Mutex

// Publish は、未配信の古いスナップショットを置き換えて最新状態を配信キューへ格納します。
func Publish(
	domainDomain *domain.Domain,
) {
	publishMu.Lock()
	defer publishMu.Unlock()
	select {
	case DomainChan <- domainDomain:
		return
	default:
	}
	select {
	case <-DomainChan:
	default:
	}
	DomainChan <- domainDomain
}
