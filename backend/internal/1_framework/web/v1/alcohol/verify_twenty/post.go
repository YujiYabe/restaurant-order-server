package verify_twenty

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

type Response struct {
	Age      string `json:"age"`
	Birthday string `json:"birthday"`
	UUID     string `json:"uuid"`
}

// post は、POSTリクエストを検証し、対応するユースケースを実行して結果を応答します。
func post(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	// 運転免許証かマイナンバーカードを写真撮影して送信する
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	request := transport.Alcohol{}
	if err := echoContext.Bind(&request); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	alcohol := request.Domain()

	id, err := uuid.NewUUID()
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	alcohol.UUID = id.String()
	pkg.Logging(ctx, alcohol.UUID)

	_, age, birthday, err := toController.VerifyTwenty(
		ctx,
		alcohol,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	response := Response{
		Age:      strconv.Itoa(age),
		Birthday: birthday.Format("2006/01/02"),
		UUID:     alcohol.UUID,
	}

	return echoContext.JSON(
		http.StatusOK,
		response,
	)
}
