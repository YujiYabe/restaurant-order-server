package admin_place

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// ---------------------------------------------------------
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

	err = toController.ResetPlaceList(ctx)
	if err != nil {
		pkg.Logging(ctx, err)
		// return echoContext.JSON(
		// 	http.StatusInternalServerError,
		// 	nil,
		// )
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)
}
