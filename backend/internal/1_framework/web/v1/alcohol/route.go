package alcohol

import (
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"

	includeList "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/alcohol/include_list"
	registerFace "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/alcohol/register_face"
	verifyFace "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/alcohol/verify_face"
	verifyTwenty "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/alcohol/verify_twenty"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/alcohol")

	includeList.NewRoute(
		echoEcho,
		toController,
		group,
	)

	registerFace.NewRoute(
		echoEcho,
		toController,
		group,
	)

	verifyFace.NewRoute(
		echoEcho,
		toController,
		group,
	)

	verifyTwenty.NewRoute(
		echoEcho,
		toController,
		group,
	)

}
