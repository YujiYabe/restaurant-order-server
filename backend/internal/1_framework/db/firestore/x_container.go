package firestore

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

const (
	defaultOrderProjectID       = "demo-restaurant-order-server-order"
	defaultCertificateProjectID = "demo-restaurant-order-server-certificate"
)

type (
	// Firestore接続
	Firestore struct {
		ConnCertificate *firestore.Client
		ConnOrder       *firestore.Client
	}
)

// NewToFirestore は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToFirestore() (
	toFirestore gateway.ToFirestore,
) {
	ctx := context.Background()
	store := new(Firestore)

	connOrder, err := newClient(
		ctx,
		"FIRESTORE_ORDER_PROJECT_ID",
		defaultOrderProjectID,
		"./secret/firestore_order.json",
	)
	if err != nil {
		pkg.Logging(ctx, err)
		log.Fatalf("initialize order firestore: %v", err)
	}
	store.ConnOrder = connOrder

	connCertificate, err := newClient(
		ctx,
		"FIRESTORE_CERTIFICATE_PROJECT_ID",
		defaultCertificateProjectID,
		"./secret/firestore_certificate.json",
	)
	if err != nil {
		pkg.Logging(ctx, err)
		log.Fatalf("initialize certificate firestore: %v", err)
	}
	store.ConnCertificate = connCertificate

	toFirestore = store
	return
}

// newClient は、Emulator設定時は認証情報なしでローカルへ接続し、
// それ以外ではサービスアカウント認証情報を使用してFirestoreへ接続します。
func newClient(
	ctx context.Context,
	projectIDEnvironment string,
	defaultProjectID string,
	credentialsFile string,
) (*firestore.Client, error) {
	emulatorHost := strings.TrimSpace(os.Getenv("FIRESTORE_EMULATOR_HOST"))
	if emulatorHost != "" {
		if strings.Contains(emulatorHost, "://") {
			return nil, fmt.Errorf("FIRESTORE_EMULATOR_HOST must not include a URL scheme: %q", emulatorHost)
		}

		projectID := strings.TrimSpace(os.Getenv(projectIDEnvironment))
		if projectID == "" {
			projectID = defaultProjectID
		}

		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
		if err != nil {
			return nil, fmt.Errorf("create firebase app for emulator project %q: %w", projectID, err)
		}

		client, err := app.Firestore(ctx)
		if err != nil {
			return nil, fmt.Errorf("create firestore emulator client for project %q: %w", projectID, err)
		}

		return client, nil
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("create firebase app with credentials %q: %w", credentialsFile, err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("create firestore client with credentials %q: %w", credentialsFile, err)
	}

	return client, nil
}
