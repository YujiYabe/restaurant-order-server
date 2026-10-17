package sse

import (
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/sse/app"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	go app.SendToAgentPlaceList()

	group := parent.Group("/sse")

	app.NewRoute(
		echoEcho,
		toController,
		group,
	)
}
