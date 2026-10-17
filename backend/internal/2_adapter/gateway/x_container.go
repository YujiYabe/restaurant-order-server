package gateway

import (
	"context"

	"gorm.io/gorm"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

type Gateway struct {
	ToSqlite
	ToRedis
	ToFirestore
	ToSmaregi
	ToGate
	ToAzure
	ToAws
}

type NewGatewayOption struct {
	ToSqlite
	ToRedis
	ToFirestore
	ToSmaregi
	ToGate
	ToAzure
	ToAws
}

// NewGateway は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewGateway(
	newGatewayOption NewGatewayOption,
) (
	gateway *Gateway,
) {
	gateway = &Gateway{
		ToSqlite:    newGatewayOption.ToSqlite,
		ToRedis:     newGatewayOption.ToRedis,
		ToFirestore: newGatewayOption.ToFirestore,
		ToSmaregi:   newGatewayOption.ToSmaregi,
		ToGate:      newGatewayOption.ToGate,
		ToAzure:     newGatewayOption.ToAzure,
		ToAws:       newGatewayOption.ToAws,
	}

	return
}

type ToSqlite interface {
	WithOutTx() (
		tx *gorm.DB,
	)

	BeginTx() (
		tx *gorm.DB,
	)

	EndTx(
		tx *gorm.DB,
		isSuccess bool,
	) error

	GetProduct(
		tx *gorm.DB,
		janCode int,
	) (
		product domain.Product,
		err error,
	)

	GetAllProductList(
		tx *gorm.DB,
	) (
		allProductList domain.AllProductList,
		err error,
	)

	UpdateProduct(
		tx *gorm.DB,
		product domain.Product,
		newProduct domain.Product,
	) (
		err error,
	)

	UpdateProductNoSaga(
		ctx context.Context,
		product domain.Product,
	) (
		err error,
	)

	GetAllTag(
		tx *gorm.DB,
	) (
		tagPrimaryList domain.TagPrimaryList,
		err error,
	)

	UpdateStock(
		tx *gorm.DB,
		product domain.Product,
		count int,
	) (
		err error,
	)
}

type ToAzure interface {
	VerifyCertificate(
		ctx context.Context,
		alcohol domain.Alcohol,
	) (
		birthday string,
		number string,
		err error,
	)
}

type ToAws interface {
	Identification(
		ctx context.Context,
		imageSet []string,
	) (
		ok bool,
		err error,
	)
}

type ToRedis interface {
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
}

type ToFirestore interface {
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

	RegisterFace(
		ctx context.Context,
		alcohol domain.Alcohol,
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
}

type ToSmaregi interface {
	SendPostTransactionTemporary(
		ctx context.Context,
		reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
	) (
		smaregiPostResponseEnquiry domain.SmaregiPostResponseEnquiry,
		err error,
	)
}

type ToGate interface {
	SendPostTransactionInfo(
		ctx context.Context,
		transactionHeadId string,
	) (
		err error,
	)
}
