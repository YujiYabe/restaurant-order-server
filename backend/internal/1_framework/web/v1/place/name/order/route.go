package order

import (
	"github.com/labstack/echo"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	number "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/place/name/order/number"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/order")
	group.POST(
		"",
		func(echoContext echo.Context) (
			err error,
		) {
			return post(echoContext, toController)
		},
		webUtil.DeviceKeyAuth(),
	)

	number.NewRoute(
		echoEcho,
		toController,
		group,
	)
}
