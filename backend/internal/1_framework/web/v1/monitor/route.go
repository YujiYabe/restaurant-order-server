package monitor

import (
	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"

	// webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	backYard "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/back_yard"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/internet"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/public"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/reception"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/repository"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/signboard"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/monitor/waiter"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group(
		"/monitor",
		// middleware.BasicAuth(webUtil.CustomBasicAuth),
	)

	backYard.NewRoute(
		echoEcho,
		toController,
		group,
	)

	internet.NewRoute(
		echoEcho,
		toController,
		group,
	)

	public.NewRoute(
		echoEcho,
		toController,
		group,
	)

	reception.NewRoute(
		echoEcho,
		toController,
		group,
	)

	repository.NewRoute(
		echoEcho,
		toController,
		group,
	)

	signboard.NewRoute(
		echoEcho,
		toController,
		group,
	)

	waiter.NewRoute(
		echoEcho,
		toController,
		group,
	)
}
