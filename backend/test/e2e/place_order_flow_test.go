//go:build e2e

package e2e

import (
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"
)

const (
	placeStatusCleaning = 0
	placeStatusOpening  = 1
	placeStatusUsing    = 2
	placeStatusBilling  = 3

	orderStatusPreparing = 1
	orderStatusPassed    = 2

	billingStageStockUpdated = 4
)

func TestPlaceOrderFlow(t *testing.T) {
	janCode, gateKey := requireBillingTestConfiguration(t)
	waitForBackend(t, 5*time.Second)
	resetPlaces(t)
	t.Cleanup(func() { resetPlaces(t) })

	places := getPlaces(t)
	placeNames := make([]string, 0, len(places))
	for name := range places {
		placeNames = append(placeNames, name)
	}
	sort.Strings(placeNames)
	if len(placeNames) == 0 {
		t.Fatal("no places are configured")
	}
	placeName := placeNames[0]
	if places[placeName] == nil || places[placeName].PlaceStatus != placeStatusOpening {
		t.Fatalf("initial place status = %#v, want Opening", places[placeName])
	}

	stream := openPlaceSSE(t, placeName)
	defer stream.close()
	initialEvent := stream.nextPlace(t)
	if initialEvent.PlaceStatus != placeStatusOpening {
		t.Fatalf("initial SSE status = %d, want %d", initialEvent.PlaceStatus, placeStatusOpening)
	}

	patchPlaceStatus(t, placeName, placeStatusUsing)
	usingEvent := stream.nextPlace(t)
	if usingEvent.PlaceStatus != placeStatusUsing {
		t.Fatalf("updated SSE status = %d, want %d", usingEvent.PlaceStatus, placeStatusUsing)
	}

	orderPath := "/v1/place/" + url.PathEscape(placeName) + "/order"
	requestAndExpectStatus(
		t,
		http.MethodPost,
		orderPath,
		map[string]int{"jan_code": janCode, "count": 1},
		http.StatusOK,
	)
	orderedEvent := stream.nextPlace(t)
	if len(orderedEvent.OrderList) != 1 {
		t.Fatalf("order count = %d, want 1", len(orderedEvent.OrderList))
	}
	order, exists := orderedEvent.OrderList[1]
	if !exists || order.JANCode != janCode || order.Count != 1 || order.OrderStatus != orderStatusPreparing {
		t.Fatalf("created order = %#v, want JAN=%d count=1 status=%d", order, janCode, orderStatusPreparing)
	}

	patchOrderStatus(t, placeName, order.OrderNumber, orderStatusPassed)
	passedEvent := stream.nextPlace(t)
	if got := passedEvent.OrderList[order.OrderNumber].OrderStatus; got != orderStatusPassed {
		t.Fatalf("passed order status = %d, want %d", got, orderStatusPassed)
	}

	patchPlaceStatus(t, placeName, placeStatusBilling)
	billingEvent := stream.nextPlace(t)
	if billingEvent.PlaceStatus != placeStatusBilling {
		t.Fatalf("billing place status = %d, want %d", billingEvent.PlaceStatus, placeStatusBilling)
	}
	if billingEvent.BillingStage != billingStageStockUpdated {
		t.Fatalf("billing stage = %d, want %d", billingEvent.BillingStage, billingStageStockUpdated)
	}
	if billingEvent.TransactionHeadID == "" {
		t.Fatal("transaction head ID is empty after billing integration")
	}

	completeBilling(t, billingEvent.TransactionHeadID, gateKey)
	cleaningEvent := stream.nextPlace(t)
	if cleaningEvent.PlaceStatus != placeStatusCleaning {
		t.Fatalf("cleaning place status = %d, want %d", cleaningEvent.PlaceStatus, placeStatusCleaning)
	}

	patchPlaceStatus(t, placeName, placeStatusOpening)
	openingEvent := stream.nextPlace(t)
	if openingEvent.PlaceStatus != placeStatusOpening {
		t.Fatalf("reopened place status = %d, want %d", openingEvent.PlaceStatus, placeStatusOpening)
	}
	if openingEvent.TransactionHeadID != "" || openingEvent.BillingStage != 0 || len(openingEvent.OrderList) != 0 {
		t.Fatalf("reopened place was not reset: %#v", openingEvent)
	}
	persistedPlace := getPlaces(t)[placeName]
	if persistedPlace == nil || persistedPlace.PlaceStatus != placeStatusOpening ||
		persistedPlace.TransactionHeadID != "" || persistedPlace.BillingStage != 0 || len(persistedPlace.OrderList) != 0 {
		t.Fatalf("persisted reopened place was not reset: %#v", persistedPlace)
	}
}
