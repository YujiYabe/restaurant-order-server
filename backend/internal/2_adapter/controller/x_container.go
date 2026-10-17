package controller

import (
	"context"
	"time"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/presenter"
	usecase "github.com/YujiYabe/restaurant-order-server/backend/internal/3_usecase"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// controller は、Web入力とユースケースを仲介します。
type controller struct {
	ToUseCase usecase.ToUseCase
}

type NewControllerOption struct {
	ToSqlite    gateway.ToSqlite
	ToRedis     gateway.ToRedis
	ToFirestore gateway.ToFirestore
	ToSmaregi   gateway.ToSmaregi
	ToGate      gateway.ToGate
	ToAzure     gateway.ToAzure
	ToAws       gateway.ToAws
	ToMonitor   presenter.ToMonitor
}

// NewController は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewController(
	newControllerOption NewControllerOption,
) (
	toController ToController,
) {

	// ドメイン ---------------------------
	newDomainOption := domain.NewDomainOption{}
	toDomain := domain.NewDomain(
		newDomainOption,
	)

	// ゲートウェイ ---------------------------
	newGatewayOption := gateway.NewGatewayOption{
		ToSqlite:    newControllerOption.ToSqlite,
		ToRedis:     newControllerOption.ToRedis,
		ToFirestore: newControllerOption.ToFirestore,
		ToSmaregi:   newControllerOption.ToSmaregi,
		ToGate:      newControllerOption.ToGate,
		ToAzure:     newControllerOption.ToAzure,
		ToAws:       newControllerOption.ToAws,
	}
	toGateway := gateway.NewGateway(
		newGatewayOption,
	)

	// プレゼンター ---------------------------
	newPresenter := presenter.NewPresenterOption{
		ToMonitor: newControllerOption.ToMonitor,
	}
	toPresenter := presenter.NewPresenter(
		newPresenter,
	)

	// ユースケース ---------------------------
	newUseCaseOption := usecase.NewUseCaseOption{
		ToDomain:    toDomain,
		ToGateway:   toGateway,
		ToPresenter: toPresenter,
	}

	toUseCase := usecase.NewUseCase(
		newUseCaseOption,
	)

	toController = &controller{
		ToUseCase: toUseCase,
	}

	return
}

// ToController は、コントローラーへ公開する処理を定義します。
type ToController interface {
	Start()

	// 商品 -----------------------
	GetProduct(
		ctx context.Context,
		productNumber int,
	) (
		product domain.Product,
	)
	GetAllProductList(
		ctx context.Context,
	) (
		allProductList domain.AllProductList,
	)

	GetMiniProductMap(
		ctx context.Context,
	) (
		miniProductMap domain.MiniProductMap,
	)

	UpdateProduct(
		ctx context.Context,
		product domain.Product,
	) (
		err error,
	)

	// 注文 -----------------------
	ResetPlaceList(
		ctx context.Context,
	) (
		err error,
	)

	GetPlaceList(
		ctx context.Context,
	) (
		placeList domain.PlaceList,
	)

	GetPlace(
		ctx context.Context,
		placeName string,
	) (
		place domain.Place,
	)

	EditPlace(
		ctx context.Context,
		place domain.Place,
	) (
		err error,
	)

	AddOrder(
		ctx context.Context,
		placeName string,
		order domain.Order,
	) (
		ok bool,
	)

	IsAlcohol(
		ctx context.Context,
		order domain.Order,
	) (
		ok bool,
	)

	IsIncludeTwentyList(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		ok bool,
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
		order domain.Order,
	) (
		ok bool,
	)

	// アレルギー -----------------------
	GetAllergyDefault(
		ctx context.Context,
	) (
		allergy domain.Allergy,
	)

	// タグ -----------------------
	GetTagList(
		ctx context.Context,
	) (
		tagPrimaryList domain.TagPrimaryList,
	)

	// 在庫 -----------------------
	GetInventoryMap(
		ctx context.Context,
	) (
		inventoryMap domain.InventoryMap,
	)

	// 酒類 -----------------------
	VerifyTwenty(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		certificateID string,
		age int,
		birthday time.Time,
		err error,
	)

	RegisterFace(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		err error,
	)

	VerifyFace(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		err error,
	)

	// SSE配信 -----------------------
	DistributeAgent(
		ctx context.Context,
	)
}
