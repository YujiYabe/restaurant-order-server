package number

import (
	"net/http"
	"strconv"

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
	// orderの編集
	request := transport.Order{}
	if err := echoContext.Bind(&request); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	order := request.Domain()
	pkg.Logging(ctx, order)

	placeName := echoContext.Param("place_name")
	orderNumber, err := strconv.Atoi(echoContext.Param("number"))
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	order.OrderNumber = orderNumber

	ok := toController.EditOrder(
		ctx,
		placeName,
		order,
	)
	if !ok {
		return echoContext.JSON(
			http.StatusBadRequest,
			nil,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)

}
