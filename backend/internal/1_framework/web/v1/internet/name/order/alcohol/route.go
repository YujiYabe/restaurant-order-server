package alcohol

import (
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/alcohol")

	group.POST(
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return post(echoContext, toController)
		},
	)
}
