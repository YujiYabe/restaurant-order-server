package webUtil

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

// CustomBasicAuth は、実行環境と設定値に基づいてBasic認証情報を検証します。
func CustomBasicAuth(
	userName string,
	password string,
	echoContext echo.Context,
) (
	ok bool,
	err error,
) {
	// localは認証パス
	if env.Env.Environment == env.LOCAL {
		ok = true
		return
	}

	isUserName := userName == env.Env.BasicAuthUserName
	isPassword := password == env.Env.BasicAuthPassword
	if isUserName && isPassword {
		ok = true
		return
	}

	return
}

// DeviceKeyAuth は、端末キーが設定値と一致するリクエストだけを許可するミドルウェアを返します。
func DeviceKeyAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) (
			err error,
		) {
			if echoContext.QueryParam("device_key") == env.Env.FromDeviceKey {
				return next(echoContext)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "")
		}
	}
}

// SmaregiKeyAuth は、スマレジ連携用APIキーを検証するミドルウェアを返します。
func SmaregiKeyAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) (
			err error,
		) {
			apiKey := echoContext.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "X-API-Key is missing")
			}

			if apiKey != env.Env.FromGateKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "X-API-Key is missing")
			}

			return next(echoContext)
		}
	}
}
