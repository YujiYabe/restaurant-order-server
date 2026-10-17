package verify_twenty

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
	group := parent.Group("/verify_twenty")
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
