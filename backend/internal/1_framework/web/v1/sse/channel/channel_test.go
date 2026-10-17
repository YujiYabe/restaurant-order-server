package ssechannel

import (
	"testing"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// TestPublishReplacesPendingSnapshot は、対象仕様の回帰動作を検証します。
func TestPublishReplacesPendingSnapshot(t *testing.T) {
	select {
	case <-DomainChan:
	default:
	}
	first := &domain.Domain{}
	latest := &domain.Domain{}
	Publish(first)
	Publish(latest)

	if got := <-DomainChan; got != latest {
		t.Fatal("Publish() did not retain the latest snapshot")
	}
}
