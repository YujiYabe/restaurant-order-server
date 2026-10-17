package admin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	place "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/admin/place"
	product "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/admin/product"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {

	group := parent.Group(
		"/admin",
		middleware.BasicAuth(webUtil.CustomBasicAuth),
		webUtil.DeviceKeyAuth(),
	)

	product.NewRoute(
		echoEcho,
		toController,
		group,
	)

	place.NewRoute(
		echoEcho,
		toController,
		group,
	)

}
