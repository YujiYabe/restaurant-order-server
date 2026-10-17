package signboard

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
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

	signboardProductJson, err := transport.MarshalJSON(env.Env.SignboardProduct)
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
		StoreID          string
		Environment      string
		URL              string
		StoreName        string
		StoreAddress     string
		Commission       float64
		Version          string
		MiniProductMap   string
		SignboardProduct string
	}{
		// PlaceName:       placeName,
		StoreID:          envUnique.StoreID,
		Environment:      envUnique.Environment,
		URL:              "https://" + envUnique.URL,
		StoreName:        envUnique.StoreName,
		StoreAddress:     envUnique.StoreAddress,
		Commission:       envUnique.Commission,
		Version:          envUnique.Version,
		MiniProductMap:   string(miniProductMapJson),
		SignboardProduct: string(signboardProductJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"signboard",
		data,
	)
}
