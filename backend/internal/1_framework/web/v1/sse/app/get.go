package app

import (
	// #nosec
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"context"
	"encoding/json"
	"log"
	"sort"
	"sync"

	"github.com/labstack/echo"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	ssechannel "github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/web/v1/sse/channel"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/controller"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

type Agent struct {
	mutex     sync.Mutex
	ID        string
	Writer    http.ResponseWriter
	Flusher   http.Flusher
	Kind      string
	PlaceName string
	Hash      [16]byte
}

var (
	Agents = make(map[string]*Agent)
	Mutex  sync.RWMutex
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

	kind := echoContext.QueryParam("kind")
	placeName := ""
	if kind == "back_yard" || kind == "waiter_place" || kind == "internet" {
		placeName = echoContext.QueryParam("place_name")
		if placeName == "" || toController.GetPlace(ctx, placeName).PlaceName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "place not found")
		}
	} else if kind != "repository" && kind != "reception" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid agent kind")
	}

	echoContext.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	echoContext.Response().Header().Set("Cache-Control", "no-cache")
	echoContext.Response().Header().Set("Connection", "keep-alive")

	flusher, ok := echoContext.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	id := echoContext.Response().Header().Get(echo.HeaderXRequestID)
	agent := &Agent{
		ID:        id,
		Writer:    echoContext.Response().Writer,
		Flusher:   flusher,
		Kind:      kind,
		PlaceName: placeName,
	}

	Mutex.Lock()
	Agents[id] = agent
	Mutex.Unlock()

	toController.DistributeAgent(ctx)

	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	// 接続を維持（簡易的な keep-alive）
	for {
		select {
		case <-ctx.Done():
			Mutex.Lock()
			delete(Agents, id)
			Mutex.Unlock()
			return nil
		case <-heartbeat.C:
			// 接続維持用の何もしないpingです。
			agent.mutex.Lock()
			_, writeErr := agent.Writer.Write([]byte(": ping\n\n"))
			if writeErr != nil {
				agent.mutex.Unlock()
				removeAgent(id)
				pkg.Logging(ctx, writeErr)
				return nil
			}
			agent.Flusher.Flush()
			agent.mutex.Unlock()
		}
	}
}

type clientResponse struct {
	Place        transport.Place
	InventoryMap domain.InventoryMap
}

type ReceptionResponseList []ReceptionResponse

type ReceptionResponse struct {
	PlaceGroup  string `json:"place_group"`
	PlaceName   string `json:"place_name"`
	Capacity    int    `json:"capacity"`
	CustomerId  int    `json:"customer_id"`
	LineNumber  int    `json:"line_number"`
	PlaceStatus int    `json:"place_status"`
}

// SendToAgentPlaceList は、状態更新を受信し、接続中の各SSEクライアントへ振り分けます。
func SendToAgentPlaceList() {
	for {
		domainDomain := <-ssechannel.DomainChan

		sortOutAgent(domainDomain)
	}
}

// sortOutAgent は、接続種別に応じたレスポンスを構築して各SSEクライアントへ送信します。
func sortOutAgent(
	domainDomain *domain.Domain,
) {
	placeList := &domainDomain.PlaceList
	inventoryMap := domainDomain.InventoryMap

	placeNameList := make([]string, 0, len(*placeList))
	for key := range *placeList {
		placeNameList = append(
			placeNameList,
			key,
		)
	}
	sort.Strings(placeNameList)

	var receptionResponseList ReceptionResponseList
	for _, placeName := range placeNameList {
		if (*placeList)[placeName] == nil {
			continue
		}
		receptionResponseList = append(
			receptionResponseList,
			ReceptionResponse{
				PlaceGroup:  (*placeList)[placeName].PlaceGroup,
				PlaceName:   (*placeList)[placeName].PlaceName,
				Capacity:    (*placeList)[placeName].Capacity,
				CustomerId:  (*placeList)[placeName].CustomerId,
				LineNumber:  (*placeList)[placeName].LineNumber,
				PlaceStatus: int((*placeList)[placeName].PlaceStatus),
			},
		)
	}

	Mutex.RLock()
	agents := make([]*Agent, 0, len(Agents))
	for _, agent := range Agents {
		agents = append(agents, agent)
	}
	Mutex.RUnlock()

	for _, agent := range agents {
		switch agent.Kind {
		// ___________________________________
		case "back_yard":
			place := (*placeList)[agent.PlaceName]
			if place == nil {
				removeAgent(agent.ID)
				continue
			}
			sendToAgent(
				agent,
				transport.PlaceFromDomain(*place),
			)

		// ___________________________________
		case "waiter_place", "internet":
			place := (*placeList)[agent.PlaceName]
			if place == nil {
				removeAgent(agent.ID)
				continue
			}
			sendToAgent(
				agent,
				clientResponse{
					Place:        transport.PlaceFromDomain(*place),
					InventoryMap: inventoryMap,
				},
			)

		// ___________________________________
		case "repository":
			sendToAgent(
				agent,
				transport.PlaceListFromDomain(*placeList),
			)

		// ___________________________________
		case "reception":
			sendToAgent(
				agent,
				receptionResponseList,
			)

		default:
			log.Println("is unknown")
		}

	}
}

// sendToAgent は、変更がある場合だけSSEメッセージをクライアントへ送信します。
func sendToAgent(
	agent *Agent,
	content interface{},
) {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()
	ctx := pkg.GetNewContext(
		context.Background(), // 親コンテキスト
		agent.ID,             // 値
	)
	byteContent, err := json.Marshal(content)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// #nosec
	hash := md5.Sum(byteContent)
	if agent.Hash == hash {
		return
	}

	message := fmt.Sprintf("data: %s\n\n", byteContent)
	_, err = agent.Writer.Write([]byte(message))
	if err != nil {
		pkg.Logging(ctx, err)
		removeAgent(agent.ID)
		return
	}

	agent.Flusher.Flush()
	agent.Hash = hash
}

// removeAgent は、指定IDのSSEクライアントを接続一覧から削除します。
func removeAgent(id string) {
	Mutex.Lock()
	delete(Agents, id)
	Mutex.Unlock()
}
