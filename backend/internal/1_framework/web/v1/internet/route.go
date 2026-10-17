package internet

import (
	"github.com/labstack/echo"

	name "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/internet/name"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/internet")
	name.NewRoute(
		echoEcho,
		toController,
		group,
	)
	// internetは会計ごとにUUIDを生成
}
