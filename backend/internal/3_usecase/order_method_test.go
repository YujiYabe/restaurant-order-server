package usecase

import (
	"context"
	"errors"
	"testing"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

type orderGatewayStub struct {
	ToGateway
	billingErr        error
	resetErr          error
	getPlaceListErr   error
	placeListResponse domain.PlaceList
	distributeSet     int
	temporaryCalls    int
	saveOrderCalls    int
	infoCalls         int
	stockCalls        int
	saveOrderFailures int
	setPlaceListErr   error
}

// SendPostTransactionTemporary は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) SendPostTransactionTemporary(
	context.Context,
	domain.ReqPostTransactionTemporary,
) (domain.SmaregiPostResponseEnquiry, error) {
	stub.temporaryCalls++
	if stub.billingErr != nil {
		return domain.SmaregiPostResponseEnquiry{}, stub.billingErr
	}
	return domain.SmaregiPostResponseEnquiry{CustomerId: 1, TransactionHeadId: "tx-1"}, nil
}

// SaveOrder は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) SaveOrder(
	context.Context,
	domain.Place,
	domain.ReqPostTransactionTemporary,
) (
	err error,
) {
	stub.saveOrderCalls++
	if stub.saveOrderCalls <= stub.saveOrderFailures {
		return errors.New("save order failed")
	}
	return nil
}

// SendPostTransactionInfo は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) SendPostTransactionInfo(
	context.Context,
	string,
) (
	err error,
) {
	stub.infoCalls++
	return nil
}

// UpdateStock は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) UpdateStock(
	context.Context,
	map[int]int,
) (
	err error,
) {
	stub.stockCalls++
	return nil
}

// SetPlaceListToRedis は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) SetPlaceListToRedis(
	context.Context,
	domain.PlaceList,
) (
	err error,
) {
	stub.distributeSet++
	return stub.setPlaceListErr
}

// ResetPlaceListInRedis は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) ResetPlaceListInRedis(
	context.Context,
) (
	err error,
) {
	return stub.resetErr
}

// GetPlaceListFromRedis は、テストで使用する依存処理を代替します。
func (stub *orderGatewayStub) GetPlaceListFromRedis(
	context.Context,
) (domain.PlaceList, error) {
	return stub.placeListResponse, stub.getPlaceListErr
}

type orderPresenterStub struct {
	distributeCount int
}

// DistributeAgent は、テストで使用する依存処理を代替します。
func (stub *orderPresenterStub) DistributeAgent(context.Context, *domain.Domain) {
	stub.distributeCount++
}

// TestAddOrderRejectsNonPositiveCount は、対象仕様の回帰動作を検証します。
func TestAddOrderRejectsNonPositiveCount(t *testing.T) {
	domainModel := newOrderTestDomain()
	gateway := &orderGatewayStub{}
	presenter := &orderPresenterStub{}
	useCase := &UseCase{ToDomain: domainModel, ToGateway: gateway, ToPresenter: presenter}

	for _, count := range []int{0, -1} {
		if useCase.AddOrder(context.Background(), "A1", domain.Order{JANCode: 100, Count: count}) {
			t.Fatalf("AddOrder() accepted count %d", count)
		}
	}

	if got := domainModel.InventoryMap[100]; got != 5 {
		t.Fatalf("inventory = %d, want 5", got)
	}
	if got := len(domainModel.PlaceList["A1"].OrderList); got != 0 {
		t.Fatalf("order count = %d, want 0", got)
	}
	if gateway.distributeSet != 0 || presenter.distributeCount != 0 {
		t.Fatal("rejected order was distributed")
	}
}

// TestAddOrderRejectsUnknownPlace は、対象仕様の回帰動作を検証します。
func TestAddOrderRejectsUnknownPlace(t *testing.T) {
	domainModel := newOrderTestDomain()
	useCase := &UseCase{
		ToDomain:    domainModel,
		ToGateway:   &orderGatewayStub{},
		ToPresenter: &orderPresenterStub{},
	}

	if useCase.AddOrder(context.Background(), "UNKNOWN", domain.Order{JANCode: 100, Count: 1}) {
		t.Fatal("AddOrder() accepted an unknown place")
	}
}

// TestEditPlaceDoesNotEnterBillingWhenIntegrationFails は、対象仕様の回帰動作を検証します。
func TestEditPlaceDoesNotEnterBillingWhenIntegrationFails(t *testing.T) {
	domainModel := newOrderTestDomain()
	billingErr := errors.New("billing unavailable")
	useCase := &UseCase{
		ToDomain:    domainModel,
		ToGateway:   &orderGatewayStub{billingErr: billingErr},
		ToPresenter: &orderPresenterStub{},
	}

	err := useCase.EditPlace(
		context.Background(),
		domain.Place{PlaceName: "A1", PlaceStatus: domain.PlaceStatusBilling},
	)
	if !errors.Is(err, billingErr) {
		t.Fatalf("EditPlace() error = %v, want %v", err, billingErr)
	}
	if got := domainModel.PlaceList["A1"].PlaceStatus; got != domain.PlaceStatusUsing {
		t.Fatalf("place status = %d, want %d", got, domain.PlaceStatusUsing)
	}
	if got := domainModel.PlaceList["A1"].TransactionHeadId; got != "" {
		t.Fatalf("transaction head ID = %q, want empty", got)
	}
}

