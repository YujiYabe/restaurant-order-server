package domain

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestAdjustmentInventoryRejectsNonPositiveCount は、対象仕様の回帰動作を検証します。
func TestAdjustmentInventoryRejectsNonPositiveCount(t *testing.T) {
	domain := &Domain{InventoryMap: InventoryMap{100: 5}}

	for _, count := range []int{0, -1} {
		if domain.AdjustmentInventory(context.Background(), Order{JANCode: 100, Count: count}) {
			t.Fatalf("AdjustmentInventory() accepted count %d", count)
		}
		if got := domain.InventoryMap[100]; got != 5 {
			t.Fatalf("inventory = %d after count %d, want 5", got, count)
		}
	}
}

// TestTryAddOrderIsAtomic は、対象仕様の回帰動作を検証します。
func TestTryAddOrderIsAtomic(t *testing.T) {
	domain := &Domain{
		PlaceList: PlaceList{
			"A1": {
				PlaceName:   "A1",
				PlaceStatus: PlaceStatusUsing,
				OrderList:   NewOrderList(),
			},
		},
		InventoryMap:   InventoryMap{100: 5},
		MiniProductMap: MiniProductMap{100: {JANCode: 100, Price: 500}},
	}

	var successes atomic.Int32
	var waitGroup sync.WaitGroup
	for range 20 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if domain.TryAddOrder(context.Background(), "A1", Order{JANCode: 100, Count: 1}) {
				successes.Add(1)
			}
		}()
	}
	waitGroup.Wait()

	if got := successes.Load(); got != 5 {
		t.Fatalf("successful orders = %d, want 5", got)
	}
	if got := domain.GetInventoryMap(context.Background())[100]; got != 0 {
		t.Fatalf("inventory = %d, want 0", got)
	}
	if got := len(domain.GetPlaceList(context.Background())["A1"].OrderList); got != 5 {
		t.Fatalf("order count = %d, want 5", got)
	}
}

func TestTryAddOrderRejectsInternalChargeJANCode(t *testing.T) {
	receiver := &Domain{
		PlaceList: PlaceList{"A1": {
			PlaceStatus: PlaceStatusUsing,
			OrderList:   NewOrderList(),
		}},
		InventoryMap:   NewInventoryMap(),
		MiniProductMap: NewMiniProductMap(),
	}

	if receiver.TryAddOrder(context.Background(), "A1", Order{JANCode: 0, Count: 1}) {
		t.Fatal("TryAddOrder() accepted internal charge JAN code")
	}
	if got := len(receiver.GetPlaceList(context.Background())["A1"].OrderList); got != 0 {
		t.Fatalf("order count = %d, want 0", got)
	}
}

func TestBuildPostTransactionTemporaryExcludesChargeFromStockUpdate(t *testing.T) {
	places := PlaceList{"A1": {
		PlaceName: "A1",
		OrderList: OrderList{
			1: {JANCode: 0, Count: 1, OrderStatus: OrderStatusPassed},
			2: {JANCode: 100, Count: 2, OrderStatus: OrderStatusPreparing},
		},
	}}

	_, counts := places.BuildPostTransactionTemporary(
		context.Background(),
		"A1",
		MiniProductMap{100: {JANCode: 100, ProductId: 10, Price: 500}},
	)
	if _, exists := counts[0]; exists {
		t.Fatal("stock update contains internal charge JAN code")
	}
	if counts[100] != 2 {
		t.Fatalf("product count = %d, want 2", counts[100])
	}
}

// TestGetPlaceListReturnsSnapshot は、対象仕様の回帰動作を検証します。
func TestGetPlaceListReturnsSnapshot(t *testing.T) {
	domain := &Domain{PlaceList: PlaceList{
		"A1": {PlaceName: "A1", PlaceStatus: PlaceStatusUsing, OrderList: NewOrderList()},
	}}

	snapshot := domain.GetPlaceList(context.Background())
	snapshot["A1"].PlaceStatus = PlaceStatusBilling

	if got := domain.GetPlace(context.Background(), "A1").PlaceStatus; got != PlaceStatusUsing {
		t.Fatalf("stored place status = %d, want %d", got, PlaceStatusUsing)
	}
}

func TestGetDomainReturnsDeepCopyOfReferenceFields(t *testing.T) {
	receiver := &Domain{
		AllergyDefault: Allergy{"egg": "卵"},
		TagPrimaryList: TagPrimaryList{{
			ID:               1,
			TagSecondaryList: TagSecondaryList{{ID: 2, Name: "hot"}},
		}},
	}

	snapshot := receiver.GetDomain(context.Background())
	snapshot.AllergyDefault["egg"] = "changed"
	snapshot.TagPrimaryList[0].TagSecondaryList[0].Name = "changed"

	if got := receiver.GetAllergyDefault(context.Background())["egg"]; got != "卵" {
		t.Fatalf("allergy changed through snapshot: %q", got)
	}
	if got := receiver.GetTagList(context.Background())[0].TagSecondaryList[0].Name; got != "hot" {
		t.Fatalf("tag changed through snapshot: %q", got)
	}
}

