package place

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// get は、GETリクエストを検証し、対象データまたは画面を応答します。
func get(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	placeName := echoContext.Param("place_name")

	miniProductMapJson, err := transport.MarshalJSON(
		toController.GetMiniProductMap(ctx),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	tagList := toController.GetTagList(ctx)
	filteredTagList := domain.TagPrimaryList{}
	for index := range tagList {
		if tagList[index].ID == env.AlcoholIndex { // アルコールは注文しないようにする
			continue
		}
		filteredTagList = append(
			filteredTagList,
			tagList[index],
		)
	}

	tagListJson, err := transport.MarshalJSON(
		filteredTagList,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	permission := toController.GetPermission(
		ctx,
		placeName,
	)

	envUnique := env.Env.EnvUnique
	data := struct {
		PlaceName  string
		Permission string

		StoreID          string
		Environment      string
		URL              string
		StoreName        string
		StoreAddress     string
		Commission       float64
		ChargeBaseMinute int
		ChargeBasePrice  int
		IsPermitAlcohol  bool
		IsEnableCharge   bool
		Version          string

		MiniProductMap string
		TagList        string
	}{
		PlaceName:  placeName,
		Permission: permission,

		StoreID:          envUnique.StoreID,
		Environment:      envUnique.Environment,
		URL:              "https://" + envUnique.URL,
		StoreName:        envUnique.StoreName,
		StoreAddress:     envUnique.StoreAddress,
		Commission:       envUnique.Commission,
		ChargeBaseMinute: envUnique.ChargeBaseMinute,
		ChargeBasePrice:  envUnique.ChargeBasePrice,
		IsPermitAlcohol:  envUnique.IsPermitAlcohol,
		IsEnableCharge:   envUnique.IsEnableCharge,
		Version:          envUnique.Version,

		MiniProductMap: string(miniProductMapJson),
		TagList:        string(tagListJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"waiter_place",
		data,
	)
}