// TestResetPlaceListPreservesStateWhenReloadFails は、対象仕様の回帰動作を検証します。
func TestResetPlaceListPreservesStateWhenReloadFails(t *testing.T) {
	domainModel := newOrderTestDomain()
	reloadErr := errors.New("redis unavailable")
	gateway := &orderGatewayStub{
		getPlaceListErr: reloadErr,
		placeListResponse: domain.PlaceList{
			"BROKEN": {PlaceName: "BROKEN"},
		},
	}
	presenter := &orderPresenterStub{}
	useCase := &UseCase{ToDomain: domainModel, ToGateway: gateway, ToPresenter: presenter}

	err := useCase.ResetPlaceList(context.Background())
	if !errors.Is(err, reloadErr) {
		t.Fatalf("ResetPlaceList() error = %v, want %v", err, reloadErr)
	}
	places := domainModel.GetPlaceList(context.Background())
	if places["A1"] == nil || places["BROKEN"] != nil {
		t.Fatalf("place list changed after reload failure: %#v", places)
	}
	if gateway.distributeSet != 0 || presenter.distributeCount != 0 {
		t.Fatal("failed reset was distributed")
	}
}

// TestBillingRetryResumesAfterCompletedStep は、対象仕様の回帰動作を検証します。
func TestBillingRetryResumesAfterCompletedStep(t *testing.T) {
	domainModel := newOrderTestDomain()
	gateway := &orderGatewayStub{saveOrderFailures: 1}
	useCase := &UseCase{ToDomain: domainModel, ToGateway: gateway, ToPresenter: &orderPresenterStub{}}
	request := domain.Place{PlaceName: "A1", PlaceStatus: domain.PlaceStatusBilling}

	if err := useCase.EditPlace(context.Background(), request); err == nil {
		t.Fatal("first EditPlace() returned no error")
	}
	if stage := domainModel.GetPlace(context.Background(), "A1").BillingStage; stage != domain.BillingStageTransactionCreated {
		t.Fatalf("billing stage after failure = %d, want %d", stage, domain.BillingStageTransactionCreated)
	}
	if err := useCase.EditPlace(context.Background(), request); err != nil {
		t.Fatalf("second EditPlace() error = %v", err)
	}
	if gateway.temporaryCalls != 1 {
		t.Fatalf("temporary transaction calls = %d, want 1", gateway.temporaryCalls)
	}
	if gateway.saveOrderCalls != 2 || gateway.infoCalls != 1 || gateway.stockCalls != 1 {
		t.Fatalf("downstream calls = save:%d info:%d stock:%d", gateway.saveOrderCalls, gateway.infoCalls, gateway.stockCalls)
	}
	place := domainModel.GetPlace(context.Background(), "A1")
	if place.PlaceStatus != domain.PlaceStatusBilling || place.BillingStage != domain.BillingStageStockUpdated {
		t.Fatalf("final place = %#v", place)
	}
}

// TestBillingStopsWhenCheckpointCannotBePersisted は、対象仕様の回帰動作を検証します。
func TestBillingStopsWhenCheckpointCannotBePersisted(t *testing.T) {
	domainModel := newOrderTestDomain()
	redisErr := errors.New("redis unavailable")
	gateway := &orderGatewayStub{setPlaceListErr: redisErr}
	useCase := &UseCase{ToDomain: domainModel, ToGateway: gateway, ToPresenter: &orderPresenterStub{}}

	err := useCase.EditPlace(context.Background(), domain.Place{PlaceName: "A1", PlaceStatus: domain.PlaceStatusBilling})
	if !errors.Is(err, redisErr) {
		t.Fatalf("EditPlace() error = %v, want %v", err, redisErr)
	}
	if gateway.temporaryCalls != 1 || gateway.saveOrderCalls != 0 {
		t.Fatalf("calls after checkpoint failure = temporary:%d save:%d", gateway.temporaryCalls, gateway.saveOrderCalls)
	}
	if got := domainModel.GetPlace(context.Background(), "A1").PlaceStatus; got != domain.PlaceStatusUsing {
		t.Fatalf("place status = %d, want using", got)
	}
}

// newOrderTestDomain は、テストで使用する依存処理を代替します。
func newOrderTestDomain() *domain.Domain {
	return &domain.Domain{
		PlaceList: domain.PlaceList{
			"A1": {
				PlaceName:   "A1",
				PlaceStatus: domain.PlaceStatusUsing,
				CustomerId:  1,
				OrderList:   domain.NewOrderList(),
			},
		},
		InventoryMap: domain.InventoryMap{100: 5},
		MiniProductMap: domain.MiniProductMap{
			100: {JANCode: 100, ProductId: 100, Price: 500},
		},
	}
}
