package smaregi

import (
	"testing"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

func TestLocalTransactionResponse(t *testing.T) {
	t.Parallel()
	request := domain.ReqPostTransactionTemporary{CustomerId: 12, TerminalTranId: 34}

	response := localTransactionResponse("shop", request)
	if response.CustomerId != request.CustomerId {
		t.Fatalf("customer ID = %d, want %d", response.CustomerId, request.CustomerId)
	}
	if response.TransactionHeadId != "local-shop-12-34" {
		t.Fatalf("transaction head ID = %q, want %q", response.TransactionHeadId, "local-shop-12-34")
	}
}
