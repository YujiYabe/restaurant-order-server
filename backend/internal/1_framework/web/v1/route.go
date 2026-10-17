package v1

import (
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/admin"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/alcohol"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/internet"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/place"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/smaregi"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/sse"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/v1")

	admin.NewRoute(
		echoEcho,
		toController,
		group,
	)

	alcohol.NewRoute(
		echoEcho,
		toController,
		group,
	)

	internet.NewRoute(
		echoEcho,
		toController,
		group,
	)

	monitor.NewRoute(
		echoEcho,
		toController,
		group,
	)

	place.NewRoute(
		echoEcho,
		toController,
		group,
	)

	smaregi.NewRoute(
		echoEcho,
		toController,
		group,
	)

	sse.NewRoute(
		echoEcho,
		toController,
		group,
	)

}