// TestEditOrderValidation は、対象仕様の回帰動作を検証します。
func TestEditOrderValidation(t *testing.T) {
	tests := []struct {
		name        string
		orders      OrderList
		target      Order
		wantUpdated bool
		wantStatus  StatusNum
	}{
		{
			name: "valid transition",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusPreparing, JANCode: 100, Count: 1},
			},
			target:      Order{OrderNumber: 1, OrderStatus: OrderStatusPassed},
			wantUpdated: true,
			wantStatus:  OrderStatusPassed,
		},
		{
			name: "missing order",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusPreparing, JANCode: 100, Count: 1},
			},
			target:     Order{OrderNumber: 2, OrderStatus: OrderStatusPassed},
			wantStatus: OrderStatusPreparing,
		},
		{
			name: "invalid status",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusPreparing, JANCode: 100, Count: 1},
			},
			target:     Order{OrderNumber: 1, OrderStatus: StatusNum(99)},
			wantStatus: OrderStatusPreparing,
		},
		{
			name: "completed order can be canceled",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusPassed, JANCode: 100, Count: 1},
			},
			target:      Order{OrderNumber: 1, OrderStatus: OrderStatusCancel},
			wantUpdated: true,
			wantStatus:  OrderStatusCancel,
		},
		{
			name: "canceled order can return to preparing",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusCancel, JANCode: 100, Count: 1},
			},
			target:      Order{OrderNumber: 1, OrderStatus: OrderStatusPreparing},
			wantUpdated: true,
			wantStatus:  OrderStatusPreparing,
		},
		{
			name: "same status",
			orders: OrderList{
				1: {OrderNumber: 1, OrderStatus: OrderStatusPreparing, JANCode: 100, Count: 1},
			},
			target:     Order{OrderNumber: 1, OrderStatus: OrderStatusPreparing},
			wantStatus: OrderStatusPreparing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeList := PlaceList{
				"A1": {PlaceName: "A1", OrderList: tt.orders},
			}
			products := MiniProductMap{100: {JANCode: 100, Price: 500}}

			updated := placeList.EditOrder(context.Background(), "A1", tt.target, products)
			if updated != tt.wantUpdated {
				t.Fatalf("EditOrder() = %v, want %v", updated, tt.wantUpdated)
			}
			if got := placeList["A1"].OrderList[1].OrderStatus; got != tt.wantStatus {
				t.Fatalf("status = %d, want %d", got, tt.wantStatus)
			}
			if got := len(placeList["A1"].OrderList); got != 1 {
				t.Fatalf("order count = %d, want 1", got)
			}
		})
	}
}

// TestEditOrderCancellationRestoresInventory は、対象仕様の回帰動作を検証します。
func TestEditOrderCancellationRestoresInventory(t *testing.T) {
	domain := &Domain{
		PlaceList: PlaceList{"A1": {
			PlaceName: "A1",
			OrderList: OrderList{1: {OrderNumber: 1, OrderStatus: OrderStatusPreparing, JANCode: 100, Count: 2}},
		}},
		InventoryMap:   InventoryMap{100: 3},
		MiniProductMap: MiniProductMap{100: {JANCode: 100, Price: 500}},
	}

	if !domain.EditOrder(context.Background(), "A1", Order{OrderNumber: 1, OrderStatus: OrderStatusCancel}, domain.MiniProductMap) {
		t.Fatal("EditOrder() did not cancel the order")
	}
	if got := domain.GetInventoryMap(context.Background())[100]; got != 5 {
		t.Fatalf("inventory = %d, want 5", got)
	}
	if domain.EditOrder(context.Background(), "A1", Order{OrderNumber: 1, OrderStatus: OrderStatusCancel}, domain.MiniProductMap) {
		t.Fatal("EditOrder() canceled an already canceled order")
	}
	if got := domain.GetInventoryMap(context.Background())[100]; got != 5 {
		t.Fatalf("inventory after duplicate cancellation = %d, want 5", got)
	}
	if !domain.EditOrder(context.Background(), "A1", Order{OrderNumber: 1, OrderStatus: OrderStatusPreparing}, domain.MiniProductMap) {
		t.Fatal("EditOrder() did not restore the canceled order")
	}
	if got := domain.GetInventoryMap(context.Background())[100]; got != 3 {
		t.Fatalf("inventory after restoring cancellation = %d, want 3", got)
	}
}

func TestEditOrderDoesNotRestoreCanceledOrderWithoutStock(t *testing.T) {
	domain := &Domain{
		PlaceList: PlaceList{"A1": {
			PlaceName: "A1",
			OrderList: OrderList{1: {OrderNumber: 1, OrderStatus: OrderStatusCancel, JANCode: 100, Count: 2}},
		}},
		InventoryMap:   InventoryMap{100: 1},
		MiniProductMap: MiniProductMap{100: {JANCode: 100, Price: 500}},
	}

	if domain.EditOrder(context.Background(), "A1", Order{OrderNumber: 1, OrderStatus: OrderStatusPreparing}, domain.MiniProductMap) {
		t.Fatal("EditOrder() restored an order without enough stock")
	}
	if got := domain.GetInventoryMap(context.Background())[100]; got != 1 {
		t.Fatalf("inventory = %d, want 1", got)
	}
	if got := domain.GetPlace(context.Background(), "A1").OrderList[1].OrderStatus; got != OrderStatusCancel {
		t.Fatalf("status = %d, want %d", got, OrderStatusCancel)
	}
}

