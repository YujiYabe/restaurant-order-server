package verify_twenty

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

// get は、GETリクエストを検証し、対象データまたは画面を応答します。
func get(
	echoContext echo.Context,
	_ controller.ToController,
) (
	err error,
) {

	envUnique := env.Env.EnvUnique
	data := struct {
		StoreID         string
		Environment     string
		URL             string
		StoreName       string
		StoreAddress    string
		Commission      float64
		ChargeBasePrice int
		IsPermitAlcohol bool
		Version         string
	}{
		StoreID:         envUnique.StoreID,
		Environment:     envUnique.Environment,
		URL:             "https://" + envUnique.URL,
		StoreName:       envUnique.StoreName,
		StoreAddress:    envUnique.StoreAddress,
		Commission:      envUnique.Commission,
		ChargeBasePrice: envUnique.ChargeBasePrice,
		IsPermitAlcohol: envUnique.IsPermitAlcohol,
		Version:         envUnique.Version,
	}

	return echoContext.Render(
		http.StatusOK,
		"verifyTwenty",
		data,
	)
}
