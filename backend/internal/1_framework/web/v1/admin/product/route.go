package admin_product

import (
	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"

	// webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/product")

	group.GET( // 管理画面 html
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return get(echoContext, toController)
		},
	)

	group.PATCH( // 管理画面 商品更新
		"/:number",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return patch(echoContext, toController)
		},
	)
}
