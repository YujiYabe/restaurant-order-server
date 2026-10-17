# restaurant-order-server

restaurant-order-serverは、店舗の席、注文、在庫、会計、成人確認を管理するGo製のWebサーバーです。ファミリーレストラン向けのテーブルオーダーと、店内の着席客向けセルフオーダーに対応しています。着席後に発行されるQRコードから利用する方式であり、店外からのモバイルオーダーは対象としていません。店舗端末と利用者端末向けのWeb画面およびAPIを提供し、スマレジ、Firestore、Redis、Azure Document Intelligenceなどの外部サービスと連携します。

名称の`restaurant-order-server`は、対象となる飲食店内の業務を`restaurant`、注文受付から提供・会計までの中心機能を`order`、Web画面とAPIを提供するバックエンドを`server`として、その役割を直接表したものです。特定の店舗やサービスのブランド名ではなく、店外注文サービスや本番インフラ一式を意味する名称でもありません。

- リポジトリ: <https://github.com/YujiYabe/restaurant-order-server>
- バグ報告・機能提案: <https://github.com/YujiYabe/restaurant-order-server/issues>
- 脆弱性の非公開報告: [SECURITY.md](SECURITY.md)
- コントリビューションガイド: [CONTRIBUTING.md](CONTRIBUTING.md)
- 行動規範: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## 目的と対象範囲

このプロジェクトの目的は、飲食店内で着席している利用者からの注文を起点に、席の利用状況、注文、在庫、商品の提供、会計までの店舗業務を一つのWebサーバーで扱うことです。店舗スタッフ向け画面、利用者向けセルフオーダー画面、API、およびリアルタイム更新を提供し、Goによる店舗注文システムの実装・検証・カスタマイズの基盤として利用することを想定しています。

対象範囲は、単一店舗内で運用する次の業務です。

- 着席済みの利用者に対する、テーブル端末またはQRコード経由の店内注文
- 受付、ホール、厨房、サイネージ、商品管理で使用する店舗スタッフ向け画面
- 席、注文、在庫、提供状況、チャージ料金、会計進捗の管理
- 酒類注文時の成人確認と、スマレジ、Firestore、Redisなどへの連携
- 架空のサンプル設定とローカルサービスを使った開発・動作検証

店外注文、配送、予約、オンライン決済、顧客アカウント、複数店舗の一括管理は対象外です。具体的な対応機能と制約は以下に記載します。また、このリポジトリが提供するのはアプリケーションと開発用構成であり、実店舗へそのまま導入できるマネージドサービスや本番インフラ一式ではありません。

## できること

- 店舗のテーブル端末から、席に紐づく注文を登録する
- 席ごとに発行した権限キーを含むQRコードから、着席客のスマートフォンで商品を注文する
- 空席、利用中、会計中、清掃中、停止中という席の状態を管理する
- 注文数と在庫を検証し、注文受付と在庫の確保を一体として処理する
- 注文を準備中、提供済み、取消済みの状態で管理し、取消時に在庫を戻す
- 厨房、ホール、受付、注文画面へSSEで状態更新を配信する
- 席の利用時間と利用人数に応じてチャージ料金を計算する
- 酒類注文時に、本人確認書類のOCRと顔照合を利用して成人確認を行う
- 会計内容をスマレジへ連携し、注文履歴をFirestoreへ保存する
- 席の利用状況をRedisへ保存し、サーバー再起動後に復元する
- 商品と在庫をSQLiteで管理する

## 現在できないこと

- 店外からのテイクアウト注文
- 配送先を指定するデリバリー注文
- 受取日時または来店日時を指定する事前注文・席予約
- クレジットカードなどによるオンライン決済
- 顧客アカウントの登録、ログイン、ポイント、顧客別注文履歴の管理
- テーブル端末とスマートフォンのどちらから登録された注文かを、注文データへ記録すること
- 複数店舗を1つのサーバープロセスで同時に管理すること

「モバイルオーダー」は、店外注文ではなく、店内で着席済みの利用者がQRコードから利用するセルフオーダーを指します。

## アーキテクチャ

サーバーは、技術詳細をドメインから分離する4層構成です。HTTPリクエストの基本的な流れと、永続化・外部連携の流れは次のとおりです。

```text
Web / API
  ↓
1_framework/web → 2_adapter/controller → 3_usecase → 4_domain
                                           ↓
                                  2_adapter/gateway
                                           ↓
                         1_framework/db・external

3_usecase → 2_adapter/presenter → SSE・店舗モニター
```

