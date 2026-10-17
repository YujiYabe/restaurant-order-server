package register_face

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// post は、POSTリクエストを検証し、対応するユースケースを実行して結果を応答します。
func post(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	// 本人確認画像
	request := transport.Alcohol{}
	if err := echoContext.Bind(&request); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	alcohol := request.Domain()

	pkg.Logging(ctx, alcohol.UUID)
	err = toController.RegisterFace(
		ctx,
		alcohol,
	)
	if err != nil {
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		alcohol.UUID,
	)
}
