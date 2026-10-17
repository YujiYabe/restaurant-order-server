package public

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
		if !env.Env.EnvUnique.IsPermitAlcohol {
			if tagList[index].ID == env.AlcoholIndex { // アルコールは注文しないようにする
				continue
			}
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

	inventoryMapJson, err := transport.MarshalJSON(
		toController.GetInventoryMap(ctx),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	envUnique := env.Env.EnvUnique
	data := struct {
		// PlaceName       string
		StoreID         string
		Environment     string
		URL             string
		StoreName       string
		StoreAddress    string
		Commission      float64
		ChargeBasePrice int
		IsPermitAlcohol bool
		Version         string
		MiniProductMap  string
		TagList         string
		InventoryMap    string
	}{
		// PlaceName:       placeName,
		StoreID:         envUnique.StoreID,
		Environment:     envUnique.Environment,
		URL:             "https://" + envUnique.URL,
		StoreName:       envUnique.StoreName,
		StoreAddress:    envUnique.StoreAddress,
		Commission:      envUnique.Commission,
		ChargeBasePrice: envUnique.ChargeBasePrice,
		IsPermitAlcohol: envUnique.IsPermitAlcohol,
		Version:         envUnique.Version,
		MiniProductMap:  string(miniProductMapJson),
		TagList:         string(tagListJson),
		InventoryMap:    string(inventoryMapJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"public",
		data,
	)
}
