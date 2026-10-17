package admin_product

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// ---------------------------------------------------------
// 管理画面 商品更新
// Patch は管理者用のAPIで、製品情報を更新します。
func patch(
	echoContext echo.Context,
	toController controller.ToController,
) (
	err error,
) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)

	// コンテキストからJANコードを取得
	number, err := strconv.Atoi(echoContext.Param("number"))
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	// // リクエストの内容を新製品オブジェクトにバインド
	request := transport.Product{}
	if err := echoContext.Bind(&request); err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}
	newProduct := request.Domain()

	newProduct.JANCode = number

	err = toController.UpdateProduct(
		ctx,
		newProduct,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	return echoContext.JSON(
		http.StatusOK,
		nil,
	)
}
