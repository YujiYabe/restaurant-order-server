#!/bin/sh

set -eu

REPOSITORY_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

if [ -f "$REPOSITORY_DIR/backend/.env" ]; then
	set -a
	. "$REPOSITORY_DIR/backend/.env"
	set +a
fi

: "${STORE_ID:=0}"
export STORE_ID

"$REPOSITORY_DIR/backend/script/doctor.sh" --prepare

cd "$REPOSITORY_DIR/backend"
exec bin/air -c ./script/air.toml
