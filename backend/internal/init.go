package internal

import (
	"log"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/db/firestore"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/db/redis"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/db/sqlite"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/external/aws"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/external/azure"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/external/gate"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/external/monitor"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/external/smaregi"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
)

type app struct {
	web     *web.Web
	monitor *monitor.Monitor
}

// NewApp は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewApp() (
	a *app,
) {

	newControllerOption := controller.NewControllerOption{
		ToSqlite:    sqlite.NewToSQLite(),
		ToRedis:     redis.NewToRedis(),
		ToFirestore: firestore.NewToFirestore(),
		ToAzure:     azure.NewToAzure(),
		ToAws:       aws.NewToAws(),
		ToSmaregi:   smaregi.NewToSmaregi(),
		ToGate:      gate.NewToGate(),
		ToMonitor:   monitor.NewToMonitor(),
	}

	ctrl := controller.NewController(
		newControllerOption,
	)
	ctrl.Start() // 初期処理 DBをインメモリに保存

	a = &app{
		web:     web.NewWeb(ctrl),
		monitor: monitor.NewMonitor(),
	}

	return
}

// Start は、必要な初期処理とバックグラウンド処理を開始します。
func (receiver *app) Start() {
	log.SetFlags(0)
	receiver.web.Start()
}
