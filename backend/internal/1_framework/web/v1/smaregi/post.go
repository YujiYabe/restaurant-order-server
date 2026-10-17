package smaregi

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

type GateWebhook struct {
	TransactionHeadId string `json:"transactionHeadId"`
}

// post は、POSTリクエストを検証し、対応するユースケースを実行して結果を応答します。
func post(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	// スマレジのwebhookAPIで支払い終了のあとに
	// テーブルステータスを会計中→掃除中に変更する
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	gateWebhook := GateWebhook{}
	if err := echoContext.Bind(&gateWebhook); err != nil {
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	pkg.Logging(ctx, gateWebhook)
	if gateWebhook.TransactionHeadId == "" {
		return echoContext.JSON(http.StatusBadRequest, "transactionHeadId is required")
	}

	newPlace := domain.Place{}

	placeList := toController.GetPlaceList(ctx)

	isExist := false
	for _, place := range placeList {
		if place == nil {
			continue
		}
		isPlaceStatusBilling := place.PlaceStatus == domain.PlaceStatusBilling
		isTargetTransactionHeadId := place.TransactionHeadId == gateWebhook.TransactionHeadId
		if isPlaceStatusBilling && isTargetTransactionHeadId {
			newPlace = *place
			newPlace.PlaceStatus = domain.PlaceStatusCleaning
			isExist = true
			break
		}
	}

	if !isExist {
		return echoContext.JSON(
			http.StatusNotFound,
			nil,
		)
	}

	if err := toController.EditPlace(ctx, newPlace); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusInternalServerError,
			err,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)

}
