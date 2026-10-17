package gateway

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *Gateway) GetAllProductList(
	ctx context.Context,
) (
	allProductList domain.AllProductList,
	err error,
) {
	allProductList, err = receiver.ToSqlite.GetAllProductList(
		receiver.ToSqlite.WithOutTx(),
	)

	return
}

// UpdateProduct は、SQLiteトランザクション内で対象商品を取得し、更新内容を保存します。
func (receiver *Gateway) UpdateProduct(
	ctx context.Context,
	newProduct domain.Product,
) (
	err error,
) {
	tx := receiver.ToSqlite.BeginTx()
	if tx == nil {
		return fmt.Errorf("begin product update transaction: nil transaction")
	}
	if tx.Error != nil {
		return fmt.Errorf("begin product update transaction: %w", tx.Error)
	}
	defer func() {
		if endErr := receiver.ToSqlite.EndTx(tx, err == nil); endErr != nil {
			err = errors.Join(err, fmt.Errorf("finish product update transaction: %w", endErr))
		}
	}()

	var product domain.Product

	product, err = receiver.ToSqlite.GetProduct(
		tx,
		newProduct.JANCode,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return err
	}

	err = receiver.ToSqlite.UpdateProduct(
		tx,
		product,
		newProduct,
	)
	if err != nil {
		pkg.Logging(ctx, err)
	}

	return err
}

// UpdateStock は、会計済みの商品数量をSQLite在庫から減算します。
func (receiver *Gateway) UpdateStock(
	ctx context.Context,
	janCodeCountMap map[int]int,
) (
	err error,
) {
	tx := receiver.ToSqlite.BeginTx()
	if tx == nil {
		return fmt.Errorf("begin stock update transaction: nil transaction")
	}
	if tx.Error != nil {
		return fmt.Errorf("begin stock update transaction: %w", tx.Error)
	}
	defer func() {
		if endErr := receiver.ToSqlite.EndTx(tx, err == nil); endErr != nil {
			err = errors.Join(err, fmt.Errorf("finish stock update transaction: %w", endErr))
		}
	}()

	var product domain.Product
	for janCode, count := range janCodeCountMap {
		if janCode == 0 {
			continue
		}
		product, err = receiver.ToSqlite.GetProduct(tx, janCode)
		if err != nil {
			pkg.Logging(ctx, err)
			return err
		}
		err = receiver.ToSqlite.UpdateStock(tx, product, count)
		if err != nil {
			pkg.Logging(ctx, err)
			return err
		}
	}

	return err
}

// GetAllTag は、関連する副タグを含む全タグをSQLiteから取得します。
func (receiver *Gateway) GetAllTag(
	ctx context.Context,
) (
	tagPrimaryList domain.TagPrimaryList,
	err error,
) {
	tagPrimaryList, err = receiver.ToSqlite.GetAllTag(
		receiver.ToSqlite.WithOutTx(),
	)

	return
}