| 層 | 責務 |
| --- | --- |
| `1_framework` | EchoによるHTTP/SSE、GORM、SQLite、Redis、Firestore、外部APIなどの技術詳細を実装する |
| `2_adapter` | Web入力とユースケースの橋渡し、保存先・外部サービスのインターフェース、表示用出力を扱う |
| `3_usecase` | 注文、会計、商品更新、成人確認など、一連のアプリケーション処理を調整する |
| `4_domain` | 席、注文、在庫、会計状態などのモデルとビジネスルールを保持する |

依存関係は外側から内側へ向けます。ドメイン層はEcho、GORM、JSON、Redisなどの外部形式へ依存せず、DTOや永続化モデルとの変換は境界側で行います。`backend/internal/init.go`がDB・外部サービスの実装を生成してcontrollerへ渡し、起動時にSQLiteとRedisからドメインのインメモリ状態を構築します。

## 主要ディレクトリ

```text
backend/
├── cmd/                 サーバーのエントリーポイント
│   └── db/              SQLite DB生成・更新用の独立Go module
├── config/              店舗別の環境設定と席の初期設定
├── db/engine/sqlite/    商品・タグのTSVと店舗別SQLite DB
├── internal/
│   ├── 1_framework/     Web、DB、外部サービスとの接続、外部形式のDTO
│   ├── 2_adapter/       コントローラー、ゲートウェイ、プレゼンター
│   ├── 3_usecase/       ユースケース
│   ├── 4_domain/        ドメインモデルとビジネスルール
│   ├── env/             実行環境別の設定読み込み
│   └── pkg/             層に属さない小さな共通処理
├── openssl/             ローカルHTTPS用の証明書
├── script/              起動診断と開発ツールの設定
├── secret/              Firestoreの認証情報（Git管理対象外）
├── test/e2e/            ローカルバックエンドに対するHTTP/SSE E2Eテスト
└── web/                 HTML、JavaScript、CSS、画像、音声などのWeb資産
```

リポジトリルートの`compose.yaml`はRedisとFirestore Emulator、`docker/firestore/`はEmulatorイメージ、`makefile`と`start.sh`はセットアップ・起動操作を管理します。

## データ保存先

| データ | 保存先 | 内容と扱い |
| --- | --- | --- |
| 実行中の業務状態 | `4_domain.Domain`のメモリー | 席、注文、在庫、会計進捗を処理する。プロセス終了で消えるため、起動時にSQLiteとRedisから復元する |
| 商品・タグ・在庫 | `backend/db/engine/sqlite/<STORE_ID>.sqlite3` | 店舗単位のSQLite DB。TSVから生成し、実店舗データを含み得るためGit管理しない |
| 席・未会計注文・会計進捗 | Redis DB 0の`placeList`キー | JSONとして保存し、再起動時に復元する。ローカルComposeでは`restaurant-order-server_redis-data` volumeへAOFを保存する |
| 会計済み注文履歴 | 注文用Firestoreプロジェクト | 日付ごとのcollectionへ保存する。ローカルではFirestore Emulatorを使用する |
| 成人確認属性・顔画像・酒類注文履歴 | 成人確認用Firestoreプロジェクト | 利用者UUIDごとのcollectionへ保存する。注文履歴とは別プロジェクトを使用する |
| 店舗設定・席の初期値 | `backend/config/<STORE_ID>_environment.yaml`、`<STORE_ID>_table_list.yaml` | Redisに席情報がない初回起動時などに読み込む設定ファイル |
| ローカルFirestoreデータ | `restaurant-order-server_firestore-data` Docker volume | Emulator終了後も保持される。実環境のFirestoreデータとは分離される |

Redisの`placeList`は通常停止・再起動で保持されますが、volumeを削除すると失われます。Firestore Emulatorも同様です。SQLite DB、Redis volume、Firestore volumeの削除や初期化は通常のセットアップ手順では行いません。TLS鍵とサービスアカウントJSONは実行設定・認証情報であり、業務データと同様に公開リポジトリへ含めません。

## Go moduleパス

公開時のGo moduleパスは、GitHub上の配置に合わせて次のように定めます。

| 対象 | 公開moduleパス |
| --- | --- |
| メインサーバー（`backend/go.mod`） | `github.com/YujiYabe/restaurant-order-server/backend` |
| DB更新ツール（`backend/cmd/db/go.mod`） | `github.com/YujiYabe/restaurant-order-server/backend/cmd/db` |

