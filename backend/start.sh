#!/bin/sh

if [ -f .env ]; then
	set -a
	. ./.env
	set +a
fi

: "${STORE_ID:=0}"
export STORE_ID

exec bin/air -c ./script/air.toml
