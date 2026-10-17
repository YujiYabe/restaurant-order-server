package domain

import (
	"context"
	"sync"
	"time"
)

// Domain は、席、注文、商品、在庫、成人確認のインメモリ状態を保持します。
// 可変フィールドはmuで保護され、チャージ自動更新ループの停止関数も保持します。
type Domain struct {
	mutex             sync.RWMutex
	autoUpdateCancel  context.CancelFunc
	PlaceList         PlaceList
	Stock             Stock
	AllergyDefault    Allergy
	TagPrimaryList    TagPrimaryList
	InventoryMap      InventoryMap
	MiniProductMap    MiniProductMap
	AlcoholList       AlcoholList
	AllowedTwentyList TwentyList
	PendingFaceList   TwentyList
}

// NewDomainOption は、Domain生成時の拡張オプションを表します。
// 現在は設定項目を持ちません。
type NewDomainOption struct{}

// NewDomain は、各一覧とマップを空の初期値で初期化したDomainを生成します。
// 成人確認済み一覧と顔登録待ち一覧は独立したマップとして初期化されます。
func NewDomain(
	newDomainOption NewDomainOption,
) (
	toDomain ToDomain,
) {
	toDomain = &Domain{
		PlaceList:         NewPlaceList(),
		Stock:             NewStock(),
		AllergyDefault:    NewAllergyDefault(),
		TagPrimaryList:    NewTagPrimaryList(),
		InventoryMap:      NewInventoryMap(),
		MiniProductMap:    NewMiniProductMap(),
		AlcoholList:       NewAlcoholList(),
		AllowedTwentyList: NewTwentyList(),
		PendingFaceList:   NewTwentyList(),
	}

	return
}

// ToDomain は、上位層が使用するドメイン操作を定義します。
// 参照系操作は原則としてスナップショットを返し、複合更新はDomain内のロック下で実行されます。
type ToDomain interface {
	GetDomain(
		ctx context.Context,
	) (
		domain *Domain,
	)

	// 商品 -----------------------
	GetProduct(
		ctx context.Context,
		productNumber int,
	) (
		product Product,
	)

	GetAllProductList(
		ctx context.Context,
	) (
		allProductList AllProductList,
	)

	GetMiniProductMap(
		ctx context.Context,
	) (
		miniProductMap MiniProductMap,
	)

	UpdateProduct(
		ctx context.Context,
		product Product,
	)

	SaveInMemory(
		ctx context.Context,
		allProductList AllProductList,
		tagPrimaryList TagPrimaryList,
		placeList PlaceList,
	)

	SetPlaceList(
		ctx context.Context,
		placeList PlaceList,
	)

	VerifyJANCodeList(
		janCodeList []int,
		isValidLangCodeList []int,
		defaultLangCode int,
	) (
		newIsValidJANCodeList []int,
		newIsValidLangCode int,
	)

	AdjustmentInventory(
		ctx context.Context,
		newOrder Order,
	) (
		ok bool,
	)

	TryAddOrder(
		ctx context.Context,
		placeName string,
		newOrder Order,
	) (
		ok bool,
	)

	GetIsValidJANCodeList() (
		isValidJANCodeList []int,
	)

	// 注文 -----------------------
	GetPlaceList(
		ctx context.Context,
	) (
		placeList PlaceList,
	)

	GetPlace(
		ctx context.Context,
		placeName string,
	) (
		place *Place,
	)

	EditPlace(
		ctx context.Context,
		place Place,
	) (ok bool)

	AllowFaceRegistration(
		ctx context.Context,
		alcohol Alcohol,
	)

	ConsumeFaceRegistration(
		ctx context.Context,
		alcohol Alcohol,
	) (
		ok bool,
	)

	UpdateBillingProgress(
		ctx context.Context,
		placeName string,
		transactionHeadID string,
		stage BillingStage,
	) (
		ok bool,
	)

	AddOrder(
		ctx context.Context,
		placeName string,
		newOrder Order,
		miniProductMap MiniProductMap,
	)

	GetPermission(
		ctx context.Context,
		placeName string,
	) (
		permission string,
	)

	EditOrder(
		ctx context.Context,
		placeName string,
		targetOrder Order,
		miniProductMap MiniProductMap,
	) (
		ok bool,
	)

	BuildPostTransactionTemporary(
		ctx context.Context,
		placeName string,
		miniProductMap MiniProductMap,
	) (
		reqPostTransactionTemporary ReqPostTransactionTemporary,
		janCodeCountMap map[int]int,
	)

	// アレルギー -----------------------
	GetAllergyDefault(
		ctx context.Context,
	) (
		allergy Allergy,
	)

	// タグ -----------------------
	GetTagList(
		ctx context.Context,
	) (
		tagPrimaryList TagPrimaryList,
	)

	// 在庫 -----------------------
	GetInventoryMap(
		ctx context.Context,
	) (
		inventoryMap InventoryMap,
	)

	// 酒類 -----------------------
	VerifyTwenty(
		ctx context.Context,
		birthdayString string,
	) (
		age int,
		birthday time.Time,
		err error,
	)

	IsAlcohol(
		ctx context.Context,
		order Order,
	) (
		ok bool,
	)

	IsIncludeTwentyList(
		ctx context.Context,
		alcohol Alcohol,
	) (
		ok bool,
	)

	AddTwentyList(
		ctx context.Context,
		alcohol Alcohol,
	)
}