module宣言と内部importは上記の公開パスへ移行済みです。DB更新ツールは独立したmoduleとして依存関係を管理しており、メインサーバーの内部packageには依存しません。

## ローカル成果物とComposeの名称

既存の開発手順と保存データとの互換性を維持するため、次の名称を継続して使用します。

| 対象 | 採用する名称 | 方針 |
| --- | --- | --- |
| ビルド済みサーバーバイナリ | `backend/app` | 既存の `make build` と起動手順を維持するため変更しない |
| Docker Compose project | `restaurant-order-server` | リポジトリのディレクトリ名から決まる現在の名称を正式な名称とする |
| Compose内のRedis volume | `redis-data` | 論理名を変更しない |
| Compose内のFirestore Emulator volume | `firestore-data` | 論理名を変更しない |

通常のclone先では、Docker Engine上のvolume名はCompose project名を含む `restaurant-order-server_redis-data` と `restaurant-order-server_firestore-data` になります。project名やvolume名を変更すると既存データが自動では引き継がれないため、移行手順を用意せずに変更しません。clone先のディレクトリ名を変更する場合は、同じvolumeを使用するため `docker compose --project-name restaurant-order-server ...` を指定してください。

## 既存利用者向けの移行

この公開リポジトリはコピー元のGit履歴を引き継いでいません。コピー元のcloneでremoteだけを変更してpullすると履歴が一致しないため、公開版を新しいディレクトリへcloneしてください。

```bash
git clone https://github.com/YujiYabe/restaurant-order-server.git
cd restaurant-order-server
```

コピー元へ独自のコード変更がある場合は、秘密情報や実データを除外したパッチとして新しいcloneへ再適用します。旧 `.git`、`.env`、`backend/secret/`、TLS秘密鍵、SQLite DBはコピーしません。

Go moduleパスは次のように変更されています。

| 対象 | 旧パス | 新パス |
| --- | --- | --- |
| メインサーバー | `backend` | `github.com/YujiYabe/restaurant-order-server/backend` |
| DB更新ツール | `backend/cmd/db` | `github.com/YujiYabe/restaurant-order-server/backend/cmd/db` |
| 内部package | `backend/internal/...` | `github.com/YujiYabe/restaurant-order-server/backend/internal/...` |

独自変更を再適用する場合は、2つの `go.mod` とGoソース内のimportを新パスへ合わせた後、それぞれのmoduleで依存関係とテストを確認します。

```bash
cd backend
go mod tidy
go test ./...

cd cmd/db
go mod tidy
go test ./...
```

バイナリ名 `backend/app`、Compose project名 `restaurant-order-server`、volumeの論理名は変更していません。同じCompose project名で利用していた環境では既存volumeを継続利用できます。コピー元で別のproject名を使っていた場合は、先に `docker volume ls` で対象を確認してバックアップし、従来のproject名を `docker compose --project-name <従来のproject名> ...` に指定してください。volumeの削除や初期化は移行手順に含みません。

## 必須サービスと任意サービス

ここでの「必須」は、このREADMEの最小構成でサーバーを起動し、基本的な席・注文・在庫管理を確認するために必要なものを指します。「任意」は、対応する外部連携機能を利用する場合だけ必要です。

### 最小構成で必須

| 種別 | 要件 | 用途 |
| --- | --- | --- |
| 開発環境 | Go 1.27.1以上、Cコンパイラー、Make、OpenSSL | ビルド、CGOを使用するSQLite、ローカルHTTPS証明書の生成 |
| コンテナ環境 | Docker Engine、Docker Compose v2 | RedisとFirestore Emulatorの起動 |
| SQLite | リポジトリ内のサンプルTSVから生成するローカルDB | 商品と在庫の管理。別途SQLiteサーバーを導入する必要はない |
| Redis | Docker Composeで起動するローカルRedis | 席と注文の一時状態の保存・復元 |
| Firestore | Docker Composeで起動するFirestore Emulator | 注文履歴などの保存。最小構成ではGoogle Cloudの認証情報は不要 |
| インターネット接続 | 初回セットアップおよびWeb画面の外部配信資産取得時 | Goツール、コンテナイメージ、Bootstrap、Alpine.jsの取得 |

### 機能利用時のみ必要