// TestDomainAutoUpdateStopsWhenContextIsCanceled は、対象仕様の回帰動作を検証します。
func TestDomainAutoUpdateStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	receiver := &Domain{
		PlaceList: PlaceList{
			"A1": {PlaceName: "A1", PlaceStatus: PlaceStatusOpening},
		},
		MiniProductMap: MiniProductMap{},
	}
	done := make(chan struct{})
	go func() {
		receiver.autoUpdate(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("AutoUpdate() did not stop after context cancellation")
	}
}

// TestVerifyTwentyRejectsMalformedBirthday は、対象仕様の回帰動作を検証します。
func TestVerifyTwentyRejectsMalformedBirthday(t *testing.T) {
	domain := &Domain{}
	malformed := []string{
		"昭和54年",
		"昭和54年9月",
		"昭和年月日",
	}

	for _, birthday := range malformed {
		t.Run(birthday, func(t *testing.T) {
			if _, _, err := domain.VerifyTwenty(context.Background(), birthday); err == nil {
				t.Fatal("VerifyTwenty() returned no error")
			}
		})
	}
}

// TestTwentyListExpiration は、対象仕様の回帰動作を検証します。
func TestTwentyListExpiration(t *testing.T) {
	domain := &Domain{AllowedTwentyList: NewTwentyList()}
	valid := Alcohol{UUID: "valid"}
	domain.AddTwentyList(context.Background(), valid)

	if !domain.IsIncludeTwentyList(context.Background(), valid) {
		t.Fatal("newly verified UUID was not accepted")
	}

	domain.mutex.Lock()
	domain.AllowedTwentyList["expired"] = time.Now().Add(-time.Second)
	domain.mutex.Unlock()
	if domain.IsIncludeTwentyList(context.Background(), Alcohol{UUID: "expired"}) {
		t.Fatal("expired UUID was accepted")
	}

	domain.mutex.RLock()
	_, remains := domain.AllowedTwentyList["expired"]
	domain.mutex.RUnlock()
	if remains {
		t.Fatal("expired UUID was not removed")
	}
}

// TestFaceRegistrationPermissionIsOneTime は、対象仕様の回帰動作を検証します。
func TestFaceRegistrationPermissionIsOneTime(t *testing.T) {
	domain := &Domain{PendingFaceList: NewTwentyList()}
	alcohol := Alcohol{UUID: "verified-document"}

	if domain.ConsumeFaceRegistration(context.Background(), alcohol) {
		t.Fatal("unapproved face registration was accepted")
	}
	domain.AllowFaceRegistration(context.Background(), alcohol)
	if !domain.ConsumeFaceRegistration(context.Background(), alcohol) {
		t.Fatal("approved face registration was rejected")
	}
	if domain.ConsumeFaceRegistration(context.Background(), alcohol) {
		t.Fatal("face registration permission was reused")
	}
}

// TestSaveInMemorySubtractsUnbilledOrders は、対象仕様の回帰動作を検証します。
func TestSaveInMemorySubtractsUnbilledOrders(t *testing.T) {
	domain := &Domain{}
	products := AllProductList{{JANCode: 100, Stock: 10, IsValid: true}}
	places := PlaceList{
		"A1": {PlaceStatus: PlaceStatusUsing, OrderList: OrderList{1: {JANCode: 100, Count: 3, OrderStatus: OrderStatusPreparing}}},
		"A2": {PlaceStatus: PlaceStatusBilling, BillingStage: BillingStageTransactionInfoSent, OrderList: OrderList{1: {JANCode: 100, Count: 2, OrderStatus: OrderStatusPreparing}}},
		"A3": {PlaceStatus: PlaceStatusBilling, BillingStage: BillingStageStockUpdated, OrderList: OrderList{1: {JANCode: 100, Count: 4, OrderStatus: OrderStatusPreparing}}},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	domain.SaveInMemory(ctx, products, nil, places)
	if got := domain.GetInventoryMap(context.Background())[100]; got != 5 {
		t.Fatalf("inventory = %d, want 5", got)
	}
}

// TestEditPlaceRejectsInvalidTransition は、対象仕様の回帰動作を検証します。
func TestEditPlaceRejectsInvalidTransition(t *testing.T) {
	domain := &Domain{PlaceList: PlaceList{"A1": {PlaceName: "A1", PlaceStatus: PlaceStatusOpening}}}
	if domain.EditPlace(context.Background(), Place{PlaceName: "A1", PlaceStatus: PlaceStatusBilling}) {
		t.Fatal("opening to billing transition was accepted")
	}
	if got := domain.GetPlace(context.Background(), "A1").PlaceStatus; got != PlaceStatusOpening {
		t.Fatalf("place status = %d, want %d", got, PlaceStatusOpening)
	}
}
