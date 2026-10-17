package place_name

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
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

	place := toController.GetPlace(
		ctx,
		placeName,
	)
	if place.PlaceName == "" {
		return echoContext.JSON(
			http.StatusNotFound,
			nil,
		)
	}
	pkg.Logging(ctx, place)

	return echoContext.JSON(
		http.StatusOK,
		transport.JSONValue(place),
	)
}