| サービス | 必要になる機能 | 必要な設定 |
| --- | --- | --- |
| Google Cloud Firestore | Emulatorではなく実環境へ注文履歴・成人確認情報を保存する場合 | 2つのプロジェクトとサービスアカウントJSON |
| スマレジ | 実際の会計・販売情報をスマレジへ連携する場合 | 契約ID、クライアントID、クライアントシークレット |
| Azure Document Intelligence | 酒類注文で本人確認書類をOCRする場合 | 解析エンドポイントとサブスクリプションキー |
| AWS Rekognition | 酒類注文で本人確認書類と撮影画像の顔照合を行う場合 | AWSの認証情報（リージョンは現在 `ap-northeast-1` 固定） |
| 会計連携先（Gate） | 会計完了情報を別サービスへ通知する場合 | 通知先URL、パス、認証キー |

ローカル環境では酒類注文を無効化し、スマレジと会計連携先への送信をローカル応答で代替します。このため、最小構成では上記の任意サービスへ接続せず、`backend/.env.example`の公開ダミー値で起動できます。実運用では利用する機能に対応した認証情報へ置き換えてください。

## 本番利用前のセキュリティ確認

このリポジトリに含まれるCompose構成、自己署名証明書、サンプル設定、公開ダミー認証情報はローカル開発専用です。本番利用では、運用者が脅威分析と法令・契約要件の確認を行い、少なくとも次の対策を実施してください。

### ネットワークとTLS

- アプリケーションをインターネットへ直接公開せず、ファイアウォール、VPN、リバースプロキシなどで店舗端末・管理者・必要な連携元だけに接続元を制限する
- 信頼された認証局が発行した有効な証明書を使用し、TLSの終端箇所、更新、自動失効監視、HTTPからHTTPSへの転送、HSTSなどを運用環境に合わせて設定する
- Redis、SQLite、Firestoreおよび管理・診断用エンドポイントを公開ネットワークから分離する。特に開発用`compose.yaml`のFirestore Emulatorポートは本番で公開しない
- `/v1/sse/app`は長時間接続になるため、プロキシのタイムアウト、同時接続数、バッファリング、切断検知を調整し、接続数と帯域を監視する
- リバースプロキシまたは同等の境界で、リクエストボディサイズ、タイムアウト、同時接続数、レート制限を設定し、DoSと総当たりを緩和する

### 認証情報とアクセス制御

- `ENVIRONMENT`を`prod`にし、`backend/.env.example`の値を使用せず、Basic認証、端末キー、連携APIキー、外部API資格情報を十分な長さのランダムな値へ置き換える
- 認証情報とFirestoreサービスアカウントJSONはシークレット管理サービスから実行時に渡し、ソース管理、コンテナイメージ、ログ、バックアップへ含めない。ファイルで渡す場合は所有者と権限を最小化する
- スマレジ、Google Cloud、Azure、AWSなどの資格情報には最小権限、利用元制限、利用量アラートを設定し、定期ローテーションと漏えい時の失効手順を用意する
- Basic認証、端末キー、席ごとのQRコードを第三者と共有しない。端末の紛失、従業員の異動、QRコードの露出時に速やかに無効化・再発行できる運用を定める
- `device_key`はURLクエリに含まれるため、ブラウザー履歴、Referer、アクセスログ、画面共有へ残る可能性がある。ログのマスキング、外部リンク制限、端末管理を行い、より強い認証や短期トークンが必要な環境では本番導入前に追加実装とレビューを行う
- 管理画面、商品更新、席操作、会計、SSEを含む全ルートについて、役割別の権限と到達可能範囲が店舗の要件を満たすことを検証する。Basic認証だけを従業員ごとの本人識別や詳細な権限管理の代替としない

### 個人情報、データ保護、運用

- 注文、売上、顔画像、本人確認書類由来の情報を取り扱う目的、保存先、保存期間、削除手順、本人への通知・同意を定め、適用される個人情報保護法、酒類販売ルール、外部サービス規約を専門家と確認する
- 成人確認機能を有効にする前に、誤判定時のスタッフ確認、撮影データの最小化、アクセス権、保存期間、削除、利用者からの問い合わせ対応を設計する。自動判定だけで法的要件を満たすと仮定しない
- SQLite、Redis、Firestoreの暗号化、アクセス制御、バックアップ、復元試験、保持期間、確実な削除を構成する。バックアップにも本番データと同じ保護を適用する
- 認証失敗、管理操作、会計、設定変更、外部APIエラーを、秘密値や本人確認画像を記録せず監査できるようにする。ログへのアクセス制限、改ざん防止、保持期間、アラートを設定する
- 外部サービス障害、二重会計、ネットワーク断、プロセス再起動、データ破損を想定した手順を用意し、`TransactionHeadId`と`BillingStage`による再開を含む会計復旧を実データ以外で訓練する
- Go、コンテナイメージ、JavaScript依存、OSを継続的に更新し、脆弱性スキャン、`go test -race ./...`、`go vet ./...`、`make lint`を本番反映前に実行する
- 本番相当の分離環境で、認証・認可、CSRF、XSS、入力値、ファイル／画像処理、エラー表示、SSE切断、負荷、バックアップ復元、外部APIの失敗を含むセキュリティレビューと受入試験を行う

