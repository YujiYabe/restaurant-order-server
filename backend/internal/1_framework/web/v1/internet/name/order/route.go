package order

import (
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/internet/name/order/alcohol"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/order")
	group.POST( // 認証は都度生成されるUUIDで行う
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return post(echoContext, toController)
		},
	)

	alcohol.NewRoute(
		echoEcho,
		toController,
		group,
	)

}
