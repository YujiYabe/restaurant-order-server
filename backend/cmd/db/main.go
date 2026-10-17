package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// システム環境変数 STORE_ID の読み込み
	storeID := os.Getenv("STORE_ID")
	err := migrateTable(storeID)
	if err != nil {
		log.Fatalf("DBの取得に失敗しました: %v", err)
	}

	nowProductList, err := getNowProductList(storeID)
	if err != nil {
		log.Fatalf("DBの取得に失敗しました: %v", err)
	}

	newProductList, err := getNewProductList()
	if err != nil {
		log.Fatalf("商品tsvの取得に失敗しました: %v", err)
	}

	productList := generateProductList(nowProductList, newProductList)

	tagPrimaryList, err := getTagPrimary()
	if err != nil {
		log.Fatalf("カテゴリtsvの取得に失敗しました: %v", err)
	}

	tagSecondaryList, err := getTagSecondary()
	if err != nil {
		log.Fatalf("サブカテゴリtsvの取得に失敗しました: %v", err)
	}

	err = replaceTable(
		storeID,
		productList,
		tagPrimaryList,
		tagSecondaryList,
	)
	if err != nil {
		log.Fatalf("トランザクションに失敗しました: %v", err)
	}

	return
}

func replaceTable(
	storeID string,
	productList []Product,
	tagPrimaryList []TagPrimary,
	tagSecondaryList []TagSecondary,
) (
	err error,
) {
	db, err := gorm.Open(
		sqlite.Open("../../db/engine/sqlite/"+storeID+".sqlite3"),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := executeTransaction(tx, "product", productList); err != nil {
		return err
	}
	if err := executeTransaction(tx, "tag_primary", tagPrimaryList); err != nil {
		return err
	}
	if err := executeTransaction(tx, "tag_secondary", tagSecondaryList); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("コミット失敗: %w", err)
	}

	return nil
}

func executeTransaction(
	tx *gorm.DB,
	tableName string,
	data interface{},
) (
	err error,
) {
	if err := tx.Exec("DELETE FROM " + tableName).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("削除失敗: %w", err)
	}
	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("挿入失敗: %w", err)
	}
	return nil
}

func getNewProductList() (
	newProductList []Product,
	err error,
) {
	file, err := os.Open("../../db/engine/sqlite/product.tsv")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) < int(UpdatedAt)+1 {
			return nil, fmt.Errorf("フィールド数が不足しています：%v", line)
		}

		product := Product{}

		product.JANCode, err = strconv.Atoi(fields[JANCode])
		if err != nil {
			return
		}

		// -----------------------------
		product.ProductId, err = strconv.Atoi(fields[ProductId])
		if err != nil {
			return
		}

		// -----------------------------
		product.IsValid, err = strconv.ParseBool(fields[IsValid])
		if err != nil {
			return
		}

		// -----------------------------
		product.Stock, err = strconv.Atoi(fields[Stock])
		if err != nil {
			return
		}

		// -----------------------------
		product.Place = fields[Place]

		// -----------------------------
		product.Calorie, err = strconv.Atoi(fields[Calorie])
		if err != nil {
			return
		}

		// -----------------------------
		product.CostPrice, err = strconv.Atoi(fields[CostPrice])
		if err != nil {
			return
		}

		// -----------------------------
		product.RetailPrice, err = strconv.Atoi(fields[RetailPrice])
		if err != nil {
			return
		}

		// -----------------------------
		product.OperationCost, err = strconv.Atoi(fields[OperationCost])
		if err != nil {
			return
		}

		// -----------------------------
		product.EnergyCost, err = strconv.Atoi(fields[EnergyCost])
		if err != nil {
			return
		}

		// -----------------------------
		product.Recipe = fields[Recipe]

		// -----------------------------
		product.Tag = fields[Tag]

		// -----------------------------
		product.Allergy = fields[Allergy]

		// -----------------------------
		product.NameJa = fields[NameJa]

		// -----------------------------
		product.NameJaShort = fields[NameJaShort]

		// -----------------------------
		product.InformationJa = fields[InformationJa]

		// -----------------------------
		product.UpdatedAt = fields[UpdatedAt]

		newProductList = append(newProductList, product)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	sortProductList(newProductList)

	return
}
func migrateTable(
	storeID string,
) (
	err error,
) {
	db, err := gorm.Open(
		sqlite.Open("../../db/engine/sqlite/"+storeID+".sqlite3"),
		&gorm.Config{},
	)
	if err != nil {
		return
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&TagPrimary{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&TagSecondary{})
	if err != nil {
		return
	}

	return
}

func getNowProductList(
	storeID string,
) (
	nowProductList []Product,
	err error,
) {
	db, err := gorm.Open(
		sqlite.Open("../../db/engine/sqlite/"+storeID+".sqlite3"),
		&gorm.Config{},
	)
	if err != nil {
		return
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		return
	}

	err = db.
		Find(&nowProductList).
		Error
	if len(nowProductList) == 0 {
		err = nil
		return
	}

	sortProductList(nowProductList)
	return
}

func sortProductList(
	productList []Product,
) {
	sort.Slice(productList, func(i, j int) bool {
		return productList[i].JANCode < productList[j].JANCode
	})
}

func generateProductList(
	nowProductList []Product,
	newProductList []Product,
) (
	productList []Product,
) {

	for _, newProductValue := range newProductList {
		for _, nowProductValue := range nowProductList {
			if newProductValue.JANCode == nowProductValue.JANCode {
				newProductValue.Stock = nowProductValue.Stock
				newProductValue.IsValid = nowProductValue.IsValid
				newProductValue.Place = nowProductValue.Place
				break
			}
		}
		productList = append(productList, newProductValue)
	}

	return
}

func getTagPrimary() (
	tagPrimaryList []TagPrimary,
	err error,
) {
	file, err := os.Open("../../db/engine/sqlite/tag_primary.tsv")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		tagPrimary := TagPrimary{}

		// -----------------------------
		tagPrimary.ID, err = strconv.Atoi(fields[0])
		if err != nil {
			return
		}

		// -----------------------------
		tagPrimary.Name = fields[1]

		tagPrimaryList = append(tagPrimaryList, tagPrimary)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func getTagSecondary() (
	tagSecondaryList []TagSecondary,
	err error,
) {
	file, err := os.Open("../../db/engine/sqlite/tag_secondary.tsv")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		tagSecondary := TagSecondary{}

		// -----------------------------
		tagSecondary.ID, err = strconv.Atoi(fields[0])
		if err != nil {
			return
		}

		// -----------------------------
		tagSecondary.TagPrimaryID, err = strconv.Atoi(fields[1])
		if err != nil {
			return
		}

		// -----------------------------
		tagSecondary.Name = fields[2]

		// -----------------------------
		tagSecondaryList = append(tagSecondaryList, tagSecondary)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