### 免責事項

本ソフトウェアはMIT Licenseに基づき、現状有姿で提供されます。商品性、特定目的への適合性、可用性、セキュリティ、法令適合性、会計・在庫・成人確認結果の正確性を保証するものではありません。READMEの確認事項は一般的な最低限の例であり、網羅的なセキュリティ監査、法務・税務・個人情報保護・酒類販売に関する助言を代替しません。

実店舗での利用、構成、監視、バックアップ、インシデント対応、外部サービスとの契約および適用法令への対応は導入・運用者の責任です。本番導入前に、独立したセキュリティレビューと、必要に応じて法律・会計・業務の専門家による確認を行ってください。障害や判定結果にかかわらず、人が注文、会計、在庫、年齢確認を訂正・停止できる手順を維持してください。

## セットアップ

### 最小構成で起動する

ローカルのサンプル店舗を起動する最短手順です。「最小構成で必須」に記載した開発環境、コンテナ環境、インターネット接続を事前に利用可能にしてください。スマレジ、Azure、AWS、実環境のFirestore、会計連携先は不要です。コマンドはリポジトリルートで実行します。

```bash
git clone https://github.com/YujiYabe/restaurant-order-server.git
cd restaurant-order-server
cp backend/.env.example backend/.env
make install-tools

mkdir -p backend/openssl
openssl req -x509 -newkey rsa:2048 -sha256 -days 365 -nodes \
  -keyout backend/openssl/server.key \
  -out backend/openssl/server.crt \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
chmod 600 backend/openssl/server.key

make prepare
./start.sh
```

`make prepare`は、サンプルSQLite DBがない場合だけ新規作成し、RedisとFirestore EmulatorをDocker Composeで起動します。既存のDBやDocker volumeは初期化しません。起動後は `https://localhost:6789/` を開きます。自己署名証明書を使用するため、ブラウザーにはローカル証明書の警告が表示されます。

`backend/.env.example`の認証値はローカル動作確認専用の公開ダミー値です。本番環境では使用しないでください。`backend/.env`、生成した秘密鍵、SQLite DBはGit管理対象外です。停止する場合は実行中の `./start.sh` を `Ctrl+C` で終了し、次のコマンドでローカルサービスを停止します。保存済みvolumeは削除されません。

```bash
make redis-down
make firestore-down
```

### 外部APIを使わずにローカルで試す

最小構成では、実際のスマレジ、Google Cloud Firestore、Azure Document Intelligence、AWS Rekognition、会計連携先へ接続せず、席・注文・在庫・会計の基本フローを確認できます。次のローカル用設定を維持してください。

- `backend/config/0_environment.yaml`の`ENVIRONMENT`を`local`、`IS_PERMIT_ALCOHOL`を`false`のまま使用する
- `backend/.env.example`を`backend/.env`へコピーし、公開ダミー値を実サービスの認証情報へ置き換えない
- `FIRESTORE_EMULATOR_HOST=127.0.0.1:57242`と`demo-`で始まる2つのFirestoreプロジェクトIDを維持する
- `make prepare`でRedisとFirestore Emulatorを起動してから`./start.sh`を実行する

ローカル設定での外部連携は次のように置き換わります。

| 本番時の連携 | ローカルでの動作 |
| --- | --- |
| スマレジへの仮販売登録 | 外部送信せず、`local-`で始まる取引IDを生成する |
| 会計連携先への完了通知 | 外部送信を省略する |
| Google Cloud Firestore | Docker上のFirestore Emulatorへ保存する |
| Azureによる本人確認書類OCR | 酒類注文を無効にするため、通常の操作手順では呼び出さない |
| AWS Rekognitionによる顔照合 | 酒類注文を無効にするため、通常の操作手順では呼び出さない |

起動後に`https://localhost:6789/`を開き、受付画面で席を利用中にし、店舗端末または着席客向け画面から酒類以外の商品を注文します。厨房・ホール画面で提供状態を更新し、会計まで進めることで、SQLite、Redis、Firestore Emulator、SSEを使った基本フローを確認できます。

