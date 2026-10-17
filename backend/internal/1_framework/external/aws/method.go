package aws

import (
	"context"
	"fmt"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

// Identification は、2枚の画像に写る顔が同一人物かAWS Rekognitionで判定します。
func (receiver *Aws) Identification(
	ctx context.Context,
	imageSet []string,
) (
	ok bool,
	err error,
) {
	if len(imageSet) != 2 {
		err = fmt.Errorf("image set must contain exactly two images")
		pkg.Logging(ctx, err)
		return
	}

	sourceImageBytes, err := pkg.DecodeBase64DataURL(imageSet[0])
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	targetImageBytes, err := pkg.DecodeBase64DataURL(imageSet[1])
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		err = fmt.Errorf("create AWS session: %w", err)
		pkg.Logging(ctx, err)
		return
	}

	// Rekognitionクライアントの作成
	svc := rekognition.New(sess)

	// 画像比較のリクエストの作成
	input := &rekognition.CompareFacesInput{
		SourceImage: &rekognition.Image{
			Bytes: sourceImageBytes,
		},
		TargetImage: &rekognition.Image{
			Bytes: targetImageBytes,
		},
		SimilarityThreshold: aws.Float64(80.0),
	}

	// 画像の比較を実行
	result, err := svc.CompareFaces(input)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// 結果の表示
	for _, match := range result.FaceMatches {
		if match != nil {
			if match.Similarity != nil {
				if *match.Similarity >= 80 {
					ok = true
					return
				}
			}
		}
	}

	return
}
