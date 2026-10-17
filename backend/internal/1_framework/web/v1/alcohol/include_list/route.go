package include_list

import (
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"

	"github.com/labstack/echo"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/include_list")

	group.GET(
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return get(echoContext, toController)
		},
	)
}
