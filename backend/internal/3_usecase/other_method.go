package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// Start は、必要な初期処理とバックグラウンド処理を開始します。
func (receiver *UseCase) Start(
	ctx context.Context,
) {
	// 初期処理
	receiver.SetUpInMemory(ctx) // DBをインメモリに保存
	envUnique := env.Env.EnvUnique

	if envUnique.IsEnableCharge {
		go receiver.DetectAutoUpdate(ctx) // チャージ料金の時間更新処理
	}
}

// DetectAutoUpdate は、チャージ料金の自動更新通知を監視して最新状態を配信します。
func (receiver *UseCase) DetectAutoUpdate(
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case isAutoUpdate := <-domain.IsAutoUpdate:
			if isAutoUpdate {
				receiver.DistributeAgent(ctx)
			}
		}
	}
}

// GetAllergyDefault は、既定のアレルギー設定のコピーを返します。
func (receiver *UseCase) GetAllergyDefault(
	ctx context.Context,
) (
	allergy domain.Allergy,
) {
	allergy = receiver.ToDomain.GetAllergyDefault(
		ctx,
	)

	return
}

// DistributeAgent は、現在のドメイン状態を永続化し、接続中の画面へ配信します。
func (receiver *UseCase) DistributeAgent(
	ctx context.Context,
) {
	domainDomain := receiver.ToDomain.GetDomain(
		ctx,
	)

	// redisに現在の状態を保存。エラーは無視する
	placeList := receiver.ToDomain.GetPlaceList(ctx)
	if err := receiver.ToGateway.SetPlaceListToRedis(
		ctx,
		placeList,
	); err != nil {
		pkg.Logging(ctx, err)
	}

	receiver.ToPresenter.DistributeAgent(
		ctx,
		domainDomain,
	)
}

// GetTagList は、商品タグ一覧のコピーを返します。
func (receiver *UseCase) GetTagList(
	ctx context.Context,
) (
	tagPrimaryList domain.TagPrimaryList,
) {
	tagPrimaryList = receiver.ToDomain.GetTagList(
		ctx,
	)

	return
}

// GetInventoryMap は、現在の在庫数マップのコピーを返します。
func (receiver *UseCase) GetInventoryMap(
	ctx context.Context,
) (
	inventoryMap domain.InventoryMap,
) {
	inventoryMap = receiver.ToDomain.GetInventoryMap(
		ctx,
	)

	return
}

// VerifyTwenty は、書類から取得した生年月日を検証し、成人であることを確認します。
func (receiver *UseCase) VerifyTwenty(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	certificateID string,
	age int,
	birthday time.Time,
	err error,
) {

	// デバッグ用
	if false {
		err = generateFile(
			"/home/yuji/Downloads/s.txt",
			[]byte(alcohol.Base64Source),
		)

		if err != nil {
			pkg.Logging(ctx, err)
			return
		}

		// // 画像をローカル保存
		data, decodeErr := pkg.DecodeBase64DataURL(alcohol.Base64Source)
		if decodeErr != nil {
			pkg.Logging(ctx, decodeErr)
			pkg.Logging(ctx, decodeErr.Error())
			err = decodeErr
			return
		}

		err = generateFile(
			"/home/yuji/Downloads/s.jpg",
			data,
		)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
	}

	birthdayString, certificateID, err := receiver.ToGateway.VerifyCertificate(
		ctx,
		alcohol,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	if birthdayString == "" {
		err = fmt.Errorf("birthdayString is empty")
		pkg.Logging(ctx, err)
		return
	}

	age, birthday, err = receiver.ToDomain.VerifyTwenty(
		ctx,
		birthdayString,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// 20歳未満であれば終了
	if 20 > age {
		err = fmt.Errorf("under 20 years old")
		pkg.Logging(ctx, err)
		return
	}

	err = receiver.ToGateway.SaveClientAttribute(
		ctx,
		alcohol,
		certificateID,
	)

	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// 書類の成人確認後に限り、1回だけ顔画像を登録できる。
	receiver.ToDomain.AllowFaceRegistration(ctx, alcohol)

	return
}

// RegisterFace は、許可済み利用者の顔画像をFirestoreへ登録します。
func (receiver *UseCase) RegisterFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {
	if !receiver.ToDomain.ConsumeFaceRegistration(ctx, alcohol) {
		return fmt.Errorf("face registration is not authorized or has expired")
	}

	err = receiver.registerFace(ctx, alcohol)
	if err != nil {
		receiver.ToDomain.AllowFaceRegistration(ctx, alcohol)
	}
	return
}

// registerFace は、画像形式を検証して顔画像を永続化します。
func (receiver *UseCase) registerFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {
	if _, err = pkg.DecodeBase64DataURL(alcohol.Base64Source); err != nil {
		return
	}
	if false {
		// Base64Sourceをテキスト保存
		err = generateFile(
			"/home/yuji/Downloads/s.txt",
			[]byte(alcohol.Base64Source),
		)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}

		// 画像をローカル保存
		// Base64文字列をバイト列へデコードします。
		data, decodeErr := pkg.DecodeBase64DataURL(alcohol.Base64Source)
		if decodeErr != nil {
			pkg.Logging(ctx, decodeErr)
			pkg.Logging(ctx, decodeErr.Error())
			err = decodeErr
			return
		}

		// データをファイルへ書き込みます。
		err = generateFile(
			"/home/yuji/Downloads/s.jpg",
			data,
		)
		if err != nil {
			pkg.Logging(ctx, err)
			return
		}
	}

	// firestoreにbase64顔画像登録する
	err = receiver.ToGateway.RegisterFace(
		ctx,
		alcohol,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	return
}

// VerifyFace は、登録済み画像との顔照合に成功した利用者へ当日の酒類注文を許可します。
func (receiver *UseCase) VerifyFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {
	// firestoreから対象UUIDの最新3件のbase64Imageを取得
	imageList, err := receiver.ToGateway.GetClientImage(
		ctx,
		alcohol,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return err
	}
	if len(imageList) == 0 {
		err = fmt.Errorf("no exist image")
		pkg.Logging(ctx, err)
		return
	}

	// awsで3件 顔画像比較
	isVerified := false
	for _, image := range imageList {
		ok, err := receiver.ToGateway.Identification(
			ctx,
			[]string{
				alcohol.Base64Source,
				image,
			},
		)
		if err != nil {
			pkg.Logging(ctx, err)
		}

		if ok {
			isVerified = true
			break
		}
	}

	if !isVerified {
		err = fmt.Errorf("can not verify")
		pkg.Logging(ctx, err)
		return
	}

	// 新しい画像を登録
	err = receiver.registerFace(
		ctx,
		alcohol,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// 当日の許可リストに追加
	receiver.ToDomain.AddTwentyList(
		ctx,
		alcohol,
	)

	return
}

// generateFile は、指定された内容を権限0600のファイルとして保存します。
func generateFile(
	filePath string,
	bytes []byte,
) (
	err error,
) {
	err = os.WriteFile(
		filePath,
		bytes,
		0600,
	)

	return
}
