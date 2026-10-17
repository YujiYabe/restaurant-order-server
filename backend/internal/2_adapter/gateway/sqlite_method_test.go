package gateway

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

type sqliteTransactionStub struct {
	tx             *gorm.DB
	endErr         error
	getProductCall int
}

func (stub *sqliteTransactionStub) WithOutTx() *gorm.DB { return stub.tx }
func (stub *sqliteTransactionStub) BeginTx() *gorm.DB   { return stub.tx }
func (stub *sqliteTransactionStub) EndTx(*gorm.DB, bool) error {
	return stub.endErr
}
func (stub *sqliteTransactionStub) GetProduct(*gorm.DB, int) (domain.Product, error) {
	stub.getProductCall++
	return domain.Product{JANCode: 100, Stock: 10}, nil
}
func (stub *sqliteTransactionStub) GetAllProductList(*gorm.DB) (domain.AllProductList, error) {
	return nil, nil
}
func (stub *sqliteTransactionStub) UpdateProduct(*gorm.DB, domain.Product, domain.Product) error {
	return nil
}
func (stub *sqliteTransactionStub) UpdateProductNoSaga(context.Context, domain.Product) error {
	return nil
}
func (stub *sqliteTransactionStub) GetAllTag(*gorm.DB) (domain.TagPrimaryList, error) {
	return nil, nil
}
func (stub *sqliteTransactionStub) UpdateStock(*gorm.DB, domain.Product, int) error { return nil }

func TestUpdateStockReturnsBeginTransactionError(t *testing.T) {
	wantErr := errors.New("begin failed")
	stub := &sqliteTransactionStub{tx: &gorm.DB{Error: wantErr}}
	receiver := &Gateway{ToSqlite: stub}

	err := receiver.UpdateStock(context.Background(), map[int]int{100: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("UpdateStock() error = %v, want %v", err, wantErr)
	}
}

func TestUpdateStockReturnsCommitError(t *testing.T) {
	wantErr := errors.New("commit failed")
	stub := &sqliteTransactionStub{tx: &gorm.DB{}, endErr: wantErr}
	receiver := &Gateway{ToSqlite: stub}

	err := receiver.UpdateStock(context.Background(), map[int]int{100: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("UpdateStock() error = %v, want %v", err, wantErr)
	}
}

func TestUpdateStockSkipsInternalChargeJANCode(t *testing.T) {
	stub := &sqliteTransactionStub{tx: &gorm.DB{}}
	receiver := &Gateway{ToSqlite: stub}

	if err := receiver.UpdateStock(context.Background(), map[int]int{0: 1}); err != nil {
		t.Fatalf("UpdateStock() error = %v", err)
	}
	if stub.getProductCall != 0 {
		t.Fatalf("GetProduct() calls = %d, want 0", stub.getProductCall)
	}
}
