package back_yard

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
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
	placeName := echoContext.Param("place_name")

	miniProductMapJson, err := transport.MarshalJSON(
		toController.GetMiniProductMap(ctx),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	envUnique := env.Env.EnvUnique
	data := struct {
		PlaceName      string
		Environment    string
		Version        string
		URL            string
		MiniProductMap string
	}{
		PlaceName:      placeName,
		Environment:    envUnique.Environment,
		URL:            "https://" + envUnique.URL,
		Version:        envUnique.Version,
		MiniProductMap: string(miniProductMapJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"backYard",
		data,
	)
}
