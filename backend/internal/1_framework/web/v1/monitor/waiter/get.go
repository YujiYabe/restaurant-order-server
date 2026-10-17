package waiter

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

type (
	AllTableList map[string][]string
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

	placeList := toController.GetPlaceList(ctx)
	allTableList := AllTableList{}

	for _, place := range placeList {
		allTableList[place.PlaceGroup] = append(allTableList[place.PlaceGroup], place.PlaceName)
	}

	allTableListJson, err := json.Marshal(
		allTableList,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return err
	}

	envUnique := env.Env.EnvUnique
	data := struct {
		Environment  string
		URL          string
		Version      string
		AllTableList string
	}{
		Environment:  envUnique.Environment,
		URL:          "https://" + envUnique.URL,
		Version:      envUnique.Version,
		AllTableList: string(allTableListJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"waiter",
		data,
	)
}
