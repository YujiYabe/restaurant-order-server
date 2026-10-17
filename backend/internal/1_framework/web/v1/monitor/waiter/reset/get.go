package reset

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

// get は、GETリクエストを検証し、対象データまたは画面を応答します。
func get(
	echoContext echo.Context,
	_ controller.ToController,
) (
	err error,
) {

	return echoContext.Render(
		http.StatusOK,
		"reset",
		nil,
	)
}
