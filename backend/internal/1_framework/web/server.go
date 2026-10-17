package web

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/arl/statsviz"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	webUtil "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/util"
	v1 "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
)

type (
	// Webサーバー
	Web struct {
		EchoEcho   *echo.Echo
		Controller controller.ToController
	}
)

// NewWeb は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewWeb(
	ctrl controller.ToController,
) (
	web *Web,
) {
	web = &Web{
		EchoEcho:   NewEcho(),
		Controller: ctrl,
	}

	return
}

type Template struct {
	templates *template.Template
}

// Render は、指定されたテンプレートをデータとともにHTTPレスポンスへ描画します。
func (receiver *Template) Render(
	w io.Writer,
	name string,
	data interface{},
	echoContext echo.Context,
) (
	err error,
) {

	err = receiver.templates.ExecuteTemplate(
		w,
		name,
		data,
	)

	return
}

// NewEcho は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewEcho() (
	echoEcho *echo.Echo,
) {
	echoEcho = echo.New()
	echoEcho.HideBanner = true
	t := &Template{
		templates: template.Must(
			template.ParseGlob("web/*.html"),
		),
	}
	echoEcho.Renderer = t
	echoEcho.Static("/js", "web/js")
	echoEcho.Static("/css", "web/css")
	echoEcho.Static("/image", "web/image")
	echoEcho.Static("/sound", "web/sound")
	echoEcho.Static("/fontawesome", "web/fontawesome")

	echoEcho.Use(
		middleware.LoggerWithConfig(
			middleware.LoggerConfig{
				Format:           "${time_custom}__${status}__${method}__${uri}\n",
				CustomTimeFormat: "06/01/02-15:04:05",
			},
		),
	)

	echoEcho.Use(middleware.Recover())
	echoEcho.Use(middleware.RequestID())

	// ガラケーからのアクセス禁止
	// レイアウトやSSEが保証できない
	echoEcho.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) (
			err error,
		) {

			userAgent := echoContext.Request().Header.Get("User-Agent")

			isDoCoMo := strings.Contains(userAgent, "DoCoMo")
			isKDDI := strings.Contains(userAgent, "UP.Browser")
			isVodafone := strings.Contains(userAgent, "Vodafone")
			isSoftBank := strings.Contains(userAgent, "SoftBank")

			if isDoCoMo || isKDDI || isVodafone || isSoftBank {
				return echoContext.String(http.StatusForbidden, "Access from this device is not allowed")
			}

			return next(echoContext)
		}
	})

	return echoEcho
}

// Start は、APIルートと診断画面を登録し、設定された証明書でHTTPSサーバーを起動します。
func (receiver *Web) Start() {
	group := receiver.EchoEcho.Group("")
	indexMiddleware := make([]echo.MiddlewareFunc, 0, 2)
	if env.Env.Environment != env.LOCAL {
		indexMiddleware = append(
			indexMiddleware,
			middleware.BasicAuth(webUtil.CustomBasicAuth),
			webUtil.DeviceKeyAuth(),
		)
	}
	group.GET(
		"/",
		receiver.getIndex,
		indexMiddleware...,
	)

	v1.NewRoute(
		receiver.EchoEcho,
		receiver.Controller,
		group,
	)

	isShowPerformance := true
	if isShowPerformance {
		mux := http.NewServeMux()
		err := statsviz.Register(mux)
		if err != nil {
			log.Printf("%#v\n", err)
		}

		group.GET("/debug/statsviz/",
			echo.WrapHandler(mux),
			middleware.BasicAuth(webUtil.CustomBasicAuth),
		)
		group.GET("/debug/statsviz/*",
			echo.WrapHandler(mux),
			middleware.BasicAuth(webUtil.CustomBasicAuth),
		)
	}

	isShowRoute := false
	if isShowRoute {
		routes := receiver.EchoEcho.Routes()
		for _, route := range routes {
			log.Printf("%#v\n", route)
		}
	}

	isValidSSL := true
	if isValidSSL {
		port := ":" + env.Env.ServerPort
		var certFile string
		var keyFile string

		switch env.Env.Environment {
		case env.LOCAL:
			certFile = "openssl/server.crt"
			keyFile = "openssl/server.key"

		default:
			certFile = "/etc/letsencrypt/live/" + env.Env.EnvUnique.CertPath + "/fullchain.pem"
			keyFile = "/etc/letsencrypt/live/" + env.Env.EnvUnique.CertPath + "/privkey.pem"
		}

		log.Printf(
			"restaurant-order-server Portal: %s\n",
			buildIndexURL(env.Env.Environment, env.Env.EnvUnique.URL, env.Env.ServerPort),
		)
		receiver.EchoEcho.Logger.Fatal(
			receiver.EchoEcho.StartTLS(
				port,
				certFile,
				keyFile,
			),
		)

	} else {
		err := receiver.EchoEcho.Start(":6789")

		if err != nil {
			log.Printf("%#v\n", err)
		}
	}

}
