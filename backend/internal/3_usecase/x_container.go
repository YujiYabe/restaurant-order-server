package usecase

import (
	"context"
	"sync"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

// UseCase は、アプリケーションの処理を調整します。
type UseCase struct {
	orderMu     sync.Mutex
	ToDomain    domain.ToDomain
	ToGateway   ToGateway
	ToPresenter ToPresenter
}

type NewUseCaseOption struct {
	ToDomain    domain.ToDomain
	ToGateway   ToGateway
	ToPresenter ToPresenter
}

// NewUseCase は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewUseCase(
	newUseCaseOption NewUseCaseOption,
) (
	toUseCase ToUseCase,
) {
	toUseCase = &UseCase{
		ToDomain:    newUseCaseOption.ToDomain,
		ToGateway:   newUseCaseOption.ToGateway,
		ToPresenter: newUseCaseOption.ToPresenter,
	}

	return
}

// ToUseCase は、ユースケースへ公開する処理を定義します。
type ToUseCase interface {
	Start(
		ctx context.Context,
	)

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

	// 席 -----------------------
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

	// SSE配信 -----------------------
	DistributeAgent(
		ctx context.Context,
	)
}

// ToGateway は、ゲートウェイへ公開する処理を定義します。
type ToGateway interface {
	GetAllProductList(
		ctx context.Context,
	) (
		allProductList domain.AllProductList,
		err error,
	)

	UpdateProduct(
		ctx context.Context,
		product domain.Product,
	) (
		err error,
	)

	GetAllTag(
		ctx context.Context,
	) (
		tagPrimaryList domain.TagPrimaryList,
		err error,
	)

	UpdateStock(
		ctx context.Context,
		janCodeCountMap map[int]int,
	) (
		err error,
	)

	SendPostTransactionTemporary(
		ctx context.Context,
		reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
	) (
		smaregiPostResponseEnquiry domain.SmaregiPostResponseEnquiry,
		err error,
	)

	SendPostTransactionInfo(
		ctx context.Context,
		transactionHeadId string,
	) (
		err error,
	)

	SaveOrder(
		ctx context.Context,
		place domain.Place,
		reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
	) (
		err error,
	)

	GetClientImage(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		imageList []string,
		err error,
	)

	SaveClientAttribute(
		ctx context.Context,
		alcohol domain.Alcohol,
		certificateID string,
	) (
		err error,
	)

	SaveClientOrder(
		ctx context.Context,
		clientUUID string,
		imageUUID string,
		janCode int,
		count int,
	) (
		err error,
	)

	GetPlaceListFromRedis(
		ctx context.Context,
	) (
		placeList domain.PlaceList,
		err error,
	)

	SetPlaceListToRedis(
		ctx context.Context,
		placeList domain.PlaceList,
	) (
		err error,
	)

	ResetPlaceListInRedis(
		ctx context.Context,
	) (
		err error,
	)

	// 酒類 -----------------------
	VerifyCertificate(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		birthday string,
		number string,
		err error,
	)

	RegisterFace(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		err error,
	)

	Identification(
		ctx context.Context,
		imageSet []string,
	) (
		ok bool,
		err error,
	)
}

// ToPresenter は、プレゼンターへ公開する処理を定義します。
type ToPresenter interface {
	DistributeAgent(
		ctx context.Context,
		domainDomain *domain.Domain,
	)
}
