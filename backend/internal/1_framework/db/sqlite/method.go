package sqlite

import (
	"context"

	"gorm.io/gorm"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// WithOutTx は、明示的なトランザクションを開始しないDB接続を返します。
func (receiver *Sqlite) WithOutTx() (
	tx *gorm.DB,
) {
	return receiver.Conn
}

// BeginTx は、SQLiteトランザクションを開始します。
func (receiver *Sqlite) BeginTx() (
	tx *gorm.DB,
) {
	return receiver.Conn.Begin()
}

// EndTx は、処理結果に応じてSQLiteトランザクションをコミットまたはロールバックします。
func (receiver *Sqlite) EndTx(
	tx *gorm.DB,
	isSuccess bool,
) error {
	if isSuccess {
		return tx.Commit().Error
	}
	return tx.Rollback().Error
}

// UpdateProductNoSaga は、既存値と更新値を比較して商品情報を更新します。
func (receiver *Sqlite) UpdateProductNoSaga(
	ctx context.Context,
	newProduct domain.Product,
) (
	err error,
) {
	err = receiver.Conn.Transaction(func(tx *gorm.DB) (
		err error,
	) {
		product, err := receiver.GetProduct(tx, newProduct.JANCode)
		if err != nil {
			return err
		}

		err = receiver.UpdateProduct(tx, product, newProduct)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		pkg.Logging(ctx, err)
	}

	return err
} //

// GetProduct は、指定JANコードの商品を返します。
func (receiver *Sqlite) GetProduct(
	tx *gorm.DB,
	janCode int,
) (
	product domain.Product,
	err error,
) {
	if tx == nil {
		tx = receiver.Conn
	}

	var stored productModel
	err = tx.
		// Debug().
		Table("product").
		Where("jan_code = ?", janCode).
		First(&stored).Error
	if err == nil {
		product = stored.domain()
	}

	return
}

// UpdateProduct は、現在の商品に指定された更新値を適用し、SQLiteへ保存します。
func (receiver *Sqlite) UpdateProduct(
	tx *gorm.DB,
	product domain.Product,
	newProduct domain.Product,
) (
	err error,
) {
	product.Place = newProduct.Place
	product.IsValid = newProduct.IsValid
	product.Stock = newProduct.Stock
	product.CostPrice = newProduct.CostPrice
	product.RetailPrice = newProduct.RetailPrice
	product.OperationCost = newProduct.OperationCost
	product.EnergyCost = newProduct.EnergyCost

	stored := productModelFromDomain(product)
	return tx.
		// Debug().
		Table("product").
		Where("jan_code = ?", newProduct.JANCode).
		Save(stored).
		Error
}

// UpdateStock は、会計済みの商品数量をSQLite在庫から減算します。
func (receiver *Sqlite) UpdateStock(
	tx *gorm.DB,
	product domain.Product,
	count int,
) (
	err error,
) {
	product.Stock -= count

	stored := productModelFromDomain(product)
	return tx.
		// Debug().
		Table("product").
		Where("jan_code = ?", product.JANCode).
		Save(stored).
		Error
}

// GetAllTag は、関連する副タグを含む全タグをSQLiteから取得します。
func (receiver *Sqlite) GetAllTag(
	tx *gorm.DB,
) (
	tagPrimaryList domain.TagPrimaryList,
	err error,
) {
	var stored []tagPrimaryModel
	err = tx.
		// Debug().
		Preload("TagSecondaryList").
		Find(&stored).
		Error
	if err != nil {
		return
	}
	tagPrimaryList = make(domain.TagPrimaryList, len(stored))
	for index, tag := range stored {
		tagPrimaryList[index] = tag.domain()
	}

	return
}

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *Sqlite) GetAllProductList(
	tx *gorm.DB,
) (
	allProductList domain.AllProductList,
	err error,
) {
	var stored []productModel
	err = tx.
		// Table("product").
		Find(&stored).
		Error
	if err != nil {
		return
	}
	allProductList = make(domain.AllProductList, len(stored))
	for index, product := range stored {
		allProductList[index] = product.domain()
	}

	return
}
