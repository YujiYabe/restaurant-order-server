package number

import (
	"github.com/labstack/echo"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/:number")

	group.PATCH(
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return patch(echoContext, toController)
		},
		webUtil.DeviceKeyAuth(),
	)
}
