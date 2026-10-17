package place

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	name "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/place/name"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {

	group := parent.Group(
		"/place",
		middleware.BasicAuth(webUtil.CustomBasicAuth),
		webUtil.DeviceKeyAuth(),
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
		webUtil.DeviceKeyAuth(),
	)

	name.NewRoute(
		echoEcho,
		toController,
		group,
	)
}
