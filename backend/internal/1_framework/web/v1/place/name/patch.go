package place_name

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// patch は、PATCHリクエストを検証し、対象データを更新します。
func patch(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)

	request := transport.Place{}
	if err := echoContext.Bind(&request); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	place := request.Domain()

	place.PlaceName = echoContext.Param("place_name")
	pkg.Logging(ctx, place)

	currentPlace := toController.GetPlace(ctx, place.PlaceName)
	if currentPlace.PlaceName == "" {
		return echoContext.JSON(
			http.StatusNotFound,
			nil,
		)
	}

	if err := toController.EditPlace(ctx, place); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadGateway,
			err,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)

}
