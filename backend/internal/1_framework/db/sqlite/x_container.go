package sqlite

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

type (
	// Sqlite は、SQLiteへの接続を保持します。
	Sqlite struct {
		Conn *gorm.DB
	}
)

// NewToSQLite は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToSQLite() (
	toSqlite gateway.ToSqlite,
) {

	sqliteFilePath := getSqlitePath() + "/" + env.Env.StoreID + ".sqlite3"

	conn, err := open(
		30,
		sqliteFilePath,
	)
	if err != nil {
		panic(err)
	}

	sqlite := new(Sqlite)
	sqlite.Conn = conn
	toSqlite = sqlite

	return
}

// init は、パッケージの初期状態を構成します。
func init() {
	// DBのアップデート

	isUpdateWithInit := false
	if isUpdateWithInit {
		err := updateStoreDB()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// updateStoreDB は、マスターDBの内容を店舗DBへ反映します。
func updateStoreDB() (
	err error,
) {

	// マスターDB =======================
	masterDB, err := open(
		30,
		getSqlitePath()+"/master.sqlite3",
	)
	if err != nil {
		panic(err)
	}

	// ローカルDB =======================
	storeDB, err := open(
		30,
		getSqlitePath()+"/"+env.Env.StoreID+".sqlite3",
	)
	if err != nil {
		panic(err)
	}

	err = updateProduct(
		masterDB,
		storeDB,
	)
	if err != nil {
		panic(err)
	}

	updateTag(
		masterDB,
		storeDB,
	)

	return
}

// getSqlitePath は、SQLiteデータベースを格納するディレクトリーパスを返します。
func getSqlitePath() (
	sqlitePath string,
) {

	currentPath, _ := os.Getwd()

	// localの場合は、db/engine/sqlite を確認
	sqlitePath = currentPath + "/db/engine/sqlite"

	// if !exists(sqlitePath + "/master.sqlite3") {
	// 	// dockerの場合は、/go/src/backend/sqlite を確認
	// 	sqlitePath = filepath.Dir(currentPath) + "/backend/sqlite"
	// }

	return
}

// open は、指定されたSQLiteデータベースを開き、接続設定を適用します。
func open(
	count uint,
	dbPath string,
) (
	db *gorm.DB,
	err error,
) {
	db, err = gorm.Open(
		sqlite.Open(dbPath),
		&gorm.Config{},
	)
	if err != nil {
		if count == 0 {
			err = fmt.Errorf("retry count over")
			return
		}
		time.Sleep(time.Second)
		// カウントダウンさせるようにする
		count--
		return open(count, dbPath) // 再試行
	}

	return
}

// updateProduct は、マスターDBの商品を店舗DBへ追加し、既存商品のマスター項目を更新します。
func updateProduct(
	masterDB, storeDB *gorm.DB,
) (
	err error,
) {
	storeJANCodeList := []int{}

	storeDB.
		// Debug().
		Table("product").
		Select("jan_code").
		Order("jan_code desc").
		Find(&storeJANCodeList)

	restJANCodeList := []int{}

	// storeDBにないjanCodeの検索
	masterDB.
		// Debug().
		Table("product").
		Select("jan_code").
		Order("jan_code desc").
		Not(map[string]interface{}{"jan_code": storeJANCodeList}).
		Find(&restJANCodeList)

	// 差分がなければ終了
	if len(restJANCodeList) == 0 {
		return
	}

	// storeDBにないデータ取得
	differentProductList := []productModel{}
	masterDB.
		// Debug().
		Table("product").
		Where("jan_code IN (?)", restJANCodeList).
		Order("jan_code desc").
		Find(&differentProductList)

	// DB更新
	storeDB.
		// Debug().
		Table("product").
		Create(&differentProductList)

	return
}

// updateTag は、マスターDBのタグ情報を店舗DBへ同期します。
func updateTag(
	masterDB, storeDB *gorm.DB,
) {
	primaryTagList := []tagPrimaryModel{}

	masterDB.
		// Debug().
		// Table("tag_primary").
		Find(&primaryTagList)

	// DB更新
	storeDB.
		// Debug().
		// Table("tag_primary").
		Save(&primaryTagList)

	tagSecondaryList := []tagSecondaryModel{}
	masterDB.
		// Debug().
		// Table("tag_secondary").
		Find(&tagSecondaryList)

	// DB更新
	storeDB.
		// Debug().
		// Table("tag_secondary").
		Save(&tagSecondaryList)
}