外部APIを呼ばない確認では、酒類の本人確認エンドポイントを直接呼び出したり、`IS_PERMIT_ALCOHOL`を`true`へ変更したりしないでください。AzureとAWSの連携自体を試すには、別途実サービスの認証情報が必要です。また、BootstrapとAlpine.jsはCDN配信の静的ファイルを使用するため、ブラウザーで画面を正しく表示するにはインターネット接続が必要です。これは業務データを送信する外部API連携ではありません。

以下では、各設定と操作を詳しく説明します。

### 1. 開発ツールをインストールする

```bash
make install-tools
```

Air、golangci-lint、swagが `backend/bin` にインストールされます。

### 2. 店舗設定を用意する

環境変数 `STORE_ID` に対応する次のファイルを `backend/config` に用意します。

```text
<STORE_ID>_environment.yaml
<STORE_ID>_table_list.yaml
```

既存店舗の設定例は `backend/config/0_environment.yaml` と `backend/config/0_table_list.yaml` です。`ENVIRONMENT` には `local`、`dev`、または `prod` を指定します。

### 3. 認証情報を設定する

ローカル動作確認では、公開ダミー値を含むサンプルをコピーします。

```bash
cp backend/.env.example backend/.env
```

実環境では `backend/.env.example` に記載された変数名を使用し、値は実行環境またはシークレット管理サービスから設定します。

```bash
export STORE_ID=0
export SERVER_PORT=6789
export BASIC_AUTH_USERNAME='...'
export BASIC_AUTH_PASSWORD='...'
export FROM_GATE_KEY='...'
export FROM_DEVICE_KEY='...'
export CONTRACT_ID='...'
export CLIENT_ID='...'
export CLIENT_SECRET='...'
export GATE_TRANSACTION_INFO_BASE_URL='http://localhost:8080'
export GATE_TRANSACTION_INFO_PATH='/transaction-info'
export AZURE_DOCUMENT_INTELLIGENCE_URL='...'
export AZURE_SUBSCRIPTION_KEY='...'
```

サーバーの待受ポートは `SERVER_PORT` で指定します。

Azure Document Intelligenceを利用する場合は、実行環境に対応する `backend/internal/env/local.yaml`、`dev.yaml`、`prod.yaml` の `AZURE_DOCUMENT_INTELLIGENCE_URL` に解析エンドポイントを設定します。環境変数で同名の値を指定した場合は、ファイルの設定より優先されます。

会計完了情報の通知先は、同じ環境別設定ファイルの `GATE_TRANSACTION_INFO_BASE_URL` と `GATE_TRANSACTION_INFO_PATH` で指定します。同名の環境変数による上書きも可能です。

アプリケーションは `.env` ファイルを自動では読み込みません。実際の認証情報をリポジトリへコミットしないでください。

`./start.sh`、`make prepare`、`make redis-up`は利便性のため `backend/.env` を読み込み、子プロセスへ環境変数として渡します。バイナリを直接起動する場合は、同じ変数をあらかじめ実行環境へ設定してください。

ローカル開発では、Firestore EmulatorをDockerで起動します。

```bash
make firestore-up
```

Emulatorは次のホスト側ポートを使用します。コンテナ内部ではFirebase標準のポートを使用しますが、ホスト側では競合しにくい動的・プライベートポートへ割り当てています。

| 用途 | ホスト側 | コンテナ内部 |
| --- | --- | --- |
| Firestore API | `127.0.0.1:57242` | `8080` |
| Emulator管理画面 | `http://127.0.0.1:54371` | `4000` |
| Emulator Hub | `127.0.0.1:53971` | `4400` |
| UI WebSocket | `127.0.0.1:58000` | `9150` |

`backend/.env`には次の設定が必要です。

```dotenv
FIRESTORE_EMULATOR_HOST=127.0.0.1:57242
FIRESTORE_ORDER_PROJECT_ID=demo-restaurant-order-server-order
FIRESTORE_CERTIFICATE_PROJECT_ID=demo-restaurant-order-server-certificate
```

`FIRESTORE_EMULATOR_HOST`には `http://` などのURLスキームを付けません。注文履歴と成人確認情報は異なるプロジェクトIDへ保存されます。EmulatorのデータはDockerボリューム `restaurant-order-server_firestore-data` にエクスポートされ、通常の停止・再起動では維持されます。

Emulatorを停止する場合は次のコマンドを実行します。

```bash
make firestore-down
```

