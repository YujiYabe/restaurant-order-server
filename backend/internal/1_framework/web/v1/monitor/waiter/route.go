package waiter

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/waiter/place"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/waiter/reset"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/waiter")

	place.NewRoute(
		echoEcho,
		toController,
		group,
	)

	reset.NewRoute(
		echoEcho,
		toController,
		group,
	)

	group.GET(
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return get(echoContext, toController)
		},
		middleware.BasicAuth(webUtil.CustomBasicAuth),
		webUtil.DeviceKeyAuth(),
	)
}
