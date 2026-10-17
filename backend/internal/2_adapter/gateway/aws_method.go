package gateway

import (
	"context"
)

// Identification は、2枚の画像に写る顔が同一人物かAWS Rekognitionで判定します。
func (receiver *Gateway) Identification(
	ctx context.Context,
	imageSet []string,
) (
	ok bool,
	err error,
) {
	ok, err = receiver.ToAws.Identification(
		ctx,
		imageSet,
	)

	return
}
