package admin_product

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// ---------------------------------------------------------
// 管理画面 html
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

	allProductListJson, err := transport.MarshalJSON(
		toController.GetAllProductList(ctx),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	allergyListJson, err := transport.MarshalJSON(
		toController.GetAllergyDefault(ctx),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return echoContext.JSON(
			http.StatusBadRequest,
			err,
		)
	}

	inventoryMapJson, err := transport.MarshalJSON(
		toController.GetInventoryMap(ctx),
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
		StoreID        string
		Environment    string
		URL            string
		StoreName      string
		StoreAddress   string
		Version        string
		AllProductList string
		AllergyList    string
		InventoryMap   string
	}{
		StoreID:        envUnique.StoreID,
		Environment:    envUnique.Environment,
		URL:            "https://" + envUnique.URL,
		StoreName:      envUnique.StoreName,
		StoreAddress:   envUnique.StoreAddress,
		Version:        envUnique.Version,
		AllProductList: string(allProductListJson),
		AllergyList:    string(allergyListJson),
		InventoryMap:   string(inventoryMapJson),
	}

	return echoContext.Render(
		http.StatusOK,
		"adminProduct",
		data,
	)
}
