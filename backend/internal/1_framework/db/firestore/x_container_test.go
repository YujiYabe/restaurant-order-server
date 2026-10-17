package firestore

import (
	"context"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestNewClientRejectsEmulatorURLScheme は、EmulatorホストにURLスキームを指定した場合に拒否することを確認します。
func TestNewClientRejectsEmulatorURLScheme(t *testing.T) {
	t.Setenv("FIRESTORE_EMULATOR_HOST", "http://127.0.0.1:8080")

	client, err := newClient(context.Background(), "TEST_PROJECT_ID", "demo-test", "unused.json")
	if err == nil {
		if client != nil {
			_ = client.Close()
		}
		t.Fatal("newClient() accepted FIRESTORE_EMULATOR_HOST with a URL scheme")
	}
}

// TestNewClientWithEmulator は、認証情報なしでEmulatorへ接続でき、プロジェクト間でデータが分離されることを確認します。
func TestNewClientWithEmulator(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	ctx := context.Background()
	const (
		collection = "integration-test"
		document   = "project-isolation"
	)

	orderClient, err := newClient(ctx, "TEST_ORDER_PROJECT_ID", "demo-restaurant-order-server-order", "unused.json")
	if err != nil {
		t.Fatalf("create order client: %v", err)
	}
	defer func() {
		_ = orderClient.Close()
	}()

	certificateClient, err := newClient(ctx, "TEST_CERTIFICATE_PROJECT_ID", "demo-restaurant-order-server-certificate", "unused.json")
	if err != nil {
		t.Fatalf("create certificate client: %v", err)
	}
	defer func() {
		_ = certificateClient.Close()
	}()

	orderDocument := orderClient.Collection(collection).Doc(document)
	if _, err = orderDocument.Set(ctx, map[string]string{"source": "order"}); err != nil {
		t.Fatalf("write order document: %v", err)
	}
	t.Cleanup(func() {
		_, _ = orderDocument.Delete(context.Background())
	})

	_, err = certificateClient.Collection(collection).Doc(document).Get(ctx)
	if status.Code(err) != codes.NotFound {
		t.Fatalf("certificate project Get() error = %v, want NotFound", err)
	}
}
