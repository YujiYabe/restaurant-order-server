# バックエンドE2Eテスト

起動済みのローカルrestaurant-order-serverへHTTP/SSEで接続し、実際のルーティング、認証、Redis永続化、Domain、SSE配信を通した動作を検証します。

## 安全制約

- `E2E_BASE_URL` は `localhost` またはloopback IPだけを許可します。
- テストは席情報をリセットして状態遷移を実行するため、開発専用のRedisとバックエンドを使用してください。
- 会計処理はスマレジ、Firestore、Gate、SQLite在庫更新を通ります。スマレジとGateはローカルスタブ、Firestoreはエミュレーター、SQLiteは破棄可能なテストDBを必ず使用してください。
- 誤って実サービスや通常の店舗DBを使わないよう、確認用の `E2E_LOCAL_BILLING_STUBS=true` がない場合はテストを停止します。
- テスト終了時に席情報を設定ファイルの初期状態へ戻します。

## 実行

バックエンドとRedisを起動し、バックエンドが使用している `FROM_DEVICE_KEY` と同じ値を指定します。

```bash
E2E_LOCAL_BILLING_STUBS=true \
E2E_DEVICE_KEY='local-device-key' \
E2E_GATE_KEY='local-gate-key' \
E2E_JAN_CODE='テストDBに在庫がある酒類以外のJANコード' \
make test-e2e
```

接続先の既定値は `https://127.0.0.1:6789` です。変更する場合もloopbackだけが利用できます。

```bash
E2E_BASE_URL='https://localhost:6789' \
E2E_LOCAL_BILLING_STUBS=true \
E2E_DEVICE_KEY='local-device-key' \
E2E_GATE_KEY='local-gate-key' \
E2E_JAN_CODE='テスト用JANコード' \
make test-e2e
```

内部では次のコマンドを実行します。

```bash
cd backend && go test -tags=e2e ./test/e2e
```
