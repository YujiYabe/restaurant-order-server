package back_yard

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group(
		"/back_yard",
		middleware.BasicAuth(webUtil.CustomBasicAuth),
		webUtil.DeviceKeyAuth(),
	)

	group.GET(
		"/:place_name",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return get(echoContext, toController)
		},
	)
}
