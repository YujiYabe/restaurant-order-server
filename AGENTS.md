# Codexリポジトリルール

- 日本語で対応すること。
- 作業開始時に `README.md`、対象ディレクトリ、既存テスト、現在の差分を確認すること。
- ユーザーの未コミット変更を保持し、依頼と無関係な差分を変更・削除しないこと。
- 既存機能と外部APIの互換性を壊さないこと。調査前に大規模な変更やリファクタリングを行わないこと。
- アプリケーションコード、設定、Docker、DBスキーマを変更する場合は、既存構成と命名を先に確認すること。

## プロジェクト概要

- restaurant-order-server は、店舗の席、注文、在庫、会計、成人確認を管理するGo製Webサーバー。
- 主なコードは `backend/` にあり、Echo、GORM、SQLite、Redis、Firestore、SSEを使用する。
- スマレジ、Azure Document Intelligence、AWS Rekognitionなどの外部サービスと連携する。
- アプリケーションは設定やWeb資産を相対パスで参照するため、実行時の作業ディレクトリに注意すること。
- SSEのエンドポイントは `/v1/sse/app` で、長時間接続を前提に扱う。

## 起動・テスト・lint

- 開発起動: リポジトリルートで `./start.sh`
- ビルド: `make build`
- Firestore Emulator起動: `make firestore-up`
- Firestore Emulator停止: `make firestore-down`
- 全テスト: `cd backend && go test ./...`
- race detector: `cd backend && go test -race ./...`
- vet: `cd backend && go vet ./...`
- lint: `make lint`
- 変更範囲に対応する小さいテストを先に実行し、完了前に原則 `cd backend && go test ./...` を実行すること。
- 並行処理、共有状態、SSEを変更した場合は、可能な限り `go test -race ./...` も実行すること。

## Goのコード規約

- `gofmt` を適用すること。
- 既存の4層構成と依存方向を尊重する。
  - `backend/internal/1_framework`: Web、DB、外部サービス、シリアライズなどの技術詳細
  - `backend/internal/2_adapter`: controller、gateway、presenter
  - `backend/internal/3_usecase`: アプリケーションのユースケース
  - `backend/internal/4_domain`: ドメインモデルとビジネスルール
- ドメイン層からフレームワーク層、アダプター層、ユースケース層を参照しないこと。
- `json`、`yaml`、`gorm` など外部表現・永続化のタグや変換責務をドメインモデルへ持ち込まないこと。DTOや永続化モデルは境界側に定義すること。
- HTTP、Redis、Firestore、SQLiteなどの形式をドメイン型へ直接結合せず、境界で明示的に変換すること。
- 主要関数は既存方針に合わせ、必要な場合は `context.Context` を第一引数にすること。
- メソッドレシーバー名は原則 `receiver` を使うこと。
- 戻り値は周辺コードのスタイルに合わせ、必要に応じて名前付き戻り値を使うこと。
- エラーを握りつぶさず、呼び出し元で判断できるように返すか、運用で追える文脈とともに記録すること。
- goroutineを追加する場合は、終了条件、`context.Context` によるキャンセル、リーク、panicの影響を確認すること。

## ドメイン状態と並行処理

- `Domain` の可変状態は `mutex` で保護すること。ロックなしで共有マップ、スライス、席、注文を読み書きしないこと。
- 参照系メソッドは、呼び出し側から内部状態を変更できないスナップショットを返すこと。
- 注文追加と在庫減算、注文取消と在庫復元など、整合性が必要な複合更新は同一ロック内で不可分に実行すること。
- ロック保持中に外部I/O、長時間処理、ブロッキング送信を行わないこと。
- 席状態、注文状態、`BillingStage` の遷移ルールを推測で緩和しないこと。
- `SaveInMemory` を変更する場合は、Redisから復元した未会計注文の在庫予約と、自動更新goroutineの停止・再起動を確認すること。

## DB・外部サービス境界

- SQLite用のGORMモデルとタグは `backend/internal/1_framework/db/sqlite` に置き、ドメインモデルとの変換を境界で行うこと。
- JSON/YAMLのDTOと変換は `backend/internal/1_framework/transport` または利用するフレームワーク境界に置くこと。
- 既存のJSON/YAMLフィールド名、Redis保存形式、DBテーブル名・カラム名を変更する場合は、後方互換性と既存データの移行を確認すること。
- 外部API呼び出しをテストで実行しないこと。必要に応じてスタブ、エミュレーター、ローカルDBを使用すること。
- スマレジ会計処理は再試行される可能性があるため、`TransactionHeadId` と `BillingStage` による再開・多重実行防止を維持すること。

## SSE固有の注意点

- SSEを通常の短時間HTTPリクエストとして扱わないこと。
- 接続終了、クライアント切断、heartbeat、書き込みエラー、goroutine終了を確認すること。
- `Agents` などの共有状態へアクセスする場合は、既存のmutex方針を維持すること。
- チャネル送信で業務処理を無期限にブロックしないこと。
- 正しい業務状態をSSE接続内だけに保持せず、ドメイン状態または永続ストレージから復元できる設計を維持すること。

## 変更してはいけない範囲

- `.env`、Firestoreサービスアカウント、秘密鍵、パスワード、トークン、外部サービスの認証情報をコミットしないこと。
- `backend/.env`、`backend/secret/`、`backend/openssl/server.key` の内容を出力・変更しないこと。明示的な依頼がある場合も秘密情報を露出させないこと。
- `backend/db/engine/sqlite/*.sqlite3` は店舗データを含む可能性があるため、明示的な依頼なしに更新・削除しないこと。
- Firestore EmulatorのDockerボリュームやRedisの保存データを、明示的な依頼なしに削除しないこと。
- `make removeall`、Docker prune、DB初期化などの破壊的操作を自動実行しないこと。
- 実環境のFirestore、Azure、AWS、スマレジへ副作用のあるリクエストを送信しないこと。

## ドキュメントと完了報告

- 起動方法、設定、API仕様、運用手順を変更した場合は `README.md` も更新すること。
- 新しい環境変数を追加した場合は、秘密値を含めず `backend/.env.example` と関連文書を更新すること。
- 完了報告には、主要な変更内容、変更ファイル、実行したテスト・lint、未実行または未解決事項を含めること。
