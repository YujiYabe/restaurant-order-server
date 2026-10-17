package smaregi

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// "github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// NewRoute は、このパッケージが提供するHTTPルートとハンドラーを登録します。
func NewRoute(
	echoEcho *echo.Echo,
	toController controller.ToController,
	parent *echo.Group,
) {
	group := parent.Group("/smaregi")

	group.POST(
		"",
		func(
			echoContext echo.Context,
		) (
			err error,
		) {
			return post(echoContext, toController)
		},
		webUtil.SmaregiKeyAuth(),
		middleware.BodyDump(bodyDumpHandler),
	)
}

// bodyDumpHandler は、受信したWebhookリクエストボディーをログへ記録します。
func bodyDumpHandler(
	echoContext echo.Context,
	reqBody,
	resBody []byte,
) {
	log.Printf("Req Body: %v\n", string(reqBody))
	log.Printf("Res Body: %v\n", string(resBody))
}
