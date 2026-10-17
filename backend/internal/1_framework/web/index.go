package web

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/labstack/echo"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

type indexLink struct {
	Name        string
	Description string
	Icon        string
	URL         string
}

type indexPlace struct {
	Name           string
	BackYardURL    string
	WaiterURL      string
	MobileOrderURL string
}

type indexPlaceGroup struct {
	Name   string
	Places []indexPlace
}

type indexPage struct {
	StoreName   string
	Version     string
	MainLinks   []indexLink
	AdminLinks  []indexLink
	PlaceGroups []indexPlaceGroup
}

// getIndex は、各店舗画面への入口を表示します。
func (receiver *Web) getIndex(echoContext echo.Context) (err error) {
	ctx := pkg.GetNewContext(
		echoContext.Request().Context(),
		echoContext.Response().Header().Get(echo.HeaderXRequestID),
	)
	deviceKey := resolveIndexDeviceKey(
		env.Env.Environment,
		echoContext.QueryParam("device_key"),
		env.Env.FromDeviceKey,
	)
	data := buildIndexPage(receiver.Controller.GetPlaceList(ctx), deviceKey)
	return echoContext.Render(http.StatusOK, "index", data)
}

func resolveIndexDeviceKey(environment string, requestedKey string, configuredKey string) string {
	if environment == env.LOCAL && requestedKey == "" {
		return configuredKey
	}
	return requestedKey
}

func buildIndexURL(environment string, configuredURL string, port string) string {
	if environment == env.LOCAL {
		return "https://localhost:" + port + "/"
	}

	baseURL := strings.TrimRight(configuredURL, "/")
	if !strings.HasPrefix(baseURL, "https://") && !strings.HasPrefix(baseURL, "http://") {
		baseURL = "https://" + baseURL
	}
	return baseURL + "/"
}

func buildIndexPage(placeList domain.PlaceList, deviceKey string) indexPage {
	deviceQuery := "?device_key=" + url.QueryEscape(deviceKey)
	data := indexPage{
		StoreName: env.Env.EnvUnique.StoreName,
		Version:   env.Env.EnvUnique.Version,
		MainLinks: []indexLink{
			{Name: "Reception", Description: "受付・席状況", Icon: "fa-door-open", URL: "/v1/monitor/reception" + deviceQuery},
			{Name: "Waiter", Description: "ウェイター用席一覧", Icon: "fa-concierge-bell", URL: "/v1/monitor/waiter" + deviceQuery},
			{Name: "Kitchen", Description: "厨房全体モニター", Icon: "fa-utensils", URL: "/v1/monitor/repository" + deviceQuery},
			{Name: "Signboard", Description: "店頭サイネージ", Icon: "fa-tv", URL: "/v1/monitor/signboard"},
			{Name: "Menu", Description: "公開メニュー", Icon: "fa-utensils", URL: "/v1/monitor/public"},
		},
		AdminLinks: []indexLink{
			{Name: "Products", Description: "商品・在庫管理", Icon: "fa-boxes", URL: "/v1/admin/product" + deviceQuery},
			{Name: "Age Check", Description: "成人確認", Icon: "fa-id-card", URL: "/v1/alcohol/verify_twenty"},
			{Name: "Performance", Description: "サーバー統計", Icon: "fa-chart-line", URL: "/debug/statsviz/"},
		},
	}

	placesByGroup := make(map[string][]*domain.Place)
	for _, place := range placeList {
		if place != nil && place.PlaceName != "" {
			placesByGroup[place.PlaceGroup] = append(placesByGroup[place.PlaceGroup], place)
		}
	}

	groupNames := make([]string, 0, len(placesByGroup))
	for groupName := range placesByGroup {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		places := placesByGroup[groupName]
		sort.Slice(places, func(i int, j int) bool {
			return places[i].PlaceName < places[j].PlaceName
		})
		placeGroup := indexPlaceGroup{Name: groupName}
		for _, place := range places {
			placeName := place.PlaceName
			escapedName := url.PathEscape(placeName)
			placeGroup.Places = append(placeGroup.Places, indexPlace{
				Name:           placeName,
				BackYardURL:    "/v1/monitor/back_yard/" + escapedName + deviceQuery,
				WaiterURL:      "/v1/monitor/waiter/place/" + escapedName + deviceQuery,
				MobileOrderURL: "/v1/monitor/internet/" + escapedName + "?permission=" + url.QueryEscape(place.Permission),
			})
		}
		data.PlaceGroups = append(data.PlaceGroups, placeGroup)
	}

	return data
}
