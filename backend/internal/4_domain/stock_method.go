package domain

import (
	"context"
	"sort"

	"golang.org/x/exp/slices"
)

// GetProduct は、指定JANコードと一致する商品を全商品一覧から返します。
// 一致する商品がない場合はProductのゼロ値を返します。
func (receiver *Stock) GetProduct(
	ctx context.Context,
	janCode int,
) (
	newProduct Product,
) {
	newProduct = Product{}
	for _, product := range receiver.AllProductList {
		if janCode == product.JANCode {
			newProduct = product
		}
	}
	return
}

// GetAllProductList は、有効・無効を含む全商品一覧のコピーを返します。
func (receiver *Stock) GetAllProductList(
	ctx context.Context,
) (
	allProductList AllProductList,
) {
	allProductList = receiver.AllProductList
	return
}

// GetProductList は、有効な商品一覧を返します。
func (receiver *Stock) GetProductList(
	ctx context.Context,
) (
	productList ProductList,
) {
	productList = receiver.ProductList
	return
}

// GetIsValidJANCodeList は、注文可能なJANコード一覧を返します。
// Domain経由の公開操作では、この結果をさらにコピーして返します。
func (receiver *Stock) GetIsValidJANCodeList() (
	isValidJANCodeList []int,
) {
	isValidJANCodeList = receiver.IsValidJANCodeList
	return
}

// UpdateProduct は、AllProductList内で同じJANコードを持つ商品の運用上更新可能な項目を置き換えます。
// 対象項目は配置、有効性、在庫、および各価格です。一致する商品がない場合は一覧を変更しません。
func (receiver *Stock) UpdateProduct(
	ctx context.Context,
	newProduct Product,
) {
	for index := range receiver.AllProductList {
		if newProduct.JANCode != receiver.AllProductList[index].JANCode {
			continue
		}

		// インメモリ更新
		receiver.AllProductList[index].Place = newProduct.Place
		receiver.AllProductList[index].IsValid = newProduct.IsValid
		receiver.AllProductList[index].Stock = newProduct.Stock
		receiver.AllProductList[index].CostPrice = newProduct.CostPrice
		receiver.AllProductList[index].RetailPrice = newProduct.RetailPrice
		receiver.AllProductList[index].OperationCost = newProduct.OperationCost
		receiver.AllProductList[index].EnergyCost = newProduct.EnergyCost
		break
	}
}

// VerifyJANCodeList は、入力整数一覧を昇順にソートし、販売可能なJANコードと言語コードを分離します。
// 有効な言語コードが1つだけ含まれる場合はその値を、それ以外はdefaultLangCodeを返します。
// janCodeListはソートにより呼び出し側のスライスも更新されます。
func (receiver *Stock) VerifyJANCodeList(
	janCodeList []int,
	isValidLangCodeList []int,
	defaultLangCode int,
) (
	newIsValidJANCodeList []int,
	newIsValidLangCode int,
) {
	newIsValidLangCodeList := []int{}

	sort.Ints(janCodeList)

	for _, janCode := range janCodeList {
		if slices.Contains(isValidLangCodeList, janCode) {
			newIsValidLangCodeList = append(
				newIsValidLangCodeList,
				janCode,
			)
		}

		if slices.Contains(receiver.IsValidJANCodeList, janCode) {
			newIsValidJANCodeList = append(
				newIsValidJANCodeList,
				janCode,
			)
		}
	}

	newIsValidLangCode = defaultLangCode // 既定値は日本語の1041

	if len(newIsValidLangCodeList) == 1 { // 言語が１つ指定されていればその言語にする。それ以外の場合変更しない。
		newIsValidLangCode = newIsValidLangCodeList[0]
	}

	return
}