実環境のFirestoreを使用する場合は、`FIRESTORE_EMULATOR_HOST`および2つのプロジェクトIDを設定せず、サービスアカウントJSONを次のパスへ配置します。

```text
backend/secret/firestore_order.json
backend/secret/firestore_certificate.json
```

### 4. SQLite DBを初期化・更新する

`backend/db/engine/sqlite` のTSVを読み込み、`STORE_ID` に対応するSQLite DBを作成または更新します。

```bash
STORE_ID=0 make sqlite-init
```

生成される `backend/db/engine/sqlite/<STORE_ID>.sqlite3` はGit管理対象外です。

### 5. Redisを起動する

アプリケーションを起動する前に、次のコマンドでRedisを起動します。データはDockerボリューム `restaurant-order-server_redis-data` に保持されます。

```bash
make redis-up
```

接続先の既定値は `127.0.0.1:6379` です。公開ポートを変更する場合は、`backend/.env`に次のように設定します。

```dotenv
REDIS_PORT=16379
```

`./start.sh`と`make redis-up`は、この値をComposeの公開ポートとアプリケーションの接続先へ反映します。外部Redisなどホスト名を含めて変更する場合は `REDIS_ADDR`を指定でき、`REDIS_PORT`より優先されます。

```dotenv
REDIS_ADDR=redis.example.internal:6379
```

`./start.sh`は起動診断を行い、Redisなどのローカル依存を準備してからアプリケーションを起動します。

ローカル環境の会計処理では、スマレジの仮販売登録と会計連携サービスへの通知をローカル応答で代替します。外部サービスへ取引データは送信されません。注文履歴はFirestore Emulatorへ保存され、在庫はローカルのSQLiteへ反映されます。

### 6. 起動条件を診断・準備する

設定ファイル、開発ツール、証明書、SQLite DB、Docker、Redis、Firestore Emulator、ポート競合をまとめて確認できます。`doctor`は診断だけを行い、ファイルやコンテナを変更しません。

```bash
make doctor
```

不足しているSQLite DBの新規作成と、設定済みのローカルサービスの起動まで行う場合は、次を実行します。既存のSQLite DB、Redis、Firestore Emulatorの保存データは初期化しません。

```bash
make prepare
```

`./start.sh`は内部で`make prepare`相当の処理を実行します。ポート`6789`が使用中の場合は、該当するプロセスを表示してアプリケーション起動前に終了します。

## 起動方法

アプリケーションは設定やWeb資産を相対パスで参照するため、`backend` を作業ディレクトリとして実行します。

開発時は、リポジトリのルートからAirによるホットリロードを開始できます。

```bash
./start.sh
```

起動時にトップページのURLが`restaurant-order-server Portal: https://localhost:6789/`の形式で表示されます。ローカル環境では認証キーをログへ出力せず、このURLをそのまま開けます。

通常のバイナリをビルドして起動する場合は、次のコマンドを実行します。

```bash
make build
cd backend
./app
```

ローカル環境では `backend/openssl/server.crt` と `backend/openssl/server.key` を使用してHTTPSサーバーを起動します。既定のローカルポートは `6789` です。

クライアントへのリアルタイム通知にはServer-Sent Eventsを使用し、エンドポイントは `/v1/sse/app` です。

## テストと静的解析

```bash
cd backend
go test ./...
go test -race ./...
go vet ./...
```

golangci-lintはリポジトリのルートで実行します。

```bash
make lint
```

起動済みのローカルバックエンドに対するHTTP/SSE E2Eテストは、開発専用のRedisを使用して次のように実行します。

```bash
E2E_LOCAL_BILLING_STUBS=true \
E2E_DEVICE_KEY='バックエンドのFROM_DEVICE_KEYと同じ値' \
E2E_GATE_KEY='バックエンドのFROM_GATE_KEYと同じ値' \
E2E_JAN_CODE='テストDBに在庫がある酒類以外のJANコード' \
make test-e2e
```

E2Eテストは接続先をlocalhostまたはloopback IPに制限します。会計フローではローカルスタブ、Firestore Emulator、破棄可能なSQLiteテストDBが必要です。詳細は `backend/test/e2e/README.md` を参照してください。

## フロントエンド

店舗画面はサーバーレンダリングを維持し、軽量なAlpine.jsをリアクティブ状態とライフサイクルの基盤として使用します。DOM操作の段階的な移行を支える軽量互換層にはCashを使用し、jQuery本体には依存しません。Bootstrap 5.3.1とAlpine.js 3.17.4はバージョン固定したjsDelivr CDNからSRI検証付きで取得します。Cash、QR生成、Font Awesomeの実行時ファイルはローカル配信します。再配布条件を確認できない顔認識モデルは同梱せず、成人確認画面は撮影後のサーバー側検証を使用します。

