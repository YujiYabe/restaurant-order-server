package order

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// post は、POSTリクエストを検証し、対応するユースケースを実行して結果を応答します。
func post(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	// orderの登録
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

	// お酒の注文は別APIなのでこのAPIでは受け付けない
	isAlcohol := toController.IsAlcohol(
		ctx,
		order,
	)
	if isAlcohol {
		return echoContext.JSON(
			http.StatusBadRequest,
			nil,
		)
	}

	// 注文
	placeName := echoContext.Param("place_name")
	ok := toController.AddOrder(
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
