package redis

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v2"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// GetPlaceListFromRedis は、Redisから席一覧を取得し、未保存の場合は設定ファイルから初期化します。
func (receiver *Redis) GetPlaceListFromRedis(
	ctx context.Context,
) (
	placeList domain.PlaceList,
	err error,
) {
	placeList = domain.PlaceList{}
	res, err := receiver.Conn.Get(
		ctx,
		"placeList",
	).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			pkg.Logging(ctx, "Key not found")
			placeList, err = newPlaceList(ctx)
			return
		}
		pkg.Logging(ctx, err)
		return
	}

	var stored transport.PlaceList
	if err = json.Unmarshal([]byte(res), &stored); err != nil {
		pkg.Logging(ctx, err)
		return
	}
	placeList = stored.Domain()

	return
}

// SetPlaceListToRedis は、席一覧をRedisへ永続化します。
func (receiver *Redis) SetPlaceListToRedis(
	ctx context.Context,
	placeList domain.PlaceList,
) (
	err error,
) {
	data, err := json.Marshal(
		transport.PlaceListFromDomain(placeList),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return err
	}
	err = receiver.Conn.Set(
		ctx,
		"placeList",
		data,
		0,
	).Err()
	if err != nil {
		pkg.Logging(ctx, err)
	}
	return
}

// ResetPlaceListInRedis は、Redisに保存された席一覧を削除します。
func (receiver *Redis) ResetPlaceListInRedis(
	ctx context.Context,
) (
	err error,
) {
	_, err = receiver.Conn.Del(
		ctx,
		"placeList",
	).Result()
	if err != nil {
		pkg.Logging(ctx, err)
	}

	return
}

// newPlaceList は、店舗別YAML設定から初期席一覧を生成します。
func newPlaceList(
	ctx context.Context,
) (
	placeList domain.PlaceList,
	err error,
) {
	placeList = domain.PlaceList{}
	data, err := os.ReadFile(
		"config/" + env.Env.StoreID + "_table_list.yaml",
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	var configured transport.PlaceList
	err = yaml.Unmarshal(
		data,
		&configured,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	placeList = configured.Domain()
	for index := range placeList {
		placeList[index].OrderList = make(domain.OrderList)
	}

	return
}
