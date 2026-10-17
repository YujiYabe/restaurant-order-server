package firestore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	fs "cloud.google.com/go/firestore"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/api/iterator"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/1_framework/transport"
	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

const (
	timeFormat = "06-01-02-15:04:05"
)

type SaveOrderRequest struct {
	OrderList               string
	EnterTime               string
	ExitTime                string
	BaseCount               int
	ReqTransactionTemporary domain.ReqPostTransactionTemporary
}

// SaveOrder は、会計対象の注文情報をFirestoreへ保存します。
func (receiver *Firestore) SaveOrder(
	ctx context.Context,
	place domain.Place,
	reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
) (
	err error,
) {
	orderListJSON, err := transport.MarshalJSON(place.OrderList)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	orderListString := string(orderListJSON)
	saveOrderRequest := SaveOrderRequest{
		OrderList:               orderListString,
		EnterTime:               place.EnterTime.Format(timeFormat),
		ExitTime:                time.Now().Format(timeFormat),
		BaseCount:               place.BaseCount,
		ReqTransactionTemporary: reqPostTransactionTemporary,
	}
	placeName := place.PlaceName
	t := time.Now()
	docRef := receiver.
		ConnOrder.
		Collection(t.Format("2006-01-02")).
		Doc(placeName + "-" + strconv.Itoa(reqPostTransactionTemporary.TerminalTranId))

	_, err = docRef.Set(ctx, saveOrderRequest)
	if err != nil {
		var ae *apierror.APIError
		if errors.As(err, &ae) {
			log.Println(ae)
			log.Println(ae.Details().ErrorInfo)

			log.Println(ae.Reason())
			log.Println(ae.Details().Help.GetLinks())
		}
		pkg.Logging(ctx, err)

		// 	return
	}

	return
}

// 年齢認証ができればfirestoreにdocumentを作成
// 個人の購入履歴に利用
// firestore構造
// document:uuid ┬ attribute───┬ 証明書種類 CertificateType:driver/mynumber
//               │                └ 一意な番号 certificateID
//               └ YYYYMMDDHHMMSS ┬ 画像ID:uuid
//                                 ├ JANcode
//                                 └ count

type ClientAttribute struct {
	CertificateID   string
	CertificateType string
}

type ClientImage struct {
	Base64Source string
	TimeStamp    time.Time
}

// GetClientImage は、利用者UUIDに紐づく最新の顔画像をFirestoreから取得します。
func (receiver *Firestore) GetClientImage(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	imageList []string,
	err error,
) {
	iter := receiver.
		ConnCertificate.
		Collection(alcohol.UUID).
		// Doc("-attribute").
		OrderBy("TimeStamp", fs.Desc). // 降順にソート
		Limit(3).                      // 3件に制限
		Documents(ctx)
	defer iter.Stop()

	// 結果の処理
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			pkg.Logging(ctx, err)
			return imageList, fmt.Errorf("iterate client images: %w", err)
		}

		var clientImage ClientImage
		err = doc.DataTo(&clientImage)
		if err != nil {
			pkg.Logging(ctx, err)
			return imageList, fmt.Errorf("map client image: %w", err)
		}

		imageList = append(
			imageList,
			clientImage.Base64Source,
		)
		fmt.Printf("Document data: %#v\n", doc.Data())
	}

	return
}

// SaveClientAttribute は、本人確認書類の属性を利用者UUIDに紐づけてFirestoreへ保存します。
func (receiver *Firestore) SaveClientAttribute(
	ctx context.Context,
	alcohol domain.Alcohol,
	certificateID string,
) (
	err error,
) {
	clientAttribute := ClientAttribute{
		CertificateID:   certificateID,
		CertificateType: alcohol.CertificateType,
	}

	docRef := receiver.
		ConnCertificate.
		Collection(alcohol.UUID).
		Doc("-attribute")

	_, err = docRef.Set(ctx, clientAttribute)
	if err != nil {
		var ae *apierror.APIError
		if errors.As(err, &ae) {
			log.Println(ae)
			log.Println(ae.Details().ErrorInfo)

			log.Println(ae.Reason())
			log.Println(ae.Details().Help.GetLinks())
		}
		pkg.Logging(ctx, err)
	}

	err = receiver.RegisterFace(ctx, alcohol)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	return
}

// RegisterFace は、許可済み利用者の顔画像をFirestoreへ登録します。
func (receiver *Firestore) RegisterFace(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	err error,
) {

	now := time.Now()

	clientImage := ClientImage{
		Base64Source: alcohol.Base64Source,
		TimeStamp:    now,
	}

	docRef := receiver.
		ConnCertificate.
		Collection(alcohol.UUID).
		Doc(now.Format("20060102_150405"))

	_, err = docRef.Set(ctx, clientImage)
	if err != nil {
		var ae *apierror.APIError
		if errors.As(err, &ae) {
			log.Println(ae)
			log.Println(ae.Details().ErrorInfo)

			log.Println(ae.Reason())
			log.Println(ae.Details().Help.GetLinks())
		}

		pkg.Logging(ctx, err)
		return
	}

	return
}

type SaveClientOrderRequest struct {
	ImageUUID string
	JANCode   int
	Count     int
}

// SaveClientOrder は、利用者の酒類注文履歴をFirestoreへ保存します。
func (receiver *Firestore) SaveClientOrder(
	ctx context.Context,
	clientUUID string,
	imageUUID string,
	janCode int,
	count int,
) (
	err error,
) {

	timeNow := time.Now()
	saveClientOrderRequest := SaveClientOrderRequest{
		ImageUUID: imageUUID,
		JANCode:   janCode,
		Count:     count,
	}

	docRef := receiver.
		ConnCertificate.
		Collection(clientUUID).
		Doc(timeNow.Format("2006-01-02-15-04-05"))

	_, err = docRef.Set(ctx, saveClientOrderRequest)
	if err != nil {
		var ae *apierror.APIError
		if errors.As(err, &ae) {
			log.Println(ae)
			log.Println(ae.Details().ErrorInfo)

			log.Println(ae.Reason())
			log.Println(ae.Details().Help.GetLinks())
		}
		pkg.Logging(ctx, err)
	}

	return
}
