package usecase

import (
	"context"
	"testing"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

type faceGatewayStub struct {
	ToGateway
	registerCalls int
}

// RegisterFace は、テストで使用する依存処理を代替します。
func (stub *faceGatewayStub) RegisterFace(
	context.Context,
	domain.Alcohol,
) (
	err error,
) {
	stub.registerCalls++
	return nil
}

// TestRegisterFaceRequiresOneTimeDocumentApproval は、対象仕様の回帰動作を検証します。
func TestRegisterFaceRequiresOneTimeDocumentApproval(t *testing.T) {
	domainModel := &domain.Domain{PendingFaceList: domain.NewTwentyList()}
	gateway := &faceGatewayStub{}
	useCase := &UseCase{ToDomain: domainModel, ToGateway: gateway}
	alcohol := domain.Alcohol{UUID: "uuid", Base64Source: "data:image/png;base64,AA=="}

	if err := useCase.RegisterFace(context.Background(), alcohol); err == nil {
		t.Fatal("unapproved face registration was accepted")
	}
	domainModel.AllowFaceRegistration(context.Background(), alcohol)
	if err := useCase.RegisterFace(context.Background(), alcohol); err != nil {
		t.Fatalf("approved face registration failed: %v", err)
	}
	if err := useCase.RegisterFace(context.Background(), alcohol); err == nil {
		t.Fatal("face registration approval was reused")
	}
	if gateway.registerCalls != 1 {
		t.Fatalf("RegisterFace gateway calls = %d, want 1", gateway.registerCalls)
	}
}
