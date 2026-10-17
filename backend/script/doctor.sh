#!/bin/sh

set -u

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
BACKEND_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
REPOSITORY_DIR=$(CDPATH= cd -- "$BACKEND_DIR/.." && pwd)
PREPARE=false
ERRORS=0

case "${1:-}" in
	"") ;;
	--prepare) PREPARE=true ;;
	-h|--help)
		printf '%s\n' "Usage: $0 [--prepare]"
		printf '%s\n' "  引数なし   起動条件を診断します（変更は行いません）"
		printf '%s\n' "  --prepare  不足しているローカル依存を安全に準備します"
		exit 0
		;;
	*)
		printf '[ERROR] 不明な引数です: %s\n' "$1" >&2
		exit 2
		;;
esac

if [ -f "$BACKEND_DIR/.env" ]; then
	set -a
	# shellcheck disable=SC1091
	. "$BACKEND_DIR/.env"
	set +a
fi

: "${STORE_ID:=0}"
: "${REDIS_PORT:=6379}"
export STORE_ID REDIS_PORT

ok() {
	printf '[OK] %s\n' "$1"
}

warn() {
	printf '[WARN] %s\n' "$1"
}

fail() {
	printf '[ERROR] %s\n' "$1" >&2
	ERRORS=$((ERRORS + 1))
}

require_command() {
	if command -v "$1" >/dev/null 2>&1; then
		ok "$1 コマンド"
	else
		fail "$1 コマンドが見つかりません"
	fi
}

require_file() {
	if [ -f "$1" ]; then
		ok "$2"
	else
		fail "$2 がありません: $1"
	fi
}

valid_port() {
	case "$1" in
		""|*[!0-9]*) return 1 ;;
	esac
	[ "$1" -ge 1 ] && [ "$1" -le 65535 ]
}

printf 'restaurant-order-server 起動診断 (STORE_ID=%s)\n' "$STORE_ID"

require_command go
require_command docker

if ! valid_port "$REDIS_PORT"; then
	fail "REDIS_PORT は1〜65535の整数で指定してください"
fi

require_file "$BACKEND_DIR/config/${STORE_ID}_environment.yaml" "店舗環境設定"
require_file "$BACKEND_DIR/config/${STORE_ID}_table_list.yaml" "店舗テーブル設定"
require_file "$BACKEND_DIR/db/engine/sqlite/product.tsv" "商品マスター"
require_file "$BACKEND_DIR/db/engine/sqlite/tag_primary.tsv" "商品プライマリタグ"
require_file "$BACKEND_DIR/db/engine/sqlite/tag_secondary.tsv" "商品セカンダリタグ"
require_file "$BACKEND_DIR/openssl/server.crt" "HTTPS証明書"
require_file "$BACKEND_DIR/openssl/server.key" "HTTPS秘密鍵"

if [ ! -f "$BACKEND_DIR/openssl/server.crt" ] || [ ! -f "$BACKEND_DIR/openssl/server.key" ]; then
	printf '%s\n' '[INFO] HTTPS証明書または秘密鍵がない場合は、次のコマンドでローカル開発用の自己署名証明書を生成してください:' >&2
	printf '  cd %s\n' "$REPOSITORY_DIR" >&2
	printf '%s\n' \
		'  mkdir -p backend/openssl' \
		'  openssl req -x509 -newkey rsa:2048 -sha256 -days 365 -nodes \' \
		'    -keyout backend/openssl/server.key \' \
		'    -out backend/openssl/server.crt \' \
		'    -subj "/CN=localhost" \' \
		'    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"' \
		'  chmod 600 backend/openssl/server.key' \
		'  make doctor' >&2
	printf '%s\n' '[INFO] server.keyは秘密鍵です。Gitへ追加したり他者と共有したりしないでください。' >&2
fi

if [ -x "$BACKEND_DIR/bin/air" ]; then
	ok "Air"
else
	fail "Airがありません。make install-tools を実行してください"
fi

if command -v ss >/dev/null 2>&1; then
	LISTENER=$(ss -ltnp '( sport = :6789 )' 2>/dev/null | sed -n '2,$p')
	if [ -n "$LISTENER" ]; then
		fail "アプリケーションのポート6789は既に使用されています"
		printf '%s\n' "$LISTENER" >&2
	else
		ok "アプリケーションのポート6789"
	fi
else
	warn "ssコマンドがないため、ポート6789の競合を確認できません"
fi

if [ "$ERRORS" -ne 0 ]; then
	printf '[ERROR] 診断で%s件の問題が見つかりました\n' "$ERRORS" >&2
	exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
	fail "Docker Compose v2を利用できません"
elif ! docker info >/dev/null 2>&1; then
	fail "Dockerデーモンへ接続できません"
else
	ok "Docker Compose / Dockerデーモン"
fi

SQLITE_FILE="$BACKEND_DIR/db/engine/sqlite/${STORE_ID}.sqlite3"
if [ -f "$SQLITE_FILE" ]; then
	ok "SQLite DB（既存DBは更新しません）"
elif [ "$PREPARE" = true ]; then
	printf '[PREPARE] SQLite DBを新規作成します: %s\n' "$SQLITE_FILE"
	if (cd "$BACKEND_DIR/cmd/db" && go run .); then
		ok "SQLite DBを作成しました"
	else
		fail "SQLite DBを作成できませんでした"
	fi
else
	fail "SQLite DBがありません。make prepare を実行してください: $SQLITE_FILE"
fi

REDIS_MANAGED=true
case "${REDIS_ADDR:-}" in
	""|"127.0.0.1:${REDIS_PORT}"|"localhost:${REDIS_PORT}") ;;
	*) REDIS_MANAGED=false ;;
esac

if [ "$REDIS_MANAGED" = true ]; then
	if [ "$PREPARE" = true ]; then
		printf '[PREPARE] Redisを起動してヘルスチェックします\n'
		if (cd "$REPOSITORY_DIR" && docker compose up -d --wait redis); then
			ok "Redis"
		else
			fail "Redisを起動できませんでした（REDIS_PORT=$REDIS_PORT の競合も確認してください）"
		fi
	elif (cd "$REPOSITORY_DIR" && docker compose ps --status running --services redis 2>/dev/null | grep -qx redis) &&
		(cd "$REPOSITORY_DIR" && docker compose exec -T redis redis-cli ping 2>/dev/null | grep -qx PONG); then
		ok "Redis"
	else
		fail "Redisが利用できません。make prepare を実行してください"
	fi
else
	warn "REDIS_ADDRで外部Redisが指定されているため、疎通確認と起動を省略します"
fi

if [ -n "${FIRESTORE_EMULATOR_HOST:-}" ]; then
	if [ "$PREPARE" = true ]; then
		printf '[PREPARE] Firestore Emulatorを起動してヘルスチェックします\n'
		if (cd "$REPOSITORY_DIR" && docker compose up -d --build --wait firestore); then
			ok "Firestore Emulator"
		else
			fail "Firestore Emulatorを起動できませんでした"
		fi
	elif (cd "$REPOSITORY_DIR" && docker compose ps --status running --services firestore 2>/dev/null | grep -qx firestore); then
		ok "Firestore Emulator"
	else
		fail "FIRESTORE_EMULATOR_HOSTが設定されていますがEmulatorが停止しています。make prepare を実行してください"
	fi
else
	warn "Firestore Emulatorは未設定です。実環境認証の有無は秘密情報保護のため確認しません"
fi

if [ "$ERRORS" -ne 0 ]; then
	printf '[ERROR] 診断で%s件の問題が見つかりました\n' "$ERRORS" >&2
	exit 1
fi

ok "起動準備が整っています"
