package include_list

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// get は、GETリクエストを検証し、対象データまたは画面を応答します。
func get(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)

	alcohol := domain.Alcohol{
		UUID: echoContext.Request().URL.Query().Get("uuidVerifyTwenty"),
	}
	pkg.Logging(ctx, alcohol.UUID)
	ok := toController.IsIncludeTwentyList(
		ctx,
		alcohol,
	)
	if !ok {
		return echoContext.JSON(
			http.StatusNotFound,
			nil,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)

}
