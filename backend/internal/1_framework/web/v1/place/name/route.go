package place_name

import (
	"github.com/labstack/echo"

	// webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/place/name/order"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/:place_name")

	group.PATCH(
		"",
		func(echoContext echo.Context) (
			err error,
		) {
			return patch(echoContext, toController)
		},
		// webUtil.DeviceKeyAuth(),
	)

	group.GET(
		"",
		func(echoContext echo.Context) (
			err error,
		) {
			return get(echoContext, toController)
		},
		// webUtil.DeviceKeyAuth(),
	)

	order.NewRoute(
		echoEcho,
		toController,
		group,
	)

}