トップページは `/?device_key=端末キー` です。受付、スタッフ、厨房、サイネージ、商品管理などの共通画面と、各席の厨房・スタッフ・モバイル注文画面へ移動できます。

### スクリーンショットの公開基準

READMEやドキュメントへ画面を掲載する場合は、実店舗や実利用者の環境ではなく、`STORE_ID=0`のサンプル設定と架空データを使用します。公開前に、次の情報が画像内や画像メタデータへ残っていないことを確認してください。

- 氏名、顔、本人確認書類、生年月日、住所などの個人情報
- 実店舗名、住所、電話番号、営業時間、席名など店舗を特定できる情報
- QRコード、端末キー、APIキー、認証情報、UUID、取引ID、注文番号
- 実際の注文内容、売上、在庫数、会計情報、管理画面のデータ
- 内部URL、IPアドレス、ホスト名、ファイルパス、ブラウザーの他タブや通知
- EXIFの位置情報、撮影日時、端末情報、作成者などの画像メタデータ

必要な箇所は切り抜くか、不透明な単色で完全に置き換えます。ぼかしや半透明の塗りつぶしだけに依存しません。QRコードは一部だけでなく全体を除去します。画像を書き出した後に目視とメタデータ確認を行い、元画像ではなく除去済みの画像だけを公開対象へ追加してください。現在のREADMEにはスクリーンショットを掲載していません。

## ドメインロジック

`backend/internal/4_domain` は、席、注文、在庫、会計、商品、および成人確認に関するビジネスルールとインメモリ状態を管理します。

### 同期とスナップショット

`Domain` の可変状態は `mutex` で保護されます。参照メソッドは、ネストしたマップやスライスを含むコピーを返し、呼び出し側から内部状態が更新されることを防ぎます。注文追加と在庫減算、注文取消と在庫復元は、それぞれ同一ロック内で実行されます。

### 外部形式との境界

ドメインモデルはJSON、YAML、GORMなどの技術固有タグを持ちません。HTTP、Redis、Firestore向けのDTOと変換は `backend/internal/1_framework/transport`、SQLite向けの永続化モデルとGORMタグは `backend/internal/1_framework/db/sqlite` で管理します。

### 席の状態遷移

席の通常フローは次のとおりです。

```text
Cleaning -> Opening -> Using -> Billing -> Cleaning
```

`Cleaning` から `Stopping` への遷移も許可されます。`Using` から `Cleaning` への直接遷移は、運用上の強制解放として許可されます。その他の不正な遷移は状態を変更せず `false` を返します。

`Opening` へ遷移すると、取引ID、会計進捗、アレルギー、言語、注文一覧、およびインターネット注文用権限キーが次の利用に向けて初期化されます。

### 注文と在庫

`TryAddOrder` は `Using` 中の席と正数の注文だけを受け付けます。通常商品は登録済み在庫と注文数を検証してから減算します。JANコード0はチャージ明細として扱い、在庫管理の対象外です。準備中の注文だけが提供済みまたは取消へ遷移でき、取消時は確保済み在庫を1度だけ復元します。

`SaveInMemory` はSQLiteの在庫を基準とし、Redisに保存された `Using` 中の注文と、在庫更新完了前の `Billing` 中注文を差し引きます。これにより、再起動や商品更新後に確保済み在庫が再販売されることを防ぎます。

### 会計の再開情報

`BillingStage` はスマレジ取引作成、注文保存、取引通知、在庫更新の完了位置を表します。ユースケース層は `TransactionHeadId` と `BillingStage` を保存し、再試行時に完了済みの外部処理を繰り返さないために使用します。

### チャージ料金

チャージ機能が有効な場合、`Using` 中の席について入店時刻からの経過時間を `ChargeBaseMinute` で区切り、区切り数、`ChargeBasePrice`、収容人数から料金を算出します。自動更新は1秒間隔で実行され、新しいマスター情報の読み込み時に旧ループが停止されます。

### 成人確認

書類の生年月日は和暦から西暦へ変換し、現在日時時点で20歳以上かを判定します。書類確認後の顔登録許可は10分間かつ1回限りです。顔照合に成功したUUIDの酒類注文許可は、日本時間の翌日0時まで有効です。
